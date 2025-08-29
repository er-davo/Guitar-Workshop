package clients

import (
	"context"
	"errors"

	"api-gateway/internal/proto/separator"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AudioSeparator interface {
	SeparateAudio(ctx context.Context, fileName string, audioData []byte) (map[string]*separator.AudioFileData, error)
	Close() error
}

type audioSeaparator struct {
	conn   *grpc.ClientConn
	client separator.AudioSeparatorClient

	log *zap.Logger
}

func NewAudioSeparator(target string, log *zap.Logger) (AudioSeparator, error) {
	conn, err := grpc.NewClient(
		target,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(250<<20), // 250 MB
			grpc.MaxCallSendMsgSize(250<<20),
		),
	)
	if err != nil {
		return nil, err
	}
	return &audioSeaparator{
		conn:   conn,
		client: separator.NewAudioSeparatorClient(conn),
		log:    log,
	}, nil
}

func (as audioSeaparator) SeparateAudio(ctx context.Context, fileName string, audioData []byte) (map[string]*separator.AudioFileData, error) {
	fileData := &separator.AudioFileData{
		FileName:   fileName,
		AudioBytes: audioData,
	}

	resp, err := as.client.SeparateAudio(ctx, &separator.SeparateRequest{
		AudioData: fileData,
	})
	if err != nil {
		as.log.Error("failed to separate audio", zap.Error(err))
		return nil, err
	}

	if resp == nil || resp.Stems == nil {
		return nil, errors.New("empty response from audio separator")
	}

	return resp.Stems, nil
}

func (as *audioSeaparator) Close() error {
	return as.conn.Close()
}
