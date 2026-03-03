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

func TestAudioSepService_HandleAudioSepTaskCompletedEvent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAudioRepo := mocks.NewMockAudioSepTaskRepository(ctrl)
	mockTabRepo := mocks.NewMockTabGenTaskRepository(ctrl)

	logger := zap.NewNop()

	svc := service.NewAudioSepService(mockAudioRepo, mockTabRepo, logger)
	ctx := context.Background()

	tests := []struct {
		name         string
		audioSepTask *models.AudioSepTask
		tabGenTask   *models.TabGenTask
		audioSepErr  error
		tabGenErr    error
		expectedID   string
		expectedErr  error
	}{
		{
			name: "normal flow - tab gen exists",
			audioSepTask: &models.AudioSepTask{
				ID: "audio1",
			},
			tabGenTask: &models.TabGenTask{
				ID:             "tab1",
				AudioSepTaskID: strPtr("audio1"),
			},
			expectedID:  "tab1",
			expectedErr: nil,
		},
		{
			name:         "audio sep not found",
			audioSepTask: nil,
			audioSepErr:  errors.New("db error"),
			expectedID:   "",
			expectedErr:  errors.New("db error"),
		},
		{
			name: "tab gen not found",
			audioSepTask: &models.AudioSepTask{
				ID: "audio2",
			},
			tabGenTask:  nil,
			tabGenErr:   errors.New("not found"),
			expectedID:  "",
			expectedErr: errors.New("not found"),
		},
		// {
		// 	name: "tab gen nil - no error",
		// 	audioSepTask: &models.AudioSepTask{
		// 		ID: "audio3",
		// 	},
		// 	tabGenTask:  nil,
		// 	tabGenErr:   service.ErrTabGenNotFound, // если возвращаем отдельный ErrNotFound
		// 	expectedID:  "",
		// 	expectedErr: service.ErrTabGenNotFound,
		// },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event := &models.AudioSepTaskCompletedEvent{ID: ""}
			if tt.audioSepTask != nil {
				event.ID = tt.audioSepTask.ID
			}

			mockAudioRepo.EXPECT().
				Get(ctx, gomock.Any()).
				Return(tt.audioSepTask, tt.audioSepErr).
				Times(1)

			if tt.audioSepTask != nil && tt.audioSepErr == nil {
				mockTabRepo.EXPECT().
					GetByAudioSepTaskID(ctx, tt.audioSepTask.ID).
					Return(tt.tabGenTask, tt.tabGenErr).
					Times(1)
			}

			id, err := svc.HandleAudioSepTaskCompletedEvent(ctx, event)
			if tt.expectedErr != nil {
				require.EqualError(t, err, tt.expectedErr.Error())
				require.Empty(t, id)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedID, id)
			}
		})
	}
}
