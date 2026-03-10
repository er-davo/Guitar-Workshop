package app

import (
	"context"
	"fmt"
	"net"

	"tab-service/internal/config"
	"tab-service/internal/database"
	"tab-service/internal/repository"
	"tab-service/internal/service"
	"tab-service/internal/storage"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/er-davo/gwcontracts/tab"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type App struct {
	log *zap.Logger
	cfg *config.Config

	db *pgxpool.Pool

	server *grpc.Server
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

	s3Storage, err := storage.NewStorage(ctx, log)
	if err != nil {
		return nil, fmt.Errorf("error on initializing data storage: %w", err)
	}

	tabService := service.NewTabService(
		tabRepo,
		manager.Must(trmpgx.NewDefaultFactory(dbConn)),
		s3Storage,
		cfg.Storage.TabBucket,
		cfg.Storage.ExpirationTime,
		retrier,
		log,
	)

	srv := grpc.NewServer()

	tab.RegisterTabServiceServer(srv, tabService)

	return &App{
		log:    log,
		cfg:    cfg,
		db:     dbConn,
		server: srv,
	}, nil
}

func (a *App) Run(ctx context.Context) error {
	a.log.Info("starting app")

	lis, err := net.Listen("tcp", ":"+a.cfg.App.Port)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	go func() {
		if err := a.server.Serve(lis); err != nil {
			a.log.Error("grpc server stopped", zap.Error(err))
		}
	}()

	<-ctx.Done()

	return a.Shutdown(ctx)
}

func (a *App) Shutdown(ctx context.Context) error {
	a.log.Info("shutting down app")

	a.server.GracefulStop()
	a.db.Close()

	return nil
}
