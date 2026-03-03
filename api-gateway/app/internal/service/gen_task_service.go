package service

import (
	"bytes"
	"context"
	"time"

	"api-gateway/internal/models"

	"github.com/er-davo/retry"
	"go.uber.org/zap"
)

type GenTaskService struct {
	tabGenRepo TabGenTaskRepository
	tabRepo    TabRepository

	trManager TxManager

	storage             Storage
	audioBucket         string
	tabBucket           string
	presignedExpiration time.Duration

	audioSepTaskPoster AudioSepTaskPoster
	producer           TabGenTaskProducer

	retrier retry.Retrier

	log *zap.Logger
}

func NewGenTaskService(
	tabGenRepo TabGenTaskRepository,
	tabRepo TabRepository,
	trManager TxManager,
	storage Storage,
	audioBucket string,
	tabBucket string,
	presignedExpiration time.Duration,
	audioSepTaskPoster AudioSepTaskPoster,
	producer TabGenTaskProducer,
	retrier retry.Retrier,
	log *zap.Logger,
) *GenTaskService {
	return &GenTaskService{
		tabGenRepo:          tabGenRepo,
		tabRepo:             tabRepo,
		trManager:           trManager,
		storage:             storage,
		audioBucket:         audioBucket,
		tabBucket:           tabBucket,
		presignedExpiration: presignedExpiration,
		audioSepTaskPoster:  audioSepTaskPoster,
		producer:            producer,
		retrier:             retrier,
		log:                 log,
	}
}

func (s *GenTaskService) PostGenTask(ctx context.Context, taskData models.PostGenTaskData) (*models.TabGenTask, error) {
	task := &models.TabGenTask{
		Status:         models.StatusCreated,
		InputAudioName: taskData.AudioFileName,
	}

	if err := s.tabGenRepo.Create(ctx, task); err != nil {
		s.log.Error("failed to create tab generation task", zap.String("task_id", task.ID), zap.Error(err))
		return nil, err
	}

	if taskData.Separation {
		audioSepTask, err := s.audioSepTaskPoster.PostAudioSepTask(ctx, models.AudioSepTaskData{
			AudioFileName: taskData.AudioFileName,
			Data:          taskData.Data,
			Size:          taskData.Size,
		})
		if err != nil {
			s.log.Error("failed to post audio separation task", zap.String("audio_file", taskData.AudioFileName), zap.Error(err))
			return nil, err
		}
		task.AudioSepTaskID = &audioSepTask.ID

		err = s.tabGenRepo.SetAudioSeparationID(ctx, task.ID, *task.AudioSepTaskID)

		s.log.Info("audio separation task posted and id set to tab generation task", zap.String("task_id", task.ID), zap.String("audio_sep_task_id", audioSepTask.ID))
	} else {
		err := s.retrier.Do(ctx, func(attempt int) error {
			innerErr := s.storage.UploadFile(
				ctx,
				s.audioBucket,
				task.AudioFileObjectName(),
				bytes.NewReader(taskData.Data),
			)
			if innerErr != nil {
				s.log.Warn("uploading audio file failed, retrying",
					zap.String("audio_file", taskData.AudioFileName),
					zap.Int("attempt", attempt),
					zap.Error(innerErr),
				)
			}
			return innerErr
		})
		if err != nil {
			s.log.Error("failed to upload audio file after retries", zap.String("audio_file", taskData.AudioFileName), zap.Error(err))
			return nil, err
		}
		s.log.Info("audio file uploaded successfully", zap.String("audio_file", taskData.AudioFileName))
	}

	if err := s.producer.Produce(ctx, &models.StartTabGenTask{ID: task.ID}); err != nil {
		s.log.Error("failed to produce tab generation message", zap.String("task_id", task.ID), zap.Error(err))
		return nil, err
	}

	s.log.Info("tab generation task created", zap.String("task_id", task.ID))
	return task, nil
}

func (s *GenTaskService) Get(ctx context.Context, id string) (*models.TabGenResponse, error) {
	task, err := s.tabGenRepo.Get(ctx, id)
	if err != nil {
		s.log.Error("failed to get tab gen task", zap.String("task_id", id), zap.Error(err))
		return nil, err
	}

	if task.Status != models.StatusDone {
		s.log.Info("tab generation task not completed", zap.String("task_id", id))
		return &models.TabGenResponse{Task: task}, nil
	}

	if task.ResultTabName == nil {
		s.log.Warn("tab task done but result tab name is nil", zap.String("task_id", id))
		return &models.TabGenResponse{Task: task}, nil
	}

	url, err := s.storage.PresignedGet(ctx, s.tabBucket, task.ResultTabObjectName(), s.presignedExpiration)
	if err != nil {
		s.log.Error("failed to generate presigned URL", zap.String("task_id", id), zap.Error(err))
		return &models.TabGenResponse{Task: task}, err
	}

	tab := &models.Tab{
		ID:           task.ID,
		Name:         task.InputAudioName,
		Path:         task.ResultTabObjectName(),
		CreatedAt:    task.CreatedAt,
		PresignedURL: url.String(),
	}

	s.log.Info("tab ready for download", zap.String("task_id", task.ID), zap.String("tab_name", tab.Name))
	return &models.TabGenResponse{Task: task, Tab: tab}, nil
}

func (s *GenTaskService) SaveTab(ctx context.Context, taskID string, name string) error {
	return s.trManager.Do(ctx, func(ctx context.Context) error {
		task, err := s.tabGenRepo.Get(ctx, taskID)
		if err != nil {
			s.log.Error("failed to fetch tab generation task for save", zap.String("task_id", taskID), zap.Error(err))
			return err
		}

		if task.Status != models.StatusDone {
			return ErrTaskNotCompleted
		}

		if task.Saved {
			s.log.Info("tab task already saved", zap.String("task_id", taskID))
			return nil
		}

		if task.ResultTabName == nil {
			return ErrNoResultTabName
		}

		tab := &models.Tab{
			Name: name,
			Path: task.ResultTabObjectName(),
		}

		if err := s.tabRepo.Create(ctx, tab); err != nil {
			s.log.Error("failed to create tab from task", zap.String("task_id", taskID), zap.Error(err))
			return err
		}

		if err := s.tabGenRepo.MarkSaved(ctx, taskID); err != nil {
			s.log.Error("failed to mark tab generation task as saved", zap.String("task_id", taskID), zap.Error(err))
			return err
		}

		s.log.Info("tab saved successfully from task", zap.String("task_id", taskID), zap.String("tab_name", tab.Name))
		return nil
	})
}
