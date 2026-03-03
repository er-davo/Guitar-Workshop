//go:build integration

package repository_test

import (
	"context"
	"testing"
	"time"

	"orchestrator/internal/models"
	"orchestrator/internal/repository"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestGenTaskRepository_CreateAndGet(t *testing.T) {
	repo := repository.NewTabGenTaskRepository(db, trmpgx.DefaultCtxGetter)
	ctx := context.Background()

	now := time.Now().UTC().Truncate(time.Second)

	task := &models.TabGenTask{
		ID:             uuid.NewString(),
		AudioSepTaskID: nil,
		Status:         models.StatusPending,
		InputAudioName: "/audio/input.wav",
		ResultTabName:  nil,
		Saved:          false,
		Error:          nil,
		CreatedAt:      now,
	}

	// CREATE
	err := repo.Create(ctx, task)
	require.NoError(t, err)

	// GET
	fromDB, err := repo.Get(ctx, task.ID)
	require.NoError(t, err)

	require.Equal(t, task.ID, fromDB.ID)
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
		ID:             uuid.NewString(),
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
func TestGenTaskRepository_TryUpdateStatus(t *testing.T) {
	repo := repository.NewTabGenTaskRepository(db, trmpgx.DefaultCtxGetter)
	ctx := context.Background()

	tests := []struct {
		name          string
		initialStatus models.Status
		fromStatus    models.Status
		toStatus      models.Status
		errMsg        *string
		resultName    *string

		expectUpdated bool
		expectStatus  models.Status
		expectErrMsg  *string
		expectResult  *string
	}{
		{
			name:          "success without optional fields",
			initialStatus: models.StatusPending,
			fromStatus:    models.StatusPending,
			toStatus:      models.StatusProcessing,
			expectUpdated: true,
			expectStatus:  models.StatusProcessing,
		},
		{
			name:          "with error message",
			initialStatus: models.StatusProcessing,
			fromStatus:    models.StatusProcessing,
			toStatus:      models.StatusFailed,
			errMsg:        strPtr("generation failed"),
			expectUpdated: true,
			expectStatus:  models.StatusFailed,
			expectErrMsg:  strPtr("generation failed"),
		},
		{
			name:          "with result tab name",
			initialStatus: models.StatusProcessing,
			fromStatus:    models.StatusProcessing,
			toStatus:      models.StatusDone,
			resultName:    strPtr("result.tab"),
			expectUpdated: true,
			expectStatus:  models.StatusDone,
			expectResult:  strPtr("result.tab"),
		},
		{
			name:          "wrong from status",
			initialStatus: models.StatusPending,
			fromStatus:    models.StatusProcessing,
			toStatus:      models.StatusDone,
			expectUpdated: false,
			expectStatus:  models.StatusPending,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task := &models.TabGenTask{
				ID:             uuid.NewString(),
				Status:         tt.initialStatus,
				InputAudioName: "/audio/input.wav",
				CreatedAt:      time.Now().UTC().Truncate(time.Second),
			}

			require.NoError(t, repo.Create(ctx, task))

			updated, err := repo.TryUpdateStatus(
				ctx,
				task.ID,
				tt.fromStatus,
				tt.toStatus,
				tt.errMsg,
				tt.resultName,
			)

			require.NoError(t, err)
			require.Equal(t, tt.expectUpdated, updated)

			fromDB, dbErr := repo.Get(ctx, task.ID)
			require.NoError(t, dbErr)

			require.Equal(t, tt.expectStatus, fromDB.Status)

			if tt.expectErrMsg != nil {
				require.NotNil(t, fromDB.Error)
				require.Equal(t, *tt.expectErrMsg, *fromDB.Error)
			} else {
				require.Nil(t, fromDB.Error)
			}

			if tt.expectResult != nil {
				require.NotNil(t, fromDB.ResultTabName)
				require.Equal(t, *tt.expectResult, *fromDB.ResultTabName)
			} else {
				require.Nil(t, fromDB.ResultTabName)
			}
		})
	}
}

func strPtr(s string) *string {
	return &s
}
