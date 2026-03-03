package app

import (
	"api-gateway/internal/broker/produser"
	"api-gateway/internal/config"
	"api-gateway/internal/database"
	"api-gateway/internal/handlers"
	"api-gateway/internal/repository"
	"api-gateway/internal/service"
	"api-gateway/internal/storage"
	"context"
	"fmt"
	"net/http"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/segmentio/kafka-go"

	"go.uber.org/zap"
)

type App struct {
	log *zap.Logger
	cfg *config.Config

	db *pgxpool.Pool

	audioSepTaskWriter *kafka.Writer
	tabGenTaskWriter   *kafka.Writer

	server *echo.Echo
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

	tabRepo := repository.NewTabRepository(dbConn, trmpgx.DefaultCtxGetter)
	tabGenTaskRepo := repository.NewTabGenTaskRepository(dbConn, trmpgx.DefaultCtxGetter)
	audioSepTaskRepo := repository.NewAudioSepTaskRepository(dbConn, trmpgx.DefaultCtxGetter)

	s3Storage, err := storage.NewStorage(ctx, log)
	if err != nil {
		return nil, fmt.Errorf("error on initializing data storage: %w", err)
	}

	audioSepTaskWriter := kafka.NewWriter(kafka.WriterConfig{
		Brokers: cfg.Kafka.Brokers,
		Topic:   cfg.Kafka.Topics.AudioSeparation,
	})

	tabGenTaskWriter := kafka.NewWriter(kafka.WriterConfig{
		Brokers: cfg.Kafka.Brokers,
		Topic:   cfg.Kafka.Topics.TabGeneration,
	})

	tabGenTaskProducer := produser.NewTabGenTaskProduser(tabGenTaskWriter, log)

	audioSepTaskProducer := produser.NewAudioSepTaskProducer(audioSepTaskWriter, log)

	tabService := service.NewTabService(
		tabRepo,
		tabGenTaskRepo,
		manager.Must(trmpgx.NewDefaultFactory(dbConn)),
		s3Storage,
		cfg.Storage.TabBucket,
		cfg.Storage.ExpirationTime,
		retrier,
		log,
	)

	audioSepTaskService := service.NewAudioSepTaskService(
		audioSepTaskRepo,
		s3Storage,
		cfg.Storage.AudioBucket,
		cfg.Storage.ExpirationTime,
		audioSepTaskProducer,
		retrier,
		log,
	)

	tabGenTaskService := service.NewGenTaskService(
		tabGenTaskRepo,
		tabRepo,
		manager.Must(trmpgx.NewDefaultFactory(dbConn)),
		s3Storage,
		cfg.Storage.AudioBucket,
		cfg.Storage.TabBucket,
		cfg.Storage.ExpirationTime,
		audioSepTaskService,
		tabGenTaskProducer,
		retrier,
		log,
	)

	tabHandler := handlers.NewTabHandler(tabService, log)
	audioSepHandler := handlers.NewAudioSepTaskHandler(audioSepTaskService, log)
	tabGenHandler := handlers.NewGenTaskHandler(tabGenTaskService, log)

	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:3000"},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodOptions, http.MethodPut, http.MethodDelete},
		AllowHeaders: []string{"Content-Type"},
	}))

	e.Use(zapLogger(log))
	e.Use(middleware.Recover())

	tabHandler.RegisterRoutes(e)
	audioSepHandler.RegisterRoutes(e)
	tabGenHandler.RegisterRoutes(e)

	return &App{
		log:                log,
		cfg:                cfg,
		audioSepTaskWriter: audioSepTaskWriter,
		tabGenTaskWriter:   tabGenTaskWriter,
		db:                 dbConn,
		server:             e,
	}, nil
}

func (a *App) Run(ctx context.Context) error {
	go func() {
		if err := a.server.Start(":" + a.cfg.App.Port); err != nil && err != http.ErrServerClosed {
			a.log.Error("server start failed", zap.Error(err))
		}
	}()

	<-ctx.Done()
	return a.Shutdown()
}

func (a *App) Shutdown() error {
	ctxTimeout, cancelTimeout := context.WithTimeout(context.Background(), a.cfg.App.ShutdownTimeout)
	defer cancelTimeout()
	if err := a.server.Shutdown(ctxTimeout); err != nil {
		return fmt.Errorf("failed to shutdown echo server: %w", err)
	}

	a.db.Close()

	if err := a.audioSepTaskWriter.Close(); err != nil {
		return fmt.Errorf("failed to close audio separation task writer: %w", err)
	}

	if err := a.tabGenTaskWriter.Close(); err != nil {
		return fmt.Errorf("failed to close tab generation task writer: %w", err)
	}

	return nil
}
