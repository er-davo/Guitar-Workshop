package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"api-gateway/internal/app"
	"api-gateway/internal/config"
	"api-gateway/internal/logger"

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

	app, err := app.New(ctx, cfg, log)
	if err != nil {
		log.Fatal("error on creating app", zap.Error(err))
	}

	if err := app.Run(ctx); err != nil {
		log.Fatal("server exited with error", zap.Error(err))
	}
}
