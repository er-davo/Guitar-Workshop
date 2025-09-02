//go:generate mockgen -source=note_analyzer.go -destination=../mocks/mock_note_analyzer.go -package=mocks
package clients

import (
	"context"

	note_analyzer "tabgen/internal/proto/note-analyzer"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type NoteAnalyzer interface {
	Analyze(ctx context.Context, in *note_analyzer.AudioRequest, opts ...grpc.CallOption) (*note_analyzer.NoteResponse, error)
	Close() error
}

type noteAnalyzer struct {
	conn   *grpc.ClientConn
	client note_analyzer.NoteAnalyzerClient

	log *zap.Logger
}

func NewNoteAnalyzerClient(target string, log *zap.Logger) (NoteAnalyzer, error) {
	conn, err := grpc.NewClient(
		target,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(100*1024*1024), // 100 MB
			grpc.MaxCallSendMsgSize(100*1024*1024),
		),
	)
	if err != nil {
		return nil, err
	}

	return &noteAnalyzer{
		conn:   conn,
		client: note_analyzer.NewNoteAnalyzerClient(conn),
		log:    log,
	}, nil
}

func (n *noteAnalyzer) Analyze(ctx context.Context, in *note_analyzer.AudioRequest, opts ...grpc.CallOption) (*note_analyzer.NoteResponse, error) {
	if in != nil {
		return n.client.Analyze(ctx, in, opts...)
	} else {
		return nil, nil
	}
}

func (n *noteAnalyzer) Close() error {
	if n.conn != nil {
		return n.conn.Close()
	} else {
		return nil
	}
}
