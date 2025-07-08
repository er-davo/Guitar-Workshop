package clients

import (
	"api-gateway/internal/proto/tab"
	"context"
)

func GenerateTab(ctx context.Context, fileName string, audioData []byte) (*tab.TabResponse, error) {
	return tabGenClient.GenerateTab(ctx, &tab.TabRequest{
		Audio: &tab.AudioFileData{
			FileName:   fileName,
			AudioBytes: audioData,
		},
	})
}
