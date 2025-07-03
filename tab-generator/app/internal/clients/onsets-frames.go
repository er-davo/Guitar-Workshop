package clients

import (
	"tabgen/internal/config"
	"tabgen/internal/logger"
	onsets_frames "tabgen/internal/proto/onsets-frames"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	OnsetsAndFramesClient onsets_frames.OnsetsAndFramesClient
	OnsetsAndFramesConn   *grpc.ClientConn
)

func InitClients() {
	var err error

	OnsetsAndFramesConn, err = grpc.NewClient(
		config.Load().AnalyzerHost+":"+config.Load().AnalyzerPort,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		logger.Log.Fatal("onsets-frames gRPC connection failed", zap.Error(err))
	}

	OnsetsAndFramesClient = onsets_frames.NewOnsetsAndFramesClient(OnsetsAndFramesConn)
}

func CloseClients() {
	if OnsetsAndFramesConn != nil {
		OnsetsAndFramesConn.Close()
	}
}
