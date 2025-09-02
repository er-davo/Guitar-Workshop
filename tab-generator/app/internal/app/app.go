package app

import (
	"context"
	"fmt"
	"net"
	"tabgen/internal/clients"
	"tabgen/internal/config"
	"tabgen/internal/handler"
	"tabgen/internal/proto/tab"
	"tabgen/internal/service"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type App struct {
	log *zap.Logger
	cfg *config.Config

	grpcServer *grpc.Server

	analyzerCliet clients.NoteAnalyzer
}

func New(ctx context.Context, cfg *config.Config, log *zap.Logger) (*App, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config is nil")
	}
	if log == nil {
		return nil, fmt.Errorf("logger is nil")
	}

	s := grpc.NewServer(
		grpc.MaxRecvMsgSize(100*1024*1024), // 100 MB
		grpc.MaxSendMsgSize(100*1024*1024), // 100 MB
	)

	analyzer, err := clients.NewNoteAnalyzerClient(
		cfg.Analyzer.Host+":"+cfg.Analyzer.Port,
		log,
	)
	if err != nil {
		return nil, fmt.Errorf("create analyzer client: %w", err)
	}

	tabService := service.NewTabService(analyzer, log)
	tabHandler := handler.NewTabHandler(tabService, log)

	tab.RegisterTabGenerateServer(s, tabHandler)

	return &App{
		log:           log,
		cfg:           cfg,
		grpcServer:    s,
		analyzerCliet: analyzer,
	}, nil
}

func (a *App) Run(ctx context.Context) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", a.cfg.App.Port))
	if err != nil {
		return fmt.Errorf("listen tcp: %w", err)
	}

	go func() {
		if err := a.grpcServer.Serve(lis); err != nil {
			a.log.Error("gRPC server stopped", zap.Error(err))
		}
	}()

	<-ctx.Done()
	a.log.Info("gracefully shutting down grpc server")

	return a.Shutdown()
}

func (a *App) Shutdown() error {
	a.grpcServer.GracefulStop()
	if err := a.analyzerCliet.Close(); err != nil {
		return fmt.Errorf("close analyzer client: %w", err)
	}
	return nil
}
