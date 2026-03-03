package storage

import (
	"context"
	"io"
	"net/url"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/transfermanager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"go.uber.org/zap"
)

type Storage struct {
	client    *s3.Client
	uploader  *transfermanager.Client
	presigner *s3.PresignClient
	log       *zap.Logger
}

func NewStorage(ctx context.Context, log *zap.Logger) (*Storage, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String("https://tecxeovxpzxajmthpagq.storage.supabase.co/storage/v1/s3")
		o.UsePathStyle = true
	})

	return &Storage{
		client:    client,
		presigner: s3.NewPresignClient(client),
		uploader:  transfermanager.New(client),
		log:       log,
	}, nil
}

func (s *Storage) UploadFile(
	ctx context.Context,
	bucketName string,
	objectName string,
	reader io.Reader,
) error {
	_, err := s.uploader.UploadObject(ctx, &transfermanager.UploadObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectName),
		Body:   reader,
	})
	return err
}

func (s *Storage) GetFile(
	ctx context.Context,
	bucketName string,
	objectName string,
) (io.ReadCloser, error) {
	out, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectName),
	})
	if err != nil {
		return nil, err
	}

	return out.Body, nil
}

func (s *Storage) ListFilesByPrefix(
	ctx context.Context,
	bucketName string,
	prefix string,
) ([]types.Object, error) {

	var result []types.Object

	paginator := s3.NewListObjectsV2Paginator(s.client, &s3.ListObjectsV2Input{
		Bucket: aws.String(bucketName),
		Prefix: aws.String(prefix),
	})

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		result = append(result, page.Contents...)
	}

	return result, nil
}

func (s *Storage) PresignedGet(
	ctx context.Context,
	bucket, object string,
	exp time.Duration,
) (*url.URL, error) {

	req, err := s.presigner.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(object),
	}, s3.WithPresignExpires(exp))

	if err != nil {
		return nil, err
	}

	return url.Parse(req.URL)
}

func (s *Storage) RemoveFile(
	ctx context.Context,
	bucketName string,
	objectName string,
) error {
	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectName),
	})
	return err
}
