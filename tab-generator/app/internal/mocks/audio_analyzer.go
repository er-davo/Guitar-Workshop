package mocks

import (
	"context"

	note_analyzer "tabgen/internal/proto/note-analyzer"

	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

type MockNoteAnalyzerClient struct {
	mock.Mock
}

func (m *MockNoteAnalyzerClient) Analyze(ctx context.Context, in *note_analyzer.AudioRequest, opts ...grpc.CallOption) (*note_analyzer.NoteResponse, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*note_analyzer.NoteResponse), args.Error(1)
}
