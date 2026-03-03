package service_test

import (
	"bytes"
	"context"
	"errors"
	"io"
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

func setupService(t *testing.T) (
	*service.TabService,
	*mocks.MockTabRepository,
	*mocks.MockTabGenTaskRepository,
	*mocks.MockStorage,
) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	tabRepo := mocks.NewMockTabRepository(ctrl)
	genTaskRepo := mocks.NewMockTabGenTaskRepository(ctrl)
	storage := mocks.NewMockStorage(ctrl)
	logger := zap.NewNop()

	svc := service.NewTabService(
		tabRepo,
		genTaskRepo,
		service.TxManagerStub{},
		storage,
		"tabs-bucket",
		time.Minute,
		retry.NoRetry(),
		logger,
	)

	return svc, tabRepo, genTaskRepo, storage
}

func TestTabService_CreateTab(t *testing.T) {
	ctx := context.Background()
	tabBody := "TAB CONTENT"

	t.Run("success", func(t *testing.T) {
		svc, tabRepo, _, storage := setupService(t)

		storage.EXPECT().
			UploadFile(ctx, "tabs-bucket", "tab.txt", gomock.Any()).
			DoAndReturn(func(ctx context.Context, bucket, object string, r io.Reader) error {
				buf := new(bytes.Buffer)
				_, _ = io.Copy(buf, r)
				require.Equal(t, tabBody, buf.String())
				return nil
			})

		tabRepo.EXPECT().
			Create(ctx, gomock.Any()).
			DoAndReturn(func(ctx context.Context, tab *models.Tab) error {
				require.Equal(t, "song.wav", tab.Name)
				require.Equal(t, "tab.txt", tab.Path)
				require.Empty(t, tab.PresignedURL)
				tab.ID = "generated-id"
				return nil
			})

		tabCreate := &models.TabCreate{
			Name: "song.wav",
			Path: "tab.txt",
			Body: tabBody,
		}

		require.NoError(t, svc.CreateTab(ctx, tabCreate))
		require.Equal(t, "generated-id", tabCreate.ID)
	})

	t.Run("db error", func(t *testing.T) {
		svc, tabRepo, _, storage := setupService(t)

		storage.EXPECT().
			UploadFile(ctx, "tabs-bucket", "tab.txt", gomock.Any()).
			Return(nil)

		tabRepo.EXPECT().
			Create(ctx, gomock.Any()).
			Return(errors.New("db fail"))

		tabCreate := &models.TabCreate{
			Name: "song.wav",
			Path: "tab.txt",
			Body: tabBody,
		}

		err := svc.CreateTab(ctx, tabCreate)
		require.ErrorContains(t, err, "db fail")
	})

	t.Run("upload error", func(t *testing.T) {
		svc, _, _, storage := setupService(t)

		storage.EXPECT().
			UploadFile(ctx, "tabs-bucket", "tab.txt", gomock.Any()).
			Return(errors.New("upload fail"))

		tabCreate := &models.TabCreate{
			Name: "song.wav",
			Path: "tab.txt",
			Body: tabBody,
		}

		err := svc.CreateTab(ctx, tabCreate)
		require.ErrorContains(t, err, "upload fail")
	})
}

func TestTabService_DeleteTab(t *testing.T) {
	ctx := context.Background()
	tab := &models.Tab{ID: "tab-id", Path: "tab.txt"}

	t.Run("success", func(t *testing.T) {
		svc, tabRepo, _, storage := setupService(t)

		tabRepo.EXPECT().Get(ctx, "tab-id").Return(tab, nil)
		tabRepo.EXPECT().MarkDeleted(ctx, "tab-id").Return(nil)
		storage.EXPECT().RemoveFile(ctx, "tabs-bucket", "tab.txt").Return(nil)

		require.NoError(t, svc.DeleteTab(ctx, "tab-id"))
	})

	t.Run("storage error is not fatal", func(t *testing.T) {
		svc, tabRepo, _, storage := setupService(t)

		tabRepo.EXPECT().Get(ctx, "tab-id").Return(tab, nil)
		tabRepo.EXPECT().MarkDeleted(ctx, "tab-id").Return(nil)
		storage.EXPECT().RemoveFile(ctx, "tabs-bucket", "tab.txt").
			Return(errors.New("fs fail"))

		require.NoError(t, svc.DeleteTab(ctx, "tab-id"))
	})

	t.Run("get tab error", func(t *testing.T) {
		svc, tabRepo, _, _ := setupService(t)

		tabRepo.EXPECT().Get(ctx, "tab-id").
			Return(nil, errors.New("not found"))

		err := svc.DeleteTab(ctx, "tab-id")
		require.ErrorContains(t, err, "not found")
	})

	t.Run("mark deleted error", func(t *testing.T) {
		svc, tabRepo, _, _ := setupService(t)

		tabRepo.EXPECT().Get(ctx, "tab-id").Return(tab, nil)
		tabRepo.EXPECT().MarkDeleted(ctx, "tab-id").
			Return(errors.New("db fail"))

		err := svc.DeleteTab(ctx, "tab-id")
		require.ErrorContains(t, err, "db fail")
	})
}

func TestTabService_FindTabsByNameLike(t *testing.T) {
	svc, tabRepo, _, _ := setupService(t)
	ctx := context.Background()

	tabs := []*models.Tab{
		{Name: "song1"},
		{Name: "song2"},
	}

	tabRepo.EXPECT().
		FindByNameLike(ctx, "song").
		Return(tabs, nil)

	res, err := svc.FindTabsByNameLike(ctx, "song")

	require.NoError(t, err)
	require.Equal(t, tabs, res)
}
