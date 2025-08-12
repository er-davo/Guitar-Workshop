package service

import (
	"api-gateway/internal/clients"
	"api-gateway/internal/models"
	"api-gateway/internal/repository"
	"bytes"
	"context"
	"database/sql"
	"errors"

	"github.com/supabase-community/supabase-go"
)

type TabService struct {
	repo           *repository.TabRepository
	supabaseClient *supabase.Client
	tabClient      clients.TabGenerator
	audioClient    clients.AudioSeparator
}

func NewTabService(
	repo *repository.TabRepository,
	supabaseClient *supabase.Client,
	tabClient clients.TabGenerator,
	audioClient clients.AudioSeparator,
) *TabService {
	return &TabService{
		repo:           repo,
		supabaseClient: supabaseClient,
		tabClient:      tabClient,
		audioClient:    audioClient,
	}
}

func (s *TabService) GenerateTab(ctx context.Context, audioFileName string, audioFileData []byte, separation bool) (*models.Tab, error) {
	if separation {
		separatedFiles, err := s.audioClient.SeparateAudio(ctx, audioFileName, audioFileData)
		if err != nil {
			return nil, err
		}
		otherStem, ok := separatedFiles["other"]
		if !ok {
			return nil, errors.New("audio separation result missing 'other' stem")
		}

		audioFileName = otherStem.FileName
		audioFileData = otherStem.AudioBytes
	}

	tabResp, err := s.tabClient.GenerateTab(ctx, audioFileName, audioFileData)
	if err != nil {
		return nil, err
	}

	tab := new(models.Tab)
	tab.Body = tabResp.Tab

	return tab, nil
}

func (s *TabService) CreateTab(ctx context.Context, tab *models.Tab) error {
	err := s.repo.Create(ctx, tab)
	if err != nil {
		if err.Error() == "The resource already exists" {
			return errors.New("tab with that name already exists")
		} else {
			return err
		}
	}
	_, err = s.supabaseClient.Storage.UploadFile(
		"tabs",
		tab.Path,
		bytes.NewReader([]byte(tab.Body)),
	)
	return err
}

func (s *TabService) DeleteTab(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *TabService) GetTabByID(ctx context.Context, id string) (*models.Tab, error) {
	tab, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("tab not found")
		} else {
			return nil, err
		}
	}
	buf, err := s.supabaseClient.Storage.DownloadFile("tabs", tab.Path)
	if err != nil {
		return nil, err
	}

	tab.Body = string(buf)

	return tab, nil
}

func (s *TabService) FindTabsByNameLike(ctx context.Context, name string) ([]*models.Tab, error) {
	return s.repo.FindByNameLike(ctx, name)
}
