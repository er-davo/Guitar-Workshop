//go:build integration

package repository_test

import (
	"context"
	"testing"
	"time"

	"tabgen/internal/models"
	"tabgen/internal/repository"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestAudioSepTaskRepository_CreateAndGet(t *testing.T) {
	repo := repository.NewAudioSepTaskRepository(db, trmpgx.DefaultCtxGetter)
	ctx := context.Background()

	now := time.Now().UTC().Truncate(time.Second)

	task := &models.AudioSepTask{
		ID:               uuid.NewString(),
		Status:           models.StatusPending,
		InputAudioName:   "/audio/input.wav",
		SeparatedDirName: "",
		Error:            nil,
		CreatedAt:        now,
	}

	// CREATE
	err := repo.Create(ctx, task)
	require.NoError(t, err)

	// GET
	fromDB, err := repo.Get(ctx, task.ID)
	require.NoError(t, err)

	require.Equal(t, task.ID, fromDB.ID)
	require.Equal(t, task.Status, fromDB.Status)
	require.Equal(t, task.InputAudioName, fromDB.InputAudioName)
	require.Equal(t, task.SeparatedDirName, fromDB.SeparatedDirName)
	require.Equal(t, task.Error, fromDB.Error)
	require.WithinDuration(t, task.CreatedAt, fromDB.CreatedAt, time.Second)
}

func TestAudioSepTaskRepository_Create_InvalidID(t *testing.T) {
	repo := repository.NewAudioSepTaskRepository(db, trmpgx.DefaultCtxGetter)
	ctx := context.Background()

	task := &models.AudioSepTask{
		ID: "",
	}

	err := repo.Create(ctx, task)
	require.Error(t, err)
}
