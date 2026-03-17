//go:build integration

package repository_test

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"audio-sep-task-service/internal/database"
	"audio-sep-task-service/internal/models"
	"audio-sep-task-service/internal/repository"

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

func TestAudioSepTaskRepository_CreateAndGet(t *testing.T) {
	repo := repository.NewAudioSepTaskRepository(db, trmpgx.DefaultCtxGetter)
	ctx := context.Background()

	task := &models.AudioSepTask{
		Status:           models.StatusPending,
		InputAudioName:   "test.wav",
		SeparatedDirName: "",
	}

	err := repo.Create(ctx, task)
	require.NoError(t, err)
	require.NotEmpty(t, task.ID)

	fromDB, err := repo.Get(ctx, task.ID)
	require.NoError(t, err)

	require.Equal(t, task.ID, fromDB.ID)
	require.Equal(t, task.Status, fromDB.Status)
	require.Equal(t, task.InputAudioName, fromDB.InputAudioName)
}

func TestAudioSepTaskRepository_TryUpdate_Success(t *testing.T) {
	repo := repository.NewAudioSepTaskRepository(db, trmpgx.DefaultCtxGetter)
	ctx := context.Background()

	task := &models.AudioSepTask{
		Status:         models.StatusPending,
		InputAudioName: "test.wav",
	}

	require.NoError(t, repo.Create(ctx, task))

	ok, err := repo.TryUpdate(
		ctx,
		task.ID,
		models.StatusPending,
		models.StatusProcessing,
		nil,
		nil,
	)

	require.NoError(t, err)
	require.True(t, ok)

	updated, err := repo.Get(ctx, task.ID)
	require.NoError(t, err)
	require.Equal(t, models.StatusProcessing, updated.Status)
}

func TestAudioSepTaskRepository_TryUpdate_WrongFromStatus(t *testing.T) {
	repo := repository.NewAudioSepTaskRepository(db, trmpgx.DefaultCtxGetter)
	ctx := context.Background()

	task := &models.AudioSepTask{
		Status:         models.StatusPending,
		InputAudioName: "test.wav",
	}

	require.NoError(t, repo.Create(ctx, task))

	ok, err := repo.TryUpdate(
		ctx,
		task.ID,
		models.StatusDone, // неправильный from
		models.StatusProcessing,
		nil,
		nil,
	)

	require.NoError(t, err)
	require.False(t, ok)
}

func TestAudioSepTaskRepository_TryUpdate_Done(t *testing.T) {
	repo := repository.NewAudioSepTaskRepository(db, trmpgx.DefaultCtxGetter)
	ctx := context.Background()

	task := &models.AudioSepTask{
		Status:         models.StatusProcessing,
		InputAudioName: "test.wav",
	}

	require.NoError(t, repo.Create(ctx, task))

	dir := "separated/123"

	ok, err := repo.TryUpdate(
		ctx,
		task.ID,
		models.StatusProcessing,
		models.StatusDone,
		&dir,
		nil,
	)

	require.NoError(t, err)
	require.True(t, ok)

	updated, err := repo.Get(ctx, task.ID)
	require.NoError(t, err)

	require.Equal(t, models.StatusDone, updated.Status)
	require.Equal(t, dir, updated.SeparatedDirName)
}

func TestAudioSepTaskRepository_TryUpdate_Error(t *testing.T) {
	repo := repository.NewAudioSepTaskRepository(db, trmpgx.DefaultCtxGetter)
	ctx := context.Background()

	task := &models.AudioSepTask{
		Status:         models.StatusProcessing,
		InputAudioName: "test.wav",
	}

	require.NoError(t, repo.Create(ctx, task))

	errMsg := "something failed"

	ok, err := repo.TryUpdate(
		ctx,
		task.ID,
		models.StatusProcessing,
		models.StatusFailed,
		nil,
		&errMsg,
	)

	require.NoError(t, err)
	require.True(t, ok)

	updated, err := repo.Get(ctx, task.ID)
	require.NoError(t, err)

	require.Equal(t, models.StatusFailed, updated.Status)
	require.Equal(t, &errMsg, updated.Error)
}
