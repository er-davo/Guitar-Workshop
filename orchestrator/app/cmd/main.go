package main

import (
	"context"
	"orchestrator/internal/app"
	"orchestrator/internal/config"
	"orchestrator/internal/logger"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
)

func main() {
	log := logger.NewLogger()
	defer log.Sync()

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH environment variable not set")
	}
	cfg, err := config.Load(configPath)
	if err != nil {
		log.Fatal("error on loading config: " + err.Error())
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	orchestratorApp, err := app.New(ctx, cfg, log)
	if err != nil {
		log.Fatal("error on creating app", zap.Error(err))
	}

	go func() {
		if err := orchestratorApp.Run(ctx); err != nil {
			log.Fatal("orchestrator application exited with error", zap.Error(err))
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.App.ShutdownTimeout)
	defer cancel()

	if err := orchestratorApp.Shutdown(shutdownCtx); err != nil {
		log.Error("error on shutting down app", zap.Error(err))
	}

}
