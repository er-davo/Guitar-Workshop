package service

import (
	"context"
	"testing"

	"tabgen/internal/mocks"
	note_analyzer "tabgen/internal/proto/note-analyzer"
	"tabgen/internal/proto/tab"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

func TestGenerateTab(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAnalyzer := mocks.NewMockNoteAnalyzer(ctrl)

	mockAnalyzer.EXPECT().
		Analyze(gomock.Any(), gomock.Any()).
		Return(&note_analyzer.NoteResponse{
			Notes: []*note_analyzer.NoteEvent{
				{MidiPitch: 60, StartSeconds: 0, DurationSeconds: 1, Velocity: 100},
			},
		}, nil)

	log := zap.NewNop()
	service := NewTabService(mockAnalyzer, log)

	resp, err := service.GenerateTab(context.Background(), &tab.TabRequest{
		Audio: &tab.AudioFileData{FileName: "test.wav"},
	})
	require.NoError(t, err)
	require.NotEmpty(t, resp.Tab)
}
