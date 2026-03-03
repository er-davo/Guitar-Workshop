package app

import (
	"context"
	"errors"
	"fmt"
	"orchestrator/internal/broker"
	"orchestrator/internal/broker/consumer"
	"orchestrator/internal/broker/producer"
	"orchestrator/internal/config"
	"orchestrator/internal/database"
	"orchestrator/internal/repository"
	"orchestrator/internal/service"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

type App struct {
	cfg *config.Config
	log *zap.Logger

	db *pgxpool.Pool

	audioSepTaskCompletedConsumer *consumer.AudioSepTaskCompletedConsumer
	tabGenTaskRequestConsumer     *consumer.TabGenTaskRequestedConsumer

	tabGenTaskStartProducer *producer.TabGenStartProducer
}

func New(ctx context.Context, cfg *config.Config, log *zap.Logger) (*App, error) {
	if cfg == nil {
		return nil, errors.New("config is nil")
	}
	if log == nil {
		return nil, errors.New("log is nil")
	}

	dbConn, err := database.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("error on initializing database connection: %w", err)
	}

	tabGenTaskRepo := repository.NewTabGenTaskRepository(dbConn, trmpgx.DefaultCtxGetter)
	audioSepTaskRepo := repository.NewAudioSepTaskRepository(dbConn, trmpgx.DefaultCtxGetter)

	audioSepService := service.NewAudioSepService(
		audioSepTaskRepo,
		tabGenTaskRepo,
		log,
	)

	tabGenService := service.NewTabGenService(tabGenTaskRepo, log)

	audioSepTaskCompletedReader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: cfg.Kafka.Brokers,
		Topic:   cfg.Kafka.Topics.AudioSeparationCompleted,
		GroupID: cfg.Kafka.GroupID,
	})

	// wait kafka to be ready
	broker.WaitKafkaConsumersGroupReadiness(
		cfg.Kafka.Brokers[0],
		cfg.Kafka.Topics.AudioSeparationCompleted,
		cfg.Kafka.Topics.TabGenerationRequested,
	)

	tabGenTaskRequestedReader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: cfg.Kafka.Brokers,
		Topic:   cfg.Kafka.Topics.TabGenerationRequested,
		GroupID: cfg.Kafka.GroupID,
	})

	tabGenTaskStartedWriter := kafka.NewWriter(kafka.WriterConfig{
		Brokers: cfg.Kafka.Brokers,
		Topic:   cfg.Kafka.Topics.TabGenerationStart,
	})

	tgStartProducer := producer.NewTabGenStartProducer(tabGenTaskStartedWriter, log)

	asComletedConsumer := consumer.NewAudioSepTaskCompletedConsumer(
		audioSepTaskCompletedReader,
		audioSepService,
		tgStartProducer,
		log,
	)
	tgRequestedConsumer := consumer.NewTabGenTaskRequstedConsumer(
		tabGenTaskRequestedReader,
		tabGenService,
		tgStartProducer,
		log,
	)

	return &App{
		cfg:                           cfg,
		log:                           log,
		db:                            dbConn,
		audioSepTaskCompletedConsumer: asComletedConsumer,
		tabGenTaskRequestConsumer:     tgRequestedConsumer,
		tabGenTaskStartProducer:       tgStartProducer,
	}, nil
}

func (a *App) Run(ctx context.Context) error {
	a.log.Info("starting app")

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		return a.audioSepTaskCompletedConsumer.Run(ctx)
	})
	g.Go(func() error {
		return a.tabGenTaskRequestConsumer.Run(ctx)
	})

	return g.Wait()
}

func (a *App) Shutdown(ctx context.Context) error {
	a.log.Info("shutting down app")

	if err := a.audioSepTaskCompletedConsumer.Close(); err != nil {
		a.log.Error("error closing audio separation task completed consumer", zap.Error(err))
		return err
	}
	if err := a.tabGenTaskRequestConsumer.Close(); err != nil {
		a.log.Error("error closing tab generation task requested consumer", zap.Error(err))
		return err
	}

	if err := a.tabGenTaskStartProducer.Close(); err != nil {
		a.log.Error("error closing tab generation task start producer", zap.Error(err))
		return err
	}

	a.db.Close()

	return nil
}
