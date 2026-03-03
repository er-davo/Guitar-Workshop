package service_test

import (
	"context"
	"errors"
	"testing"

	"orchestrator/internal/mocks"
	"orchestrator/internal/models"
	"orchestrator/internal/service"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

func TestTabGenService_HandleTabGenEvent(t *testing.T) {
	tests := []struct {
		name           string
		task           *models.TabGenTask
		getErr         error
		tryUpdateOk    bool
		tryUpdateErr   error
		expectedBool   bool
		expectedErr    error
		expectedStatus models.Status
	}{
		{
			name: "no audio sep -> processing",
			task: &models.TabGenTask{
				ID:             "1",
				Status:         models.StatusPending,
				AudioSepTaskID: nil,
			},
			tryUpdateOk:    true,
			expectedBool:   true,
			expectedStatus: models.StatusProcessing,
		},
		{
			name: "with audio sep -> waiting for separation",
			task: &models.TabGenTask{
				ID:             "2",
				Status:         models.StatusPending,
				AudioSepTaskID: strPtr("audio1"),
			},
			tryUpdateOk:    true,
			expectedBool:   true,
			expectedStatus: models.StatusWaitingForSeparation,
		},
		{
			name: "optimistic lock - no update",
			task: &models.TabGenTask{
				ID:             "3",
				Status:         models.StatusPending,
				AudioSepTaskID: nil,
			},
			tryUpdateOk:    false,
			expectedBool:   false,
			expectedStatus: models.StatusProcessing,
		},
		{
			name:        "Get error",
			getErr:      errors.New("db error"),
			expectedErr: errors.New("db error"),
		},
		{
			name: "TryUpdateStatus error",
			task: &models.TabGenTask{
				ID:             "4",
				Status:         models.StatusPending,
				AudioSepTaskID: nil,
			},
			tryUpdateErr:   errors.New("update failed"),
			expectedErr:    errors.New("update failed"),
			expectedStatus: models.StatusProcessing,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mocks.NewMockTabGenTaskRepository(ctrl)
			logger := zap.NewNop()

			svc := service.NewTabGenService(mockRepo, logger)
			ctx := context.Background()

			event := &models.TabGenTaskRequestedEvent{}
			if tt.task != nil {
				event.ID = tt.task.ID
			}

			mockRepo.EXPECT().
				Get(ctx, event.ID).
				Return(tt.task, tt.getErr).
				Times(1)

			if tt.task != nil && tt.getErr == nil {
				mockRepo.EXPECT().
					TryUpdateStatus(
						ctx,
						tt.task.ID,
						models.StatusPending,
						tt.expectedStatus,
						nil,
						nil,
					).
					Return(tt.tryUpdateOk, tt.tryUpdateErr).
					Times(1)
			}

			ok, err := svc.HandleTabGenEvent(ctx, event)

			if tt.expectedErr != nil {
				require.EqualError(t, err, tt.expectedErr.Error())
				require.False(t, ok)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.expectedBool, ok)
		})
	}
}

func strPtr(s string) *string { return &s }
