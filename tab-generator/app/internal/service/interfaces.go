//go:generate mockgen -source=interfaces.go -destination=../mocks/service_mocks.go -package=mocks
package service

import (
	"context"
	"io"
	"net/url"
	"tabgen/internal/models"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type TabGenTaskRepository interface {
	Create(ctx context.Context, task *models.TabGenTask) error
	Get(ctx context.Context, id string) (*models.TabGenTask, error)
	GetByAudioSepTaskID(ctx context.Context, id string) (*models.TabGenTask, error)
	MarkSaved(ctx context.Context, id string) error
	TryUpdateStatus(
		ctx context.Context,
		id string,
		fromStatus []models.Status,
		to models.Status,
		errMsg *string,
		resultTabName *string,
	) (bool, error)
}

type AudioSepTaskRepository interface {
	Create(ctx context.Context, task *models.AudioSepTask) error
	Get(ctx context.Context, id string) (*models.AudioSepTask, error)
}

type Storage interface {
	UploadFile(
		ctx context.Context,
		bucketName string,
		objectName string,
		reader io.Reader,
	) error
	GetFile(
		ctx context.Context,
		bucketName string,
		objectName string,
	) (io.ReadCloser, error)
	RemoveFile(
		ctx context.Context,
		bucketName string,
		objectName string,
	) error
	ListFilesByPrefix(
		ctx context.Context,
		bucketName string,
		prefix string,
	) ([]types.Object, error)
	PresignedGet(
		ctx context.Context,
		bucket, object string,
		exp time.Duration,
	) (*url.URL, error)
}

type TxManager interface {
	Do(ctx context.Context, fn func(ctx context.Context) error) error
}

var _ TxManager = TxManagerStub{}

type TxManagerStub struct{}

func (TxManagerStub) Do(ctx context.Context, fn func(context.Context) error) error {
	return fn(ctx)
}
