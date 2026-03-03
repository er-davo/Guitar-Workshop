//go:build integration

package repository_test

import (
	"context"
	"testing"
	"time"

	"api-gateway/internal/models"
	"api-gateway/internal/repository"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/stretchr/testify/require"
)

func TestGenTaskRepository_CreateAndGet(t *testing.T) {
	repo := repository.NewTabGenTaskRepository(db, trmpgx.DefaultCtxGetter)
	asRepo := repository.NewAudioSepTaskRepository(db, trmpgx.DefaultCtxGetter)
	ctx := context.Background()

	now := time.Now().UTC().Truncate(time.Second)

	asTask := models.AudioSepTask{
		Status:         models.StatusPending,
		InputAudioName: "/audio/input.wav",
		CreatedAt:      now,
	}

	err := asRepo.Create(ctx, &asTask)
	require.NoError(t, err)

	task := &models.TabGenTask{
		AudioSepTaskID: &asTask.ID,
		Status:         models.StatusPending,
		InputAudioName: "/audio/input.wav",
		ResultTabName:  nil,
		Saved:          false,
		Error:          nil,
		CreatedAt:      now,
	}

	// CREATE
	err = repo.Create(ctx, task)
	require.NoError(t, err)

	// GET
	fromDB, err := repo.Get(ctx, task.ID)
	require.NoError(t, err)

	require.NotEqual(t, "", fromDB.ID)
	require.Equal(t, task.AudioSepTaskID, fromDB.AudioSepTaskID)
	require.Equal(t, task.Status, fromDB.Status)
	require.Equal(t, task.InputAudioName, fromDB.InputAudioName)
	require.Equal(t, task.ResultTabName, fromDB.ResultTabName)
	require.Equal(t, task.Saved, fromDB.Saved)
	require.Equal(t, task.Error, fromDB.Error)
	require.WithinDuration(t, task.CreatedAt, fromDB.CreatedAt, time.Second)
}

func TestGenTaskRepository_Create_InvalidID(t *testing.T) {
	repo := repository.NewTabGenTaskRepository(db, trmpgx.DefaultCtxGetter)
	ctx := context.Background()

	task := &models.TabGenTask{
		ID: "",
	}

	err := repo.Create(ctx, task)
	require.Error(t, err)
}

func TestGenTaskRepository_MarkSaved(t *testing.T) {
	repo := repository.NewTabGenTaskRepository(db, trmpgx.DefaultCtxGetter)
	ctx := context.Background()

	task := &models.TabGenTask{
		Status:         models.StatusPending,
		InputAudioName: "/audio/input.wav",
		Saved:          false,
		CreatedAt:      time.Now().UTC().Truncate(time.Second),
	}

	err := repo.Create(ctx, task)
	require.NoError(t, err)

	fromDB, err := repo.Get(ctx, task.ID)
	require.NoError(t, err)
	require.False(t, fromDB.Saved)

	err = repo.MarkSaved(ctx, task.ID)
	require.NoError(t, err)

	updated, err := repo.Get(ctx, task.ID)
	require.NoError(t, err)
	require.True(t, updated.Saved)
}
