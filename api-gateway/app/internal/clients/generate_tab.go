package clients

import (
	"context"

	"api-gateway/internal/proto/tab"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type TabGenerator interface {
	GenerateTab(ctx context.Context, fileName string, audioData []byte) (*tab.TabResponse, error)
	Close() error
}

type tabGenerator struct {
	conn   *grpc.ClientConn
	client tab.TabGenerateClient

	log *zap.Logger
}

func NewTabGenerator(target string, log *zap.Logger) (TabGenerator, error) {
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
	return &tabGenerator{
		conn:   conn,
		client: tab.NewTabGenerateClient(conn),
		log:    log,
	}, nil
}

func (tb *tabGenerator) Close() error {
	return tb.conn.Close()
}

func (tb tabGenerator) GenerateTab(ctx context.Context, fileName string, audioData []byte) (*tab.TabResponse, error) {
	return tb.client.GenerateTab(ctx, &tab.TabRequest{
		Audio: &tab.AudioFileData{
			FileName:   fileName,
			AudioBytes: audioData,
		},
	})
}
