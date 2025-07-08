package clients

import (
	"api-gateway/internal/logger"
	"api-gateway/internal/proto/separator"
	"context"
	"errors"

	"go.uber.org/zap"
)

func SeparateAudio(ctx context.Context, fileName string, audioData []byte) (map[string]*separator.AudioFileData, error) {
	fileData := &separator.AudioFileData{
		FileName:   fileName,
		AudioBytes: audioData,
	}

	resp, err := audioSeparatorClient.SeparateAudio(ctx, &separator.SeparateRequest{
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
