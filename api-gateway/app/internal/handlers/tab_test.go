package handlers

import (
	"api-gateway/internal/clients"
	"api-gateway/internal/mocks"
	"api-gateway/internal/proto/separator"
	"api-gateway/internal/proto/tab"

	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestTabGenerate(t *testing.T) {
	fileName := "input.wav"
	originalAudio := []byte("original_audio")
	separatedFileName := "guitar.wav"
	separatedAudio := []byte("guitar_audio")

	mockAudioSeparatorClient := new(mocks.MockAudioSeparator)
	clients.AudioSeparatorClient = mockAudioSeparatorClient
	t.Cleanup(func() { clients.AudioSeparatorClient = nil })

	mockAudioSeparatorClient.
		On("SeparateAudio", mock.Anything, mock.MatchedBy(func(req *separator.SeparateRequest) bool {
			return req.AudioData != nil && req.AudioData.FileName == fileName
		}), mock.Anything).
		Return(&separator.SeparateResponse{
			Stems: map[string]*separator.AudioFileData{
				"other": {
					FileName:   separatedFileName,
					AudioBytes: separatedAudio,
				},
			},
		}, nil)

	mockTabGenClient := new(mocks.MockTabGenClient)
	clients.TabGenClient = mockTabGenClient
	t.Cleanup(func() { clients.TabGenClient = nil })

	mockTabGenClient.
		On("GenerateTab", mock.Anything, mock.MatchedBy(func(req *tab.TabRequest) bool {
			return req.Audio != nil &&
				req.Audio.FileName == separatedFileName &&
				bytes.Equal(req.Audio.AudioBytes, separatedAudio)
		}), mock.Anything).
		Return(&tab.TabResponse{Tab: "TAB123"}, nil)

	// multipart/form-data
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	fileWriter, err := writer.CreateFormFile("audio_file", fileName)
	require.NoError(t, err)
	fileWriter.Write(originalAudio)

	writer.WriteField("separation", "1")
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/generate", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rec := httptest.NewRecorder()

	e := echo.New()
	c := e.NewContext(req, rec)

	err = TabGenerate(c)

	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rec.Code)

	var result map[string]string
	err = json.Unmarshal(rec.Body.Bytes(), &result)
	require.NoError(t, err)
	require.Equal(t, "TAB123", result["tab"])

	mockAudioSeparatorClient.AssertExpectations(t)
	mockTabGenClient.AssertExpectations(t)
}
