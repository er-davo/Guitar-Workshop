package clients

import (
	"api-gateway/internal/proto/tab"
	"context"

	"google.golang.org/grpc"
)

type TabGenerator interface {
	GenerateTab(ctx context.Context, in *tab.TabRequest, opts ...grpc.CallOption) (*tab.TabResponse, error)
}

func GenerateTab(ctx context.Context, fileName string, audioData []byte) (*tab.TabResponse, error) {
	return TabGenClient.GenerateTab(ctx, &tab.TabRequest{
		Audio: &tab.AudioFileData{
			FileName:   fileName,
			AudioBytes: audioData,
		},
	})
}
