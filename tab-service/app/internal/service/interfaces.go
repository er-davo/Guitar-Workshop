package service

import (
	"context"
	"io"
	"net/url"
	"tab-service/internal/models"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type TabRepository interface {
	Create(ctx context.Context, tab *models.Tab) error
	Delete(ctx context.Context, id string) error
	MarkDeleted(ctx context.Context, id string) error
	Get(ctx context.Context, id string) (*models.Tab, error)
	Search(ctx context.Context, name string, limit, offset uint64) ([]*models.Tab, error)
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
