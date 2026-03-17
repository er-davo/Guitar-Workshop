package app

import (
	"context"
	"fmt"
	"net"

	"audio-sep-task-service/internal/broker/producer"
	"audio-sep-task-service/internal/config"
	"audio-sep-task-service/internal/database"
	"audio-sep-task-service/internal/repository"
	"audio-sep-task-service/internal/service"
	"audio-sep-task-service/internal/storage"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/er-davo/gwcontracts/audiosep"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type App struct {
	log *zap.Logger
	cfg *config.Config

	db *pgxpool.Pool

	audioSepTaskProducer *producer.AudioSepTaskProducer
	server               *grpc.Server
}

func New(ctx context.Context, cfg *config.Config, log *zap.Logger) (*App, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config is nil")
	}
	if log == nil {
		return nil, fmt.Errorf("logger is nil")
	}

	dbConn, err := database.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("error on initializing database connection: %w", err)
	}

	retrier := newRetrier(cfg.Retry)

	audioSepTaskRepo := repository.NewAudioSepTaskRepository(dbConn, trmpgx.DefaultCtxGetter)

	s3Storage, err := storage.NewStorage(ctx, log)
	if err != nil {
		return nil, fmt.Errorf("error on initializing data storage: %w", err)
	}

	audioSepTaskWriter := kafka.NewWriter(kafka.WriterConfig{
		Topic: cfg.Kafka.Topics.AudioSeparation,
	})

	audioSepTaskProducer := producer.NewAudioSepTaskProducer(audioSepTaskWriter, log)

	audioSepTaskService := service.NewAudioSepTaskService(
		audioSepTaskRepo,
		s3Storage,
		cfg.Storage.AudioBucket,
		cfg.Storage.ExpirationTime,
		audioSepTaskProducer,
		retrier,
		log,
	)

	server := grpc.NewServer()
	audiosep.RegisterAudioSepTaskServiceServer(server, audioSepTaskService)

	return &App{
		log:                  log,
		cfg:                  cfg,
		audioSepTaskProducer: audioSepTaskProducer,
		db:                   dbConn,
		server:               server,
	}, nil
}

func (a *App) Run(ctx context.Context) error {
	a.log.Info("starting audio separation task service")

	lis, err := net.Listen("tcp", ":"+a.cfg.App.Port)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	go func() {
		if err := a.server.Serve(lis); err != nil {
			a.log.Error("error on starting audio separation task service", zap.Error(err))
		}
	}()

	<-ctx.Done()

	return a.Shutdown()
}

func (a *App) Shutdown() error {
	a.server.GracefulStop()
	if err := a.audioSepTaskProducer.Close(); err != nil {
		a.log.Error("error on closing audio separation task producer", zap.Error(err))
		return err
	}
	a.db.Close()

	return nil
}
