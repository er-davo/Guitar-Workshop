package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"api-gateway/internal/mocks"
	"api-gateway/internal/models"
	"api-gateway/internal/service"

	"github.com/er-davo/retry"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

func newTestService(t *testing.T) (
	*service.GenTaskService,
	*mocks.MockTabGenTaskRepository,
	*mocks.MockTabRepository,
	*mocks.MockStorage,
	*mocks.MockAudioSepTaskPoster,
	*mocks.MockTabGenTaskProducer,
) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	tabGenRepo := mocks.NewMockTabGenTaskRepository(ctrl)
	tabRepo := mocks.NewMockTabRepository(ctrl)
	storage := mocks.NewMockStorage(ctrl)
	audioSep := mocks.NewMockAudioSepTaskPoster(ctrl)
	producer := mocks.NewMockTabGenTaskProducer(ctrl)

	svc := service.NewGenTaskService(
		tabGenRepo,
		tabRepo,
		service.TxManagerStub{},
		storage,
		"audio-bucket",
		"tab-bucket",
		time.Minute,
		audioSep,
		producer,
		retry.NoRetry(),
		zap.NewNop(),
	)

	return svc, tabGenRepo, tabRepo, storage, audioSep, producer
}

func TestGenTaskService_PostGenTask(t *testing.T) {
	tests := []struct {
		name  string
		data  models.PostGenTaskData
		setup func(
			repo *mocks.MockTabGenTaskRepository,
			storage *mocks.MockStorage,
			audioSep *mocks.MockAudioSepTaskPoster,
			producer *mocks.MockTabGenTaskProducer,
		)
		wantErr   bool
		checkTask func(t *testing.T, task *models.TabGenTask)
	}{
		{
			name: "success without separation",
			data: models.PostGenTaskData{
				AudioFileName: "song.wav",
				Data:          []byte("audio"),
				Size:          5,
				Separation:    false,
			},
			setup: func(repo *mocks.MockTabGenTaskRepository, storage *mocks.MockStorage, audioSep *mocks.MockAudioSepTaskPoster, producer *mocks.MockTabGenTaskProducer) {
				audioSep.EXPECT().PostAudioSepTask(gomock.Any(), gomock.Any()).Times(0)

				storage.EXPECT().
					UploadFile(gomock.Any(), "audio-bucket", gomock.Any(), gomock.Any()).
					Return(nil)

				repo.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					DoAndReturn(func(_ context.Context, task *models.TabGenTask) error {
						task.ID = "task-id"
						return nil
					})

				producer.EXPECT().
					Produce(gomock.Any(), gomock.AssignableToTypeOf(&models.StartTabGenTask{})).
					Return(nil)
			},
			checkTask: func(t *testing.T, task *models.TabGenTask) {
				require.Equal(t, models.StatusCreated, task.Status)
				require.Nil(t, task.AudioSepTaskID)
			},
		},
		{
			name: "success with separation",
			data: models.PostGenTaskData{
				AudioFileName: "song.wav",
				Data:          []byte("audio"),
				Size:          5,
				Separation:    true,
			},
			setup: func(repo *mocks.MockTabGenTaskRepository, storage *mocks.MockStorage, audioSep *mocks.MockAudioSepTaskPoster, producer *mocks.MockTabGenTaskProducer) {
				audioSep.EXPECT().
					PostAudioSepTask(gomock.Any(), gomock.Any()).
					Return(&models.AudioSepTask{ID: "sep-id"}, nil)

				storage.EXPECT().
					UploadFile(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Times(0)

				repo.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					DoAndReturn(func(_ context.Context, task *models.TabGenTask) error {
						task.ID = "task-id"
						return nil
					})

				repo.EXPECT().
					SetAudioSeparationID(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil)

				producer.EXPECT().
					Produce(gomock.Any(), gomock.Any()).
					Return(nil)
			},
			checkTask: func(t *testing.T, task *models.TabGenTask) {
				require.NotNil(t, task.AudioSepTaskID)
				require.Equal(t, "sep-id", *task.AudioSepTaskID)
			},
		},
		{
			name: "upload error",
			data: models.PostGenTaskData{
				AudioFileName: "song.wav",
				Data:          []byte("audio"),
				Size:          5,
				Separation:    false,
			},
			setup: func(repo *mocks.MockTabGenTaskRepository, storage *mocks.MockStorage, _ *mocks.MockAudioSepTaskPoster, _ *mocks.MockTabGenTaskProducer) {
				repo.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(nil)

				storage.EXPECT().
					UploadFile(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(errors.New("upload failed"))
			},
			wantErr: true,
		},
		{
			name: "repo error",
			data: models.PostGenTaskData{
				AudioFileName: "song.wav",
				Data:          []byte("audio"),
				Size:          5,
				Separation:    false,
			},
			setup: func(repo *mocks.MockTabGenTaskRepository, _ *mocks.MockStorage, _ *mocks.MockAudioSepTaskPoster, _ *mocks.MockTabGenTaskProducer) {
				repo.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(errors.New("db error"))
			},
			wantErr: true,
		},
		{
			name: "producer error",
			data: models.PostGenTaskData{
				AudioFileName: "song.wav",
				Data:          []byte("audio"),
				Size:          5,
				Separation:    false,
			},
			setup: func(repo *mocks.MockTabGenTaskRepository, storage *mocks.MockStorage, _ *mocks.MockAudioSepTaskPoster, producer *mocks.MockTabGenTaskProducer) {
				storage.EXPECT().
					UploadFile(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil)

				repo.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					DoAndReturn(func(_ context.Context, task *models.TabGenTask) error {
						task.ID = "task-id"
						return nil
					})

				producer.EXPECT().
					Produce(gomock.Any(), gomock.Any()).
					Return(errors.New("kafka down"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			svc, repo, _, storage, audioSep, producer := newTestService(t)

			if tt.setup != nil {
				tt.setup(repo, storage, audioSep, producer)
			}

			task, err := svc.PostGenTask(ctx, tt.data)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, task)

			if tt.checkTask != nil {
				tt.checkTask(t, task)
			}
		})
	}
}
