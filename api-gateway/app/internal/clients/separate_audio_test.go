package clients

import (
	"api-gateway/internal/mocks"
	"api-gateway/internal/proto/separator"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestSeparateAudio_Success(t *testing.T) {
	mockClient := new(mocks.MockAudioSeparator)

	audioData := []byte("fake_audio")
	fileName := "test.wav"

	expectedResp := &separator.SeparateResponse{
		Stems: map[string]*separator.AudioFileData{
			"guitar": {
				FileName:   "guitar.wav",
				AudioBytes: []byte("guitar_audio"),
			},
		},
	}

	mockClient.
		On("SeparateAudio", mock.Anything, mock.MatchedBy(func(req *separator.SeparateRequest) bool {
			return req.AudioData != nil && req.AudioData.FileName == fileName
		}), mock.Anything).
		Return(expectedResp, nil)

	audioSeparatorClient = mockClient

	resp, err := SeparateAudio(context.Background(), fileName, audioData)

	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, "guitar.wav", resp["guitar"].FileName)

	mockClient.AssertExpectations(t)
}

func TestSeparateAudio_Error(t *testing.T) {
	mockClient := new(mocks.MockAudioSeparator)

	mockClient.
		On("SeparateAudio", mock.Anything, mock.Anything, mock.Anything).
		Return((*separator.SeparateResponse)(nil), assert.AnError)

	audioSeparatorClient = mockClient

	resp, err := SeparateAudio(context.Background(), "badfile.wav", []byte("bad_audio"))

	require.Error(t, err)
	require.Nil(t, resp)

	mockClient.AssertExpectations(t)
}

func TestSeparateAudio_EmptyResponse(t *testing.T) {
	mockClient := new(mocks.MockAudioSeparator)

	mockClient.
		On("SeparateAudio", mock.Anything, mock.Anything, mock.Anything).
		Return(&separator.SeparateResponse{Stems: nil}, nil)

	audioSeparatorClient = mockClient

	resp, err := SeparateAudio(context.Background(), "file.wav", []byte("data"))

	require.Error(t, err)
	require.Nil(t, resp)
	require.EqualError(t, err, "empty response from audio separator")

	mockClient.AssertExpectations(t)
}
