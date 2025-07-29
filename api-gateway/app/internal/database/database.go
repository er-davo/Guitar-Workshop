package database

import (
	"api-gateway/internal/logger"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

func Connect(dbURL string) *pgxpool.Pool {
	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		logger.Fatal("DB connection error", zap.Error(err))
	}
	return pool
}
