package service

import (
	"api-gateway/internal/models"
	"bytes"
	"context"
	"time"

	"github.com/er-davo/retry"
	"go.uber.org/zap"
)

type TabService struct {
	tabRepo   TabRepository
	trManager TxManager

	storage             Storage
	tabsBucket          string
	presignedExpiration time.Duration

	retrier retry.Retrier

	log *zap.Logger
}

func NewTabService(
	tabRepo TabRepository,
	genTaskRepo TabGenTaskRepository,
	trManager TxManager,
	storage Storage,
	tabsBucket string,
	presignedExpiration time.Duration,
	retrier retry.Retrier,
	log *zap.Logger,
) *TabService {
	return &TabService{
		tabRepo:             tabRepo,
		storage:             storage,
		trManager:           trManager,
		tabsBucket:          tabsBucket,
		presignedExpiration: presignedExpiration,
		retrier:             retrier,
		log:                 log,
	}
}

// CreateTab uploads the tab file to storage and saves metadata in DB
func (s *TabService) CreateTab(ctx context.Context, tabc *models.TabCreate) error {
	err := s.retrier.Do(ctx, func(attempt int) error {
		innerErr := s.storage.UploadFile(ctx,
			s.tabsBucket,
			tabc.Path,
			bytes.NewReader([]byte(tabc.Body)),
		)
		if innerErr != nil {
			s.log.Warn("upload tab file failed, retrying",
				zap.String("tab_path", tabc.Path),
				zap.Int("attempt", attempt),
				zap.Error(innerErr),
			)
		}
		return innerErr
	})
	if err != nil {
		s.log.Error("failed to upload tab file after retries", zap.String("tab_path", tabc.Path), zap.Error(err))
		return err
	}

	tab := &models.Tab{
		Name: tabc.Name,
		Path: tabc.Path,
	}

	if err := s.tabRepo.Create(ctx, tab); err != nil {
		s.log.Error("failed to save tab metadata", zap.String("tab_name", tabc.Name), zap.Error(err))
		return err
	}

	tabc.ID = tab.ID
	s.log.Info("tab created successfully", zap.String("tab_id", tab.ID), zap.String("tab_name", tab.Name))

	return nil
}

// DeleteTab marks tab as deleted and removes file from storage
func (s *TabService) DeleteTab(ctx context.Context, id string) error {
	tab, err := s.tabRepo.Get(ctx, id)
	if err != nil {
		s.log.Error("failed to fetch tab for deletion", zap.String("tab_id", id), zap.Error(err))
		return err
	}

	if err := s.tabRepo.MarkDeleted(ctx, id); err != nil {
		s.log.Error("failed to mark tab as deleted", zap.String("tab_id", id), zap.Error(err))
		return err
	}

	err = s.retrier.Do(ctx, func(attempt int) error {
		removeErr := s.storage.RemoveFile(ctx, s.tabsBucket, tab.Path)
		if removeErr != nil {
			s.log.Warn("failed to remove tab file, retrying",
				zap.String("tab_id", tab.ID),
				zap.String("tab_path", tab.Path),
				zap.Int("attempt", attempt),
				zap.Error(removeErr),
			)
		}
		return removeErr
	})
	if err != nil {
		s.log.Error("tab marked deleted but failed to remove file from storage",
			zap.String("tab_id", tab.ID),
			zap.String("tab_path", tab.Path),
			zap.Error(err),
		)
		return nil // background cleanup can handle file removal
	}

	s.log.Info("tab deleted successfully", zap.String("tab_id", tab.ID))
	return nil
}

// GetTabByID fetches tab and generates presigned URL
func (s *TabService) GetTabByID(ctx context.Context, id string) (*models.Tab, error) {
	tab, err := s.tabRepo.Get(ctx, id)
	if err != nil {
		s.log.Error("failed to fetch tab by ID", zap.String("tab_id", id), zap.Error(err))
		return nil, err
	}

	url, err := s.storage.PresignedGet(ctx, s.tabsBucket, tab.Path, s.presignedExpiration)
	if err != nil {
		s.log.Error("failed to generate presigned URL", zap.String("tab_id", tab.ID), zap.Error(err))
		return nil, err
	}

	tab.PresignedURL = url.String()
	s.log.Info("tab fetched successfully", zap.String("tab_id", tab.ID), zap.String("tab_name", tab.Name))
	return tab, nil
}

// FindTabsByNameLike searches tabs by name
func (s *TabService) FindTabsByNameLike(ctx context.Context, name string) ([]*models.Tab, error) {
	s.log.Info("searching tabs", zap.String("query", name))
	tabs, err := s.tabRepo.FindByNameLike(ctx, name)
	if err != nil {
		s.log.Error("failed to search tabs", zap.String("query", name), zap.Error(err))
		return nil, err
	}
	s.log.Info("tabs search result", zap.String("query", name), zap.Int("count", len(tabs)))
	return tabs, nil
}
