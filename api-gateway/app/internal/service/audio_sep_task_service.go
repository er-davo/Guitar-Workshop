package service

import (
	"bytes"
	"context"
	"path"
	"sort"
	"time"

	"api-gateway/internal/models"

	"github.com/er-davo/retry"
	"go.uber.org/zap"
)

type AudioSepTaskService struct {
	repo AudioSepTaskRepository

	storage             Storage
	audioBucket         string
	presignedExpiration time.Duration

	producer AudioSepTaskProducer

	retrier retry.Retrier

	log *zap.Logger
}

func NewAudioSepTaskService(
	repo AudioSepTaskRepository,
	storage Storage,
	audioBucket string,
	presignedExpiration time.Duration,
	audioSepTaskProducer AudioSepTaskProducer,
	retrier retry.Retrier,
	log *zap.Logger,
) *AudioSepTaskService {
	return &AudioSepTaskService{
		repo:                repo,
		storage:             storage,
		audioBucket:         audioBucket,
		presignedExpiration: presignedExpiration,
		producer:            audioSepTaskProducer,
		retrier:             retrier,
		log:                 log,
	}
}

func (s *AudioSepTaskService) PostAudioSepTask(ctx context.Context, taskData models.AudioSepTaskData) (*models.AudioSepTask, error) {
	task := &models.AudioSepTask{
		Status:           models.StatusPending,
		InputAudioName:   taskData.AudioFileName,
		SeparatedDirName: "separated",
	}

	if err := s.repo.Create(ctx, task); err != nil {
		s.log.Error("failed to create audio separation task", zap.String("audio_file", taskData.AudioFileName), zap.Error(err))
		return nil, err
	}
	s.log.Info("audio separation task created", zap.String("task_id", task.ID), zap.String("audio_file", taskData.AudioFileName))

	if err := s.retrier.Do(ctx, func(attempt int) error {
		innerErr := s.storage.UploadFile(
			ctx,
			s.audioBucket,
			task.AudioFileObjectName(),
			bytes.NewReader(taskData.Data),
		)
		if innerErr != nil {
			s.log.Warn("failed to upload audio file, retrying",
				zap.String("audio_file", taskData.AudioFileName),
				zap.Int("attempt", attempt),
				zap.Error(innerErr),
			)
		}
		return innerErr
	}); err != nil {
		s.log.Error("failed to upload audio file after retries", zap.String("task_id", task.ID), zap.String("audio_file", taskData.AudioFileName), zap.Error(err))
		return nil, err
	}
	s.log.Info("audio file uploaded successfully", zap.String("task_id", task.ID))

	if err := s.producer.Produce(ctx, &models.StartAudioSepTask{ID: task.ID}); err != nil {
		s.log.Error("failed to produce audio separation task message", zap.String("task_id", task.ID), zap.Error(err))
		return nil, err
	}
	s.log.Info("audio separation task message produced", zap.String("task_id", task.ID))

	return task, nil
}

func (s *AudioSepTaskService) Get(ctx context.Context, id string) (*models.AudioSepTask, error) {
	task, err := s.repo.Get(ctx, id)
	if err != nil {
		s.log.Error("failed to get audio separation task", zap.String("task_id", id), zap.Error(err))
		return nil, err
	}

	if task.Status != models.StatusDone {
		s.log.Info("audio separation task not completed yet", zap.String("task_id", id))
		return task, nil
	}

	filesInfo, err := s.storage.ListFilesByPrefix(ctx, s.audioBucket, task.SeparatedPrefix())
	if err != nil {
		s.log.Error("failed to list separated audio files", zap.String("task_id", id), zap.Error(err))
		return nil, err
	}

	if len(filesInfo) == 0 {
		s.log.Warn("no separated audio files found", zap.String("task_id", id))
	}

	sort.Slice(filesInfo, func(i, j int) bool {
		return *filesInfo[i].Key < *filesInfo[j].Key
	})

	task.SeparatedAudioSignedURLs = make(map[string]string)
	for _, fileInfo := range filesInfo {
		url, err := s.storage.PresignedGet(ctx, s.audioBucket, *fileInfo.Key, s.presignedExpiration)
		if err != nil {
			s.log.Error("failed to generate presigned URL", zap.String("task_id", id), zap.String("object", *fileInfo.Key), zap.Error(err))
			return nil, err
		}
		name := path.Base(*fileInfo.Key)
		task.SeparatedAudioSignedURLs[name] = url.String()
	}

	s.log.Info("audio separation task ready", zap.String("task_id", id), zap.Int("files_count", len(task.SeparatedAudioSignedURLs)))
	return task, nil
}
