package service_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/url"
	"testing"
	"time"

	"api-gateway/internal/mocks"
	"api-gateway/internal/models"
	"api-gateway/internal/service"

	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/er-davo/retry"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

func setupAudioSepService(t *testing.T) (
	*service.AudioSepTaskService,
	*mocks.MockAudioSepTaskRepository,
	*mocks.MockStorage,
	*mocks.MockAudioSepTaskProducer,
) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	repo := mocks.NewMockAudioSepTaskRepository(ctrl)
	storage := mocks.NewMockStorage(ctrl)
	producer := mocks.NewMockAudioSepTaskProducer(ctrl)

	svc := service.NewAudioSepTaskService(
		repo,
		storage,
		"audio-bucket",
		time.Minute,
		producer,
		retry.NoRetry(),
		zap.NewNop(),
	)

	return svc, repo, storage, producer
}

func TestAudioSepService_PostAudioSepTask(t *testing.T) {
	ctx := context.Background()
	fileContent := []byte("audio content")

	t.Run("success", func(t *testing.T) {
		svc, repo, storage, producer := setupAudioSepService(t)

		repo.EXPECT().
			Create(ctx, gomock.Any()).
			DoAndReturn(func(_ context.Context, tsk *models.AudioSepTask) error {
				tsk.ID = "task1"
				return nil
			})

		storage.EXPECT().
			UploadFile(ctx, "audio-bucket", gomock.Any(), gomock.Any()).
			DoAndReturn(func(_ context.Context, _ string, _ string, r io.Reader) error {
				buf := new(bytes.Buffer)
				_, _ = io.Copy(buf, r)
				require.Equal(t, string(fileContent), buf.String())
				return nil
			})

		producer.EXPECT().
			Produce(ctx, gomock.Any()).
			Return(nil)

		result, err := svc.PostAudioSepTask(ctx, models.AudioSepTaskData{
			AudioFileName: "song.wav",
			Data:          fileContent,
			Size:          int64(len(fileContent)),
		})

		require.NoError(t, err)
		require.Equal(t, "task1", result.ID)
	})

	t.Run("repo create error", func(t *testing.T) {
		svc, repo, _, _ := setupAudioSepService(t)

		repo.EXPECT().
			Create(ctx, gomock.Any()).
			Return(errors.New("db fail"))

		_, err := svc.PostAudioSepTask(ctx, models.AudioSepTaskData{
			AudioFileName: "song.wav",
			Data:          fileContent,
			Size:          int64(len(fileContent)),
		})

		require.ErrorContains(t, err, "db fail")
	})

	t.Run("upload error", func(t *testing.T) {
		svc, repo, storage, _ := setupAudioSepService(t)

		repo.EXPECT().Create(ctx, gomock.Any()).Return(nil)

		storage.EXPECT().
			UploadFile(ctx, "audio-bucket", gomock.Any(), gomock.Any()).
			Return(errors.New("upload fail"))

		_, err := svc.PostAudioSepTask(ctx, models.AudioSepTaskData{
			AudioFileName: "song.wav",
			Data:          fileContent,
			Size:          int64(len(fileContent)),
		})

		require.ErrorContains(t, err, "upload fail")
	})

	t.Run("produce error", func(t *testing.T) {
		svc, repo, storage, producer := setupAudioSepService(t)

		repo.EXPECT().Create(ctx, gomock.Any()).Return(nil)

		storage.EXPECT().
			UploadFile(ctx, "audio-bucket", gomock.Any(), gomock.Any()).
			Return(nil)

		producer.EXPECT().
			Produce(ctx, gomock.Any()).
			Return(errors.New("kafka fail"))

		_, err := svc.PostAudioSepTask(ctx, models.AudioSepTaskData{
			AudioFileName: "song.wav",
			Data:          fileContent,
			Size:          int64(len(fileContent)),
		})

		require.ErrorContains(t, err, "kafka fail")
	})
}

func TestAudioSepService_Get(t *testing.T) {
	ctx := context.Background()

	task := &models.AudioSepTask{
		ID:               "task1",
		Status:           models.StatusDone,
		SeparatedDirName: "sep",
	}

	t.Run("success", func(t *testing.T) {
		svc, repo, storage, _ := setupAudioSepService(t)

		repo.EXPECT().
			Get(ctx, "task1").
			Return(task, nil)

		storage.EXPECT().
			ListFilesByPrefix(ctx, "audio-bucket", gomock.Any()).
			Return([]types.Object{
				{Key: ptrString("task1/sep/track1.wav")},
				{Key: ptrString("task1/sep/track2.wav")},
			}, nil)

		storage.EXPECT().
			PresignedGet(ctx, "audio-bucket", "task1/sep/track1.wav", gomock.Any()).
			Return(ptrURL("http://track1"), nil)

		storage.EXPECT().
			PresignedGet(ctx, "audio-bucket", "task1/sep/track2.wav", gomock.Any()).
			Return(ptrURL("http://track2"), nil)

		result, err := svc.Get(ctx, "task1")

		require.NoError(t, err)
		require.Equal(t, 2, len(result.SeparatedAudioSignedURLs))
		require.Contains(t, result.SeparatedAudioSignedURLs, "track1.wav")
		require.Contains(t, result.SeparatedAudioSignedURLs, "track2.wav")
	})

	t.Run("repo get error", func(t *testing.T) {
		svc, repo, _, _ := setupAudioSepService(t)

		repo.EXPECT().
			Get(ctx, "task1").
			Return(nil, errors.New("db fail"))

		_, err := svc.Get(ctx, "task1")
		require.ErrorContains(t, err, "db fail")
	})
}

func ptrURL(s string) *url.URL {
	u, _ := url.Parse(s)
	return u
}

func ptrString(s string) *string {
	return &s
}
