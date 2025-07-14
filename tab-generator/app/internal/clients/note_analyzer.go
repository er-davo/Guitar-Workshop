package clients

import (
	"context"
	"tabgen/internal/config"
	"tabgen/internal/logger"
	note_analyzer "tabgen/internal/proto/note-analyzer"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	NoteAnalyzerClient NoteAnalyzer
	NoteAnalyzerConn   *grpc.ClientConn
)

func init() {
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

type NoteAnalyzer interface {
	Analyze(ctx context.Context, in *note_analyzer.AudioRequest, opts ...grpc.CallOption) (*note_analyzer.NoteResponse, error)
}

func CloseClients() {
	if NoteAnalyzerConn != nil {
		NoteAnalyzerConn.Close()
	}
}
