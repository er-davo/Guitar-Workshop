package database

import (
	"context"

	"api-gateway/internal/logger"

	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

func Connect(dbURL string) *pgx.Conn {
	conn, err := pgx.Connect(context.Background(), dbURL)
	if err != nil {
		logger.Fatal("DB connection error", zap.Error(err))
	}
	return conn
	// config, _ := pgx.ParseConfig(dbURL)
	// config.StatementCacheCapacity = 0
	// conn, err := pgx.ConnectConfig(context.Background(), config)
	// if err != nil {
	// 	logger.Fatal("DB connection error", zap.Error(err))
	// }
	// return conn
}
