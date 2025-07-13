package mocks

import (
	"api-gateway/internal/proto/separator"
	"api-gateway/internal/proto/tab"
	"context"

	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

type MockTabGenClient struct {
	mock.Mock
}

func (m *MockTabGenClient) GenerateTab(ctx context.Context, in *tab.TabRequest, opts ...grpc.CallOption) (*tab.TabResponse, error) {
	args := m.Called(ctx, in, opts)

	return args.Get(0).(*tab.TabResponse), args.Error(1)
}

type MockAudioSeparator struct {
	mock.Mock
}

func (m *MockAudioSeparator) SeparateAudio(ctx context.Context, in *separator.SeparateRequest, opts ...grpc.CallOption) (*separator.SeparateResponse, error) {
	args := m.Called(ctx, in, opts)

	return args.Get(0).(*separator.SeparateResponse), args.Error(1)
}
