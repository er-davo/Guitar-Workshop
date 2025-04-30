package service

import (
	"context"

	audiopb "tabgen/internal/audioproto"
	"tabgen/internal/logger"
	"tabgen/internal/models"
	tabpb "tabgen/internal/tabproto"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type TabService struct {
	tabpb.UnimplementedTabGenerateServer
}

func (s *TabService) GenerateTab(ctx context.Context, req *tabpb.TabRequest) (*tabpb.TabResponse, error) {
	conn, err := grpc.NewClient(
		"audio-analyzer:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		logger.Log.Fatal("gRPC connection failed", zap.Error(err))
	}
	defer conn.Close()

	client := audiopb.NewAudioAnalyzerClient(conn)

	audioResp, err := client.ProcessAudio(context.Background(), &audiopb.AudioRequest{
		AudioPath: req.AudioUrl,
	})
	if err != nil {
		logger.Log.Error("error on process audio", zap.Error(err))
		return nil, err
	}

	tab, err := models.GenerateTab(audioResp)
	if err != nil {
		return nil, err
	}

	return &tabpb.TabResponse{Tab: tab}, nil
}
