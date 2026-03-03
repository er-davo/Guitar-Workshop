package service

import (
	"context"
	"orchestrator/internal/models"

	"go.uber.org/zap"
)

type TabGenService struct {
	repo TabGenTaskRepository

	log *zap.Logger
}

func NewTabGenService(
	repo TabGenTaskRepository,
	log *zap.Logger,
) *TabGenService {
	return &TabGenService{
		repo: repo,
		log:  log,
	}
}

func (s *TabGenService) HandleTabGenEvent(
	ctx context.Context,
	event *models.TabGenTaskRequestedEvent,
) (bool, error) {

	log := s.log.With(
		zap.String("tab_gen_task_id", event.ID),
	)

	log.Info("handling TabGenTaskRequestedEvent")

	tgTask, err := s.repo.Get(ctx, event.ID)
	if err != nil {
		log.Error("failed to get tab gen task from repository",
			zap.Error(err),
		)
		return false, err
	}

	var targetStatus models.Status
	fromStatus := models.StatusCreated

	if tgTask.AudioSepTaskID != nil {
		targetStatus = models.StatusWaitingForSeparation
		log.Debug("audio separation required for tab generation",
			zap.String("audio_sep_task_id", *tgTask.AudioSepTaskID),
		)
	} else {
		targetStatus = models.StatusPending
		log.Debug("no audio separation required, switching to pending")
	}

	log.Debug("attempting status transition",
		zap.String("from_status", string(fromStatus)),
		zap.String("to_status", string(targetStatus)),
	)

	ok, err := s.repo.TryUpdateStatus(
		ctx,
		event.ID,
		fromStatus,
		targetStatus,
		nil,
		nil,
	)
	if err != nil {
		log.Error("failed to update tab gen task status",
			zap.Error(err),
		)
		return false, err
	}

	if !ok {
		log.Warn("status transition skipped (already processed or changed)")
		return false, nil
	}

	log.Info("tab gen task status successfully updated",
		zap.String("new_status", string(targetStatus)),
	)

	okResponse := true

	if targetStatus == models.StatusWaitingForSeparation {
		okResponse = false
	}

	return okResponse, nil
}
