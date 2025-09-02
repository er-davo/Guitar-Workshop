package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"tabgen/internal/app"
	"tabgen/internal/config"
	"tabgen/internal/logger"

	"go.uber.org/zap"
)

func main() {
	log := logger.NewLogger()

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH environment variable not set")
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		log.Fatal("error on loading config", zap.Error(err))
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	application, err := app.New(ctx, cfg, log)
	if err != nil {
		log.Fatal("error on creating app", zap.Error(err))
	}

	log.Info("statring tab generator service", zap.String("port", cfg.App.Port))

	if err := application.Run(ctx); err != nil {
		log.Error("server exited with error", zap.Error(err))
	}
}
