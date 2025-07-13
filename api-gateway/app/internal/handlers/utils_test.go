package handlers

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/require"
)

func TestParseAudioInput_Success(t *testing.T) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("audio_file", "test.wav")
	require.NoError(t, err)
	_, err = io.Copy(part, bytes.NewReader([]byte("FAKEAUDIODATA")))
	require.NoError(t, err)

	writer.Close()

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/upload", body)
	req.Header.Set(echo.HeaderContentType, writer.FormDataContentType())
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	filename, data, err := parseAudioInput(c)

	require.NoError(t, err)
	require.Equal(t, "test.wav", filename)
	require.Equal(t, []byte("FAKEAUDIODATA"), data)
}

func TestParseAudioInput_NoFile(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/upload", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	filename, data, err := parseAudioInput(c)

	require.Error(t, err)
	require.Empty(t, filename)
	require.Nil(t, data)
}
