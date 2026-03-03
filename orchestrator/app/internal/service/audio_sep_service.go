package service

import (
	"context"
	"errors"
	"orchestrator/internal/models"
	"orchestrator/internal/repository"

	"go.uber.org/zap"
)

type AudioSepService struct {
	audioSepRepo AudioSepTaskRepository
	tabGenRepo   TabGenTaskRepository

	log *zap.Logger
}

func NewAudioSepService(
	audioSeprRepo AudioSepTaskRepository,
	tabGenRepo TabGenTaskRepository,
	log *zap.Logger,
) *AudioSepService {
	return &AudioSepService{
		audioSepRepo: audioSeprRepo,
		tabGenRepo:   tabGenRepo,
		log:          log,
	}
}
func (s *AudioSepService) HandleAudioSepTaskCompletedEvent(
	ctx context.Context,
	event *models.AudioSepTaskCompletedEvent,
) (string, error) {
	log := s.log.With(
		zap.String("audio_sep_task_id", event.ID),
	)

	log.Info("handling AudioSepTaskCompletedEvent")

	asTask, err := s.audioSepRepo.Get(ctx, event.ID)
	if err != nil {
		log.Error("failed to get audio sep task from repository",
			zap.Error(err),
		)
		return "", err
	}

	log.Debug("audio sep task loaded")

	tgTask, err := s.tabGenRepo.GetByAudioSepTaskID(ctx, asTask.ID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			log.Warn("no tab gen task linked to audio sep task")
			return "", nil
		}

		log.Error("failed to get tab gen task by audio sep id",
			zap.Error(err),
		)
		return "", err
	}

	log = log.With(
		zap.String("tab_gen_task_id", tgTask.ID),
	)

	log.Info("tab gen task associated with completed audio sep")

	return tgTask.ID, nil
}
