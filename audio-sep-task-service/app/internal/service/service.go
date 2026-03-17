package service

import (
	"audio-sep-task-service/internal/models"
	"bytes"
	"context"
	"path"
	"sort"
	"time"

	"github.com/er-davo/gwcontracts/audiosep"
	"github.com/er-davo/retry"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AudioSepTaskService struct {
	audiosep.UnimplementedAudioSepTaskServiceServer

	repo AudioSepTaskRepository

	storage             Storage
	audioBucket         string
	presignedExpiration time.Duration

	producer AudioSepTaskProducer
	retrier  retry.Retrier

	log *zap.Logger
}

func NewAudioSepTaskService(
	repo AudioSepTaskRepository,
	storage Storage,
	audioBucket string,
	presignedExpiration time.Duration,
	producer AudioSepTaskProducer,
	retrier retry.Retrier,
	log *zap.Logger,
) *AudioSepTaskService {
	return &AudioSepTaskService{
		repo:                repo,
		storage:             storage,
		audioBucket:         audioBucket,
		presignedExpiration: presignedExpiration,
		producer:            producer,
		retrier:             retrier,
		log:                 log,
	}
}

func (s *AudioSepTaskService) PostAudioSepTask(ctx context.Context, req *audiosep.PostAudioSepTaskRequest) (*audiosep.PostAudioSepTaskResponse, error) {
	start := time.Now()

	s.log.Info("creating audio separation task",
		zap.String("audio_file", req.FileName),
		zap.Int("size_bytes", len(req.Data)),
	)

	task := &models.AudioSepTask{
		Status:           models.StatusPending,
		InputAudioName:   req.FileName,
		SeparatedDirName: "separated",
	}

	if err := s.repo.Create(ctx, task); err != nil {
		s.log.Error("failed to create task",
			zap.String("audio_file", req.FileName),
			zap.Error(err),
		)
		return nil, status.Error(codes.Internal, "failed to create task")
	}

	s.log.Info("task created",
		zap.String("task_id", task.ID),
	)

	uploadStart := time.Now()

	if err := s.retrier.Do(ctx, func(attempt int) error {
		err := s.storage.UploadFile(
			ctx,
			s.audioBucket,
			task.AudioFileObjectName(),
			bytes.NewReader(req.Data),
		)

		if err != nil {
			s.log.Warn("upload attempt failed",
				zap.String("task_id", task.ID),
				zap.Int("attempt", attempt),
				zap.Error(err),
			)
		}

		return err
	}); err != nil {

		s.log.Error("upload failed after retries",
			zap.String("task_id", task.ID),
			zap.Duration("duration", time.Since(uploadStart)),
			zap.Error(err),
		)

		return nil, status.Error(codes.Internal, "failed to upload audio")
	}

	s.log.Info("audio uploaded",
		zap.String("task_id", task.ID),
		zap.Duration("duration", time.Since(uploadStart)),
	)

	if err := s.producer.Produce(ctx, &models.StartAudioSepTask{ID: task.ID}); err != nil {
		s.log.Error("failed to produce event",
			zap.String("task_id", task.ID),
			zap.Error(err),
		)
		return nil, status.Error(codes.Internal, "failed to produce event")
	}

	s.log.Info("audio separation event produced",
		zap.String("task_id", task.ID),
		zap.Duration("total_time", time.Since(start)),
	)

	return &audiosep.PostAudioSepTaskResponse{
		Task: toAudioSepTask(task),
	}, nil
}

func (s *AudioSepTaskService) GetAudioSepTask(ctx context.Context, req *audiosep.GetAudioSepTaskRequest) (*audiosep.GetAudioSepTaskResponse, error) {
	s.log.Info("fetching task",
		zap.String("task_id", req.Id),
	)

	task, err := s.repo.Get(ctx, req.Id)
	if err != nil {
		s.log.Error("failed to get task",
			zap.String("task_id", req.Id),
			zap.Error(err),
		)
		return nil, status.Error(codes.Internal, "failed to get task")
	}

	s.log.Info("task fetched",
		zap.String("task_id", req.Id),
		zap.String("status", task.Status.String()),
	)

	if task.Status != models.StatusDone {
		return &audiosep.GetAudioSepTaskResponse{
			Task: toAudioSepTask(task),
		}, nil
	}

	filesInfo, err := s.storage.ListFilesByPrefix(ctx, s.audioBucket, task.SeparatedPrefix())
	if err != nil {
		s.log.Error("failed to list separated files",
			zap.String("task_id", req.Id),
			zap.Error(err),
		)
		return nil, status.Error(codes.Internal, "failed to list stems")
	}

	s.log.Info("separated files listed",
		zap.String("task_id", req.Id),
		zap.Int("files_count", len(filesInfo)),
	)

	sort.Slice(filesInfo, func(i, j int) bool {
		return *filesInfo[i].Key < *filesInfo[j].Key
	})

	task.SeparatedAudioSignedURLs = make(map[string]string)

	for _, fileInfo := range filesInfo {
		url, err := s.storage.PresignedGet(
			ctx,
			s.audioBucket,
			*fileInfo.Key,
			s.presignedExpiration,
		)

		if err != nil {
			s.log.Error("failed to generate presigned url",
				zap.String("task_id", req.Id),
				zap.String("object", *fileInfo.Key),
				zap.Error(err),
			)
			return nil, status.Error(codes.Internal, "failed to generate url")
		}

		name := path.Base(*fileInfo.Key)
		task.SeparatedAudioSignedURLs[name] = url.String()
	}

	return &audiosep.GetAudioSepTaskResponse{
		Task: toAudioSepTask(task),
	}, nil
}

func (s *AudioSepTaskService) UploadStemsURLs(ctx context.Context, req *audiosep.UploadStemsURLsRequest) (*audiosep.UploadStemsURLsResponse, error) {
	s.log.Info("request stems upload urls",
		zap.String("task_id", req.Id),
		zap.Int("stems_count", len(req.StemsNames)),
	)

	task, err := s.repo.Get(ctx, req.Id)
	if err != nil {
		s.log.Error("failed to get task",
			zap.String("task_id", req.Id),
			zap.Error(err),
		)
		return nil, status.Error(codes.Internal, "failed to get task")
	}

	res := &audiosep.UploadStemsURLsResponse{
		StemsUploadUrls: make(map[string]string),
	}

	for _, stem := range req.StemsNames {
		object := path.Join(task.SeparatedPrefix(), stem)

		url, err := s.storage.PresignedPut(
			ctx,
			s.audioBucket,
			object,
			s.presignedExpiration,
		)

		if err != nil {
			s.log.Error("failed to generate upload url",
				zap.String("task_id", req.Id),
				zap.String("object", object),
				zap.Error(err),
			)
			return nil, status.Error(codes.Internal, "failed to generate upload url")
		}

		res.StemsUploadUrls[stem] = url.String()
	}

	s.log.Info("stems upload urls generated",
		zap.String("task_id", req.Id),
	)

	return res, nil
}

func (s *AudioSepTaskService) TryUpdate(ctx context.Context, req *audiosep.TryUpdateRequest) (*audiosep.TryUpdateResponse, error) {
	if req.Payload == nil {
		return nil, status.Error(codes.InvalidArgument, "payload required")
	}

	s.log.Info("try update task",
		zap.String("task_id", req.Id),
		zap.String("status_from", req.StatusFrom.String()),
		zap.String("status_to", req.StatusTo.String()),
	)

	from := convertStatusFromProto(req.StatusFrom)
	to := convertStatusFromProto(req.StatusTo)

	var (
		sepDirName   *string
		errorMessage *string
	)

	switch pl := req.Payload.(type) {
	case *audiosep.TryUpdateRequest_DoneData:
		if pl.DoneData == nil {
			return nil, status.Error(codes.InvalidArgument, "done_data missing")
		}

		sepDirName = &pl.DoneData.SeparatedDirName

		s.log.Info("marking task as done",
			zap.String("task_id", req.Id),
			zap.String("separated_dir", pl.DoneData.SeparatedDirName),
		)

	case *audiosep.TryUpdateRequest_ErrorData:
		if pl.ErrorData == nil {
			return nil, status.Error(codes.InvalidArgument, "error_data missing")
		}

		errorMessage = &pl.ErrorData.ErrorMessage

		s.log.Warn("marking task as failed",
			zap.String("task_id", req.Id),
			zap.String("error_message", pl.ErrorData.ErrorMessage),
		)

	case *audiosep.TryUpdateRequest_ProcessingData:
		s.log.Info("marking task as processing",
			zap.String("task_id", req.Id),
		)

	default:
		return nil, status.Error(codes.InvalidArgument, "unknown payload type")
	}

	ok, err := s.repo.TryUpdate(ctx, req.Id, from, to, sepDirName, errorMessage)
	if err != nil {
		s.log.Error("failed to update task",
			zap.String("task_id", req.Id),
			zap.String("status_from", req.StatusFrom.String()),
			zap.String("status_to", req.StatusTo.String()),
			zap.Error(err),
		)
		return nil, status.Error(codes.Internal, "failed to update task")
	}

	if !ok {
		s.log.Warn("task update rejected",
			zap.String("task_id", req.Id),
			zap.String("status_from", req.StatusFrom.String()),
			zap.String("status_to", req.StatusTo.String()),
		)
	} else {
		s.log.Info("task updated",
			zap.String("task_id", req.Id),
			zap.String("status_to", req.StatusTo.String()),
		)
	}

	return &audiosep.TryUpdateResponse{
		Success: ok,
	}, nil
}

func toAudioSepTask(task *models.AudioSepTask) *audiosep.AudioSepTask {
	return &audiosep.AudioSepTask{
		Id:               task.ID,
		InputAudioName:   task.InputAudioName,
		Status:           convertStatusToProto(task.Status),
		SeparatedDirName: task.SeparatedDirName,
		ErrorMessage:     task.Error,
	}
}

func convertStatusToProto(status models.Status) audiosep.Status {
	switch status {
	case models.StatusCreated:
		return audiosep.Status_CREATED
	case models.StatusPending:
		return audiosep.Status_PENDING
	case models.StatusProcessing:
		return audiosep.Status_PROCESSING
	case models.StatusDone:
		return audiosep.Status_DONE
	case models.StatusFailed:
		return audiosep.Status_FAILED
	default:
		return audiosep.Status_CREATED
	}
}

func convertStatusFromProto(status audiosep.Status) models.Status {
	switch status {
	case audiosep.Status_CREATED:
		return models.StatusCreated
	case audiosep.Status_PENDING:
		return models.StatusPending
	case audiosep.Status_PROCESSING:
		return models.StatusProcessing
	case audiosep.Status_DONE:
		return models.StatusDone
	case audiosep.Status_FAILED:
		return models.StatusFailed
	default:
		return models.StatusCreated
	}
}
