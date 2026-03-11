package service

import (
	"context"
	"tab-service/internal/models"
	"time"

	"github.com/er-davo/gwcontracts/tab"
	"github.com/er-davo/retry"
	"go.uber.org/zap"
)

type TabService struct {
	tab.UnimplementedTabServiceServer

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
	trManager TxManager,
	storage Storage,
	tabsBucket string,
	presignedExpiration time.Duration,
	retrier retry.Retrier,
	log *zap.Logger,
) *TabService {
	return &TabService{
		tabRepo:             tabRepo,
		trManager:           trManager,
		storage:             storage,
		tabsBucket:          tabsBucket,
		presignedExpiration: presignedExpiration,
		retrier:             retrier,
		log:                 log,
	}
}

func (s *TabService) SaveTab(ctx context.Context, req *tab.SaveTabRequest) (*tab.SaveTabResponse, error) {
	t := &models.Tab{
		Name: req.Name,
		Path: req.Path,
	}

	if err := s.tabRepo.Create(ctx, t); err != nil {
		s.log.Error("failed to save tab", zap.String("tab_id", t.ID), zap.Error(err))
		return nil, err
	}

	s.log.Info("tab saved successfully", zap.String("tab_id", t.ID), zap.String("tab_name", t.Name))
	return &tab.SaveTabResponse{
		Id: t.ID,
	}, nil
}

func (s *TabService) GetTab(ctx context.Context, req *tab.GetTabRequest) (*tab.GetTabResponse, error) {
	t, err := s.tabRepo.Get(ctx, req.Id)
	if err != nil {
		s.log.Error("failed to fetch tab by ID", zap.String("tab_id", req.Id), zap.Error(err))
		return nil, err
	}

	url, err := s.storage.PresignedGet(ctx, s.tabsBucket, t.Path, s.presignedExpiration)
	if err != nil {
		s.log.Error("failed to generate presigned URL", zap.String("tab_id", t.ID), zap.Error(err))
		return nil, err
	}

	t.PresignedURL = url.String()
	s.log.Info("tab fetched successfully", zap.String("tab_id", t.ID), zap.String("tab_name", t.Name))
	return &tab.GetTabResponse{
		Tab: &tab.Tab{
			Id:           t.ID,
			Name:         t.Name,
			Path:         t.Path,
			PresignedUrl: t.PresignedURL,
		},
	}, nil
}

func (s *TabService) DeleteTab(ctx context.Context, req *tab.DeleteTabRequest) (*tab.DeleteTabResponse, error) {
	t, err := s.tabRepo.Get(ctx, req.Id)
	if err != nil {
		s.log.Error("failed to fetch tab for deletion", zap.String("tab_id", req.Id), zap.Error(err))
		return nil, err
	}

	if err := s.tabRepo.MarkDeleted(ctx, req.Id); err != nil {
		s.log.Error("failed to mark tab as deleted", zap.String("tab_id", req.Id), zap.Error(err))
		return nil, err
	}

	err = s.retrier.Do(ctx, func(attempt int) error {
		removeErr := s.storage.RemoveFile(ctx, s.tabsBucket, t.Path)
		if removeErr != nil {
			s.log.Warn("failed to remove tab file, retrying",
				zap.String("tab_id", t.ID),
				zap.String("tab_path", t.Path),
				zap.Int("attempt", attempt),
				zap.Error(removeErr),
			)
		}
		return removeErr
	})
	if err != nil {
		s.log.Error("tab marked deleted but failed to remove file from storage",
			zap.String("tab_id", t.ID),
			zap.String("tab_path", t.Path),
			zap.Error(err),
		)
		return nil, nil // background cleanup can handle file removal
	}

	s.log.Info("tab deleted successfully", zap.String("tab_id", t.ID))
	return &tab.DeleteTabResponse{}, nil
}

func (s *TabService) SearchTabs(ctx context.Context, req *tab.SearchTabsRequest) (*tab.SearchTabsResponse, error) {
	s.log.Info("searching tabs", zap.String("query", req.NameQuery))
	tabs, err := s.tabRepo.Search(ctx, req.NameQuery, req.Limit, req.Offset)
	if err != nil {
		s.log.Error("failed to search tabs", zap.String("query", req.NameQuery), zap.Error(err))
		return nil, err
	}
	s.log.Info("tabs search result", zap.String("query", req.NameQuery), zap.Int("count", len(tabs)))

	rawTabsResponse := &tab.SearchTabsResponse{
		Tabs: make([]*tab.RawTab, len(tabs)),
	}

	for _, t := range tabs {
		rawTabsResponse.Tabs = append(rawTabsResponse.Tabs, &tab.RawTab{
			Id:   t.ID,
			Name: t.Name,
			Path: t.Path,
		})
	}

	return rawTabsResponse, nil
}
