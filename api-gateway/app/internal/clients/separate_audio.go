package clients

import (
	"context"
	"errors"

	"api-gateway/internal/logger"
	"api-gateway/internal/proto/separator"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type AudioSeparator interface {
	SeparateAudio(ctx context.Context, in *separator.SeparateRequest, opts ...grpc.CallOption) (*separator.SeparateResponse, error)
}

func SeparateAudio(ctx context.Context, fileName string, audioData []byte) (map[string]*separator.AudioFileData, error) {
	fileData := &separator.AudioFileData{
		FileName:   fileName,
		AudioBytes: audioData,
	}

	resp, err := AudioSeparatorClient.SeparateAudio(ctx, &separator.SeparateRequest{
		AudioData: fileData,
	})
	if err != nil {
		logger.Log.Error("failed to separate audio", zap.Error(err))
		return nil, err
	}

	if resp == nil || resp.Stems == nil {
		return nil, errors.New("empty response from audio separator")
	}

	return resp.Stems, nil
}
