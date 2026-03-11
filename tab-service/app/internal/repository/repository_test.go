//go:build integration

package repository_test

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"tab-service/internal/database"
	"tab-service/internal/models"
	"tab-service/internal/repository"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
	tc "github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

var db *pgxpool.Pool

func TestMain(m *testing.M) {
	ctx := context.Background()

	pgContainer, err := postgres.Run(ctx,
		"postgres:15.3-alpine",
		postgres.WithDatabase("test_db"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		tc.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(10*time.Second)),
	)
	if err != nil {
		log.Fatal(err)
	}

	dsn, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	if err := database.Migrate("../../../../migrations", dsn); err != nil {
		log.Fatal(err)
	}

	db, err = database.Connect(ctx, dsn)
	if err != nil {
		log.Fatal(err)
	}

	code := m.Run()

	_ = pgContainer.Terminate(ctx)
	os.Exit(code)
}

func TestTabRepository_CRUD(t *testing.T) {
	repo := repository.NewTabRepository(db, trmpgx.DefaultCtxGetter)
	ctx := context.Background()

	// CREATE
	tab := &models.Tab{
		Name: "Metallica Nothing Else Matters",
		Path: "/tabs/nem.gp",
	}

	err := repo.Create(ctx, tab)
	require.NoError(t, err)
	require.NotEmpty(t, tab.ID)

	// GET BY ID
	fromDB, err := repo.Get(ctx, tab.ID)
	require.NoError(t, err)
	require.Equal(t, tab.Name, fromDB.Name)
	require.Equal(t, tab.Path, fromDB.Path)

	// CREATE SECOND TAB (чтобы проверить поиск)
	tab2 := &models.Tab{
		Name: "Metallica Master of Puppets",
		Path: "/tabs/mop.gp",
	}

	err = repo.Create(ctx, tab2)
	require.NoError(t, err)

	// SEARCH (FTS)
	list, err := repo.Search(ctx, "metallica nothing", 10, 0)
	require.NoError(t, err)
	require.Len(t, list, 1)
	require.Equal(t, tab.ID, list[0].ID)

	// SEARCH BY SINGLE WORD
	list, err = repo.Search(ctx, "metallica", 10, 0)
	require.NoError(t, err)
	require.Len(t, list, 2)

	// PAGINATION
	list, err = repo.Search(ctx, "metallica", 1, 0)
	require.NoError(t, err)
	require.Len(t, list, 1)

	list2, err := repo.Search(ctx, "metallica", 1, 1)
	require.NoError(t, err)
	require.Len(t, list2, 1)
	require.NotEqual(t, list[0].ID, list2[0].ID)

	// DELETE
	err = repo.Delete(ctx, tab.ID)
	require.NoError(t, err)

	// GET AFTER DELETE
	_, err = repo.Get(ctx, tab.ID)
	require.Error(t, err)
}
