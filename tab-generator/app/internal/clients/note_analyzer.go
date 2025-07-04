package clients

import (
	"tabgen/internal/config"
	"tabgen/internal/logger"
	"tabgen/internal/proto/note-analyzer"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	NoteAnalyzerClient note_analyzer.NoteAnalyzerClient
	NoteAnalyzerConn   *grpc.ClientConn
)

func InitClients() {
	var err error

	NoteAnalyzerConn, err = grpc.NewClient(
		config.Load().AnalyzerHost+":"+config.Load().AnalyzerPort,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(100*1024*1024), // 100 MB
			grpc.MaxCallSendMsgSize(100*1024*1024),
		),
	)
	if err != nil {
		logger.Log.Fatal("onsets-frames gRPC connection failed", zap.Error(err))
	}

	NoteAnalyzerClient = note_analyzer.NewNoteAnalyzerClient(NoteAnalyzerConn)
}

func CloseClients() {
	if NoteAnalyzerConn != nil {
		NoteAnalyzerConn.Close()
	}
}
