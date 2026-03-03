package service

import (
	"context"
	"io"
	"path"
	"path/filepath"
	"strings"
	"time"

	"tabgen/internal/clients"
	"tabgen/internal/models"
	"tabgen/internal/music"
	"tabgen/internal/processor"
	note_analyzer "tabgen/internal/proto/note-analyzer"

	"go.uber.org/zap"
	"golang.org/x/sync/semaphore"
)

/*
last error caused fail:
22:36:55	ERROR	service/tab_gen_service.go:261	failed to get audio	{
	"error": "operation error S3: GetObject, https response error StatusCode: 404, RequestID: , HostID: , NoSuchKey: Object not found",
	"audio_object_name": "3a1a9413-7065-45f9-8f32-e60cdaca3e03/3a1a9413-7065-45f9-8f32-e60cdaca3e03/separated/other.wav"
}
*/

type TabGenServiceConfig struct {
	TaskTimeout    time.Duration
	DBTimeout      time.Duration
	StorageTimeout time.Duration
	MLTimeout      time.Duration
}

type TabGenService struct {
	cfg *TabGenServiceConfig

	tgRepo TabGenTaskRepository
	asRepo AudioSepTaskRepository

	mlSem      *semaphore.Weighted
	mlAnalyzer clients.NoteAnalyzer
	mlTimeout  time.Duration

	storage     Storage
	audioBucket string
	tabBucket   string

	tabProc processor.TabProcessor

	log *zap.Logger
}

func NewTabGenStartService(
	cfg *TabGenServiceConfig,
	tgRepo TabGenTaskRepository,
	asRepo AudioSepTaskRepository,
	mlSemSize int64,
	mlAnalyzer clients.NoteAnalyzer,
	strg Storage,
	audioBucket string,
	tabBucket string,
	tabProc processor.TabProcessor,
	log *zap.Logger,
) *TabGenService {
	return &TabGenService{
		cfg:         cfg,
		tgRepo:      tgRepo,
		asRepo:      asRepo,
		mlSem:       semaphore.NewWeighted(mlSemSize),
		mlAnalyzer:  mlAnalyzer,
		storage:     strg,
		audioBucket: audioBucket,
		tabBucket:   tabBucket,
		tabProc:     tabProc,
		log:         log,
	}
}

func (s *TabGenService) GenerateTab(ctx context.Context, task *models.TabGenTaskStartEvent) error {
	ctx, cancel := context.WithTimeout(ctx, s.cfg.TaskTimeout)
	defer cancel()

	start := time.Now()

	log := s.log.With(
		zap.String("task_id", task.ID),
	)

	log.Info("tab generation pipeline started")

	// --- STATUS UPDATE
	dbStart := time.Now()
	dbCtx, dbCancel := context.WithTimeout(ctx, s.cfg.DBTimeout)
	defer dbCancel()

	ok, err := s.startProcessing(dbCtx, task)
	if err != nil {
		log.Error("failed to switch task to processing",
			zap.Error(err),
			zap.Duration("db_duration", time.Since(dbStart)),
		)
		s.markError(ctx, task.ID, err.Error())
		return err
	}
	if !ok {
		log.Warn("task skipped due to status mismatch")
		return nil
	}

	log.Debug("status switched to processing",
		zap.Duration("db_duration", time.Since(dbStart)),
	)

	// --- LOAD AUDIO
	storageStart := time.Now()
	storageCtx, storageCancel := context.WithTimeout(ctx, s.cfg.StorageTimeout)
	defer storageCancel()

	audioName, audio, err := s.loadAudio(storageCtx, task.ID)
	if err != nil {
		log.Error("failed to load audio",
			zap.Error(err),
			zap.Duration("storage_duration", time.Since(storageStart)),
		)
		s.markError(ctx, task.ID, err.Error())
		return err
	}

	log.Info("audio loaded",
		zap.String("audio_name", audioName),
		zap.Duration("storage_duration", time.Since(storageStart)),
	)

	// --- ML
	mlStart := time.Now()
	mlCtx, mlCancel := context.WithTimeout(ctx, s.cfg.MLTimeout)
	defer mlCancel()

	notes, err := s.analyze(mlCtx, audio)
	if err != nil {
		log.Error("ml analyze failed",
			zap.Error(err),
			zap.Duration("ml_duration", time.Since(mlStart)),
		)
		s.markError(ctx, task.ID, err.Error())
		return err
	}

	log.Info("ml analyze completed",
		zap.Int("notes_count", len(notes.Notes)),
		zap.Duration("ml_duration", time.Since(mlStart)),
	)

	seq := music.NewNoteSequence(len(notes.Notes))
	for i, note := range notes.Notes {
		name, octave := music.MidiToNote(int(note.MidiPitch))
		seq.Notes[i] = music.NoteEvent{
			Name:      name,
			Octave:    octave,
			MidiPitch: int(note.MidiPitch),
			StartTime: note.StartSeconds,
			EndTime:   note.DurationSeconds,
			Velocity:  note.Velocity,
		}
	}

	// --- TAB GENERATION
	genStart := time.Now()
	tab, err := s.tabProc.GenerateTab(&seq)
	if err != nil {
		log.Error("tab generation failed",
			zap.Error(err),
			zap.Duration("generation_duration", time.Since(genStart)),
		)
		s.markError(ctx, task.ID, err.Error())
		return err
	}

	log.Info("tab generated",
		zap.Duration("generation_duration", time.Since(genStart)),
	)

	// --- UPLOAD
	uploadStart := time.Now()
	uploadCtx, cancel := context.WithTimeout(ctx, s.cfg.StorageTimeout)
	defer cancel()

	err = s.uploadTab(uploadCtx, task.ID, tab, audioName)
	if err != nil {
		log.Error("failed to upload tab",
			zap.Error(err),
			zap.Duration("upload_duration", time.Since(uploadStart)),
		)
		s.markError(ctx, task.ID, err.Error())
		return err
	}

	log.Info("tab generation pipeline finished successfully",
		zap.Duration("total_duration", time.Since(start)),
	)

	return nil
}

func (s *TabGenService) startProcessing(ctx context.Context, task *models.TabGenTaskStartEvent) (bool, error) {
	ok, err := s.tgRepo.TryUpdateStatus(
		ctx,
		task.ID,
		[]models.Status{
			models.StatusPending,
			models.StatusWaitingForSeparation,
		},
		models.StatusProcessing,
		nil,
		nil,
	)
	if err != nil {
		s.log.Error("failed to update task status", zap.Error(err))
		return false, err
	}
	if !ok {
		s.log.Debug("task status not updated", zap.String("task_id", task.ID))
		return false, nil
	}

	return true, nil
}

func (s *TabGenService) analyze(ctx context.Context, audio *note_analyzer.AudioRequest) (*note_analyzer.NoteResponse, error) {
	waitStart := time.Now()

	if err := s.mlSem.Acquire(ctx, 1); err != nil {
		return nil, err
	}

	waitDuration := time.Since(waitStart)
	s.log.Debug("ml semaphore acquired",
		zap.Duration("wait_duration", waitDuration),
	)

	defer s.mlSem.Release(1)

	return s.mlAnalyzer.Analyze(ctx, audio)
}

func (s *TabGenService) loadAudio(ctx context.Context, tgID string) (string, *note_analyzer.AudioRequest, error) {
	tgTask, err := s.tgRepo.Get(ctx, tgID)
	if err != nil {
		s.log.Error("failed to get task", zap.Error(err), zap.String("task_id", tgID))
		return "", nil, err
	}

	var audioObjectName string

	if tgTask.AudioSepTaskID == nil {
		audioObjectName = tgTask.AudioFileObjectName()
	} else {
		asTask, err := s.asRepo.Get(ctx, *tgTask.AudioSepTaskID)
		if err != nil {
			s.log.Error("failed to get audio separation task", zap.Error(err), zap.String("task_id", *tgTask.AudioSepTaskID))
			return "", nil, err
		}
		audioObjectName = asTask.SeparatedPrefix() + "/other" + filepath.Ext(asTask.InputAudioName)
	}

	file, err := s.storage.GetFile(ctx, s.audioBucket, audioObjectName)
	if err != nil {
		s.log.Error("failed to get audio", zap.Error(err), zap.String("audio_object_name", audioObjectName))
		return "", nil, err
	}
	defer file.Close()

	_, audioName := path.Split(audioObjectName)
	data, err := io.ReadAll(file)
	if err != nil {
		s.log.Error("failed to read audio", zap.Error(err), zap.String("audio_object_name", audioObjectName))
		return "", nil, err
	}

	audio := &note_analyzer.AudioRequest{
		AudioData: &note_analyzer.AudioFileData{
			FileName:   audioName,
			AudioBytes: data,
		},
	}

	return audioName, audio, nil
}

func (s *TabGenService) uploadTab(ctx context.Context, tgID string, tab string, name string) error {
	task, err := s.tgRepo.Get(ctx, tgID)
	if err != nil {
		s.log.Error("failed to get task", zap.Error(err), zap.String("task_id", tgID))
		return err
	}

	name = tabFileName(name, ".txt")

	tabObjectName := task.ID + "/" + name
	r := strings.NewReader(tab)

	if err := s.storage.UploadFile(
		ctx,
		s.tabBucket,
		tabObjectName,
		r,
	); err != nil {
		s.log.Error("failed to put tab", zap.Error(err), zap.String("tab_object_name", tabObjectName))
		return err
	}

	task.ResultTabName = &name
	task.Status = models.StatusDone

	ok, err := s.tgRepo.TryUpdateStatus(
		ctx,
		task.ID,
		[]models.Status{
			models.StatusProcessing,
		},
		models.StatusDone,
		nil,
		&name,
	)
	if err != nil {
		s.log.Error("failed to update task", zap.Error(err), zap.String("task_id", tgID))
		return err
	}
	if !ok {
		s.log.Error("failed to update task", zap.Error(err), zap.String("task_id", tgID))
		return err
	}

	return nil
}

func (s *TabGenService) markError(ctx context.Context, id string, errMsg string) {
	log := s.log.With(zap.String("task_id", id))

	_, err := s.tgRepo.TryUpdateStatus(
		ctx,
		id,
		[]models.Status{
			models.StatusPending,
			models.StatusProcessing,
		},
		models.StatusFailed,
		&errMsg,
		nil,
	)

	if err != nil {
		log.Error("failed to mark task as failed",
			zap.Error(err),
		)
		return
	}

	log.Info("task marked as failed")
}

func tabFileName(audioFileName string, newExt string) string {
	ext := filepath.Ext(audioFileName)
	base := audioFileName[:len(audioFileName)-len(ext)]
	return base + newExt
}
