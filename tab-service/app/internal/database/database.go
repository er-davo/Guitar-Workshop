package database

import (
	"context"

	"github.com/golang-migrate/migrate/v4"
	"github.com/jackc/pgx/v5/pgxpool"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func Connect(ctx context.Context, dbURL string) (*pgxpool.Pool, error) {
	conn, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		return nil, err
	}
	if err := conn.Ping(ctx); err != nil {
		return nil, err
	}
	return conn, nil
}

func Migrate(migrationDir string, dbURL string) error {
	m, err := migrate.New("file://"+migrationDir, dbURL)
	if err != nil {
		return err
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}
