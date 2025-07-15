package handlers

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"api-gateway/internal/clients"
	"api-gateway/internal/mocks"
	"api-gateway/internal/proto/separator"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestSeparateAudio_Success(t *testing.T) {
	fileName := "input.wav"
	originalAudio := []byte("original_audio")
	guitarAudio := []byte("guitar_bytes")

	mockAudioSeparatorClient := new(mocks.MockAudioSeparator)
	clients.AudioSeparatorClient = mockAudioSeparatorClient
	t.Cleanup(func() { clients.AudioSeparatorClient = nil })

	mockAudioSeparatorClient.
		On("SeparateAudio", mock.Anything, mock.MatchedBy(func(req *separator.SeparateRequest) bool {
			return req.AudioData != nil &&
				req.AudioData.FileName == fileName &&
				bytes.Equal(req.AudioData.AudioBytes, originalAudio)
		}), mock.Anything).
		Return(&separator.SeparateResponse{
			Stems: map[string]*separator.AudioFileData{
				"guitar": {
					FileName:   "guitar.wav",
					AudioBytes: guitarAudio,
				},
			},
		}, nil)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	fileWriter, err := writer.CreateFormFile("audio_file", fileName)
	require.NoError(t, err)
	fileWriter.Write(originalAudio)

	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/separate", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rec := httptest.NewRecorder()

	e := echo.New()
	c := e.NewContext(req, rec)

	err = SeparateAudio(c)

	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rec.Code)

	var resp map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	require.NoError(t, err)

	stems, ok := resp["stems"].(map[string]interface{})
	require.True(t, ok)
	guitarEncoded, ok := stems["guitar"].(string)
	require.True(t, ok)

	expectedEncoded := "data:audio/wav;base64," + base64.StdEncoding.EncodeToString(guitarAudio)
	require.Equal(t, expectedEncoded, guitarEncoded)

	mockAudioSeparatorClient.AssertExpectations(t)
}
