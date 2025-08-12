package service

import (
	"api-gateway/internal/clients"
	"api-gateway/internal/proto/separator"
	"context"
)

type AudioService struct {
	audioClient clients.AudioSeparator
}

func NewAudioService(audioClient clients.AudioSeparator) *AudioService {
	return &AudioService{audioClient: audioClient}
}

func (s *AudioService) SeparateAudio(ctx context.Context, audioFileName string, audioFileData []byte) (map[string]*separator.AudioFileData, error) {
	return s.audioClient.SeparateAudio(ctx, audioFileName, audioFileData)
}
