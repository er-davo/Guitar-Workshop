package service

import (
	"context"
	"testing"

	"tabgen/internal/clients"
	"tabgen/internal/mocks"
	note_analyzer "tabgen/internal/proto/note-analyzer"
	"tabgen/internal/proto/tab"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestTabService_GenerateTab_Success(t *testing.T) {
	inputAudio := []byte("fake_wav_bytes")
	fileName := "test.wav"
	req := &tab.TabRequest{
		Audio: &tab.AudioFileData{
			FileName:   fileName,
			AudioBytes: inputAudio,
		},
	}

	mockClient := new(mocks.MockNoteAnalyzerClient)
	clients.NoteAnalyzerClient = mockClient

	mockClient.
		On("Analyze", mock.Anything, mock.MatchedBy(func(a *note_analyzer.AudioRequest) bool {
			return a.AudioData != nil && a.AudioData.FileName == fileName
		})).
		Return(&note_analyzer.NoteResponse{
			Notes: []*note_analyzer.NoteEvent{
				{MidiPitch: 60, StartSeconds: 0.0, DurationSeconds: 1.0, Velocity: 0.9},
				{MidiPitch: 62, StartSeconds: 1.1, DurationSeconds: 1.0, Velocity: 0.85},
			},
		}, nil)

	service := &TabService{}
	resp, err := service.GenerateTab(context.Background(), req)

	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Contains(t, resp.Tab, "e|")
}
