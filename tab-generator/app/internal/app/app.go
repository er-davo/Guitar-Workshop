package app

import (
	"context"
	"fmt"

	"tabgen/internal/broker"
	"tabgen/internal/broker/consumer"
	"tabgen/internal/clients"
	"tabgen/internal/config"
	"tabgen/internal/database"
	"tabgen/internal/processor"
	"tabgen/internal/repository"
	"tabgen/internal/service"
	"tabgen/internal/storage"
	"tabgen/internal/worker"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type TabGenStartConsumer interface {
	Run(ctx context.Context)
	Close() error
}

type App struct {
	log *zap.Logger
	cfg *config.Config

	tgStartConsumer TabGenStartConsumer

	analyzerCliet clients.NoteAnalyzer
}

func New(ctx context.Context, cfg *config.Config, log *zap.Logger) (*App, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config is nil")
	}
	if log == nil {
		return nil, fmt.Errorf("logger is nil")
	}

	analyzer, err := clients.NewNoteAnalyzerClient(
		cfg.Analyzer.Host+":"+cfg.Analyzer.Port,
		log,
	)
	if err != nil {
		return nil, fmt.Errorf("create analyzer client: %w", err)
	}

	db, err := database.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("connect to database: %w", err)
	}

	tgRepo := repository.NewTabGenTaskRepository(db, trmpgx.DefaultCtxGetter)
	asRepo := repository.NewAudioSepTaskRepository(db, trmpgx.DefaultCtxGetter)

	dataStorage, err := storage.NewStorage(ctx, log)
	if err != nil {
		return nil, fmt.Errorf("error on initializing data storage: %w", err)
	}

	tabGenService := service.NewTabGenStartService(
		&service.TabGenServiceConfig{
			TaskTimeout:    cfg.App.Service.Config.TaskTimeout,
			DBTimeout:      cfg.App.Service.Config.DBTimeout,
			StorageTimeout: cfg.App.Service.Config.StorageTimeout,
			MLTimeout:      cfg.App.Service.Config.MLTimeout,
		},
		tgRepo,
		asRepo,
		cfg.App.Service.MaxMLRequests,
		analyzer,
		dataStorage,
		cfg.Storage.AudioBucket,
		cfg.Storage.TabBucket,
		*processor.NewTabProcessor(),
		log,
	)

	workerPool := worker.NewPool(cfg.App.MaxWorkers)

	// wait kafka to be ready
	broker.WaitKafkaConsumersGroupReadiness(
		cfg.Kafka.Brokers[0],
		cfg.Kafka.Topics.TabGenerationStart,
	)

	consumer := consumer.NewTabGenTaskStartConsumer(kafka.NewReader(
		kafka.ReaderConfig{
			Brokers: cfg.Kafka.Brokers,
			Topic:   cfg.Kafka.Topics.TabGenerationStart,
			GroupID: cfg.Kafka.GroupID,
		},
	),
		tabGenService,
		workerPool,
		log,
	)

	return &App{
		log:             log,
		cfg:             cfg,
		tgStartConsumer: consumer,
		analyzerCliet:   analyzer,
	}, nil
}

func (a *App) Run(ctx context.Context) error {

	go a.tgStartConsumer.Run(ctx)

	<-ctx.Done()
	a.log.Info("gracefully shutting down grpc server")

	return a.Shutdown()
}

func (a *App) Shutdown() error {
	if err := a.analyzerCliet.Close(); err != nil {
		return fmt.Errorf("close analyzer client: %w", err)
	}
	if err := a.tgStartConsumer.Close(); err != nil {
		return fmt.Errorf("close tab gen start consumer: %w", err)
	}
	return nil
}
