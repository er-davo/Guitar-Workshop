package handlers

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"api-gateway/internal/logger"
	"api-gateway/internal/proto/separator"
	"api-gateway/internal/proto/tab"

	"github.com/labstack/echo"
)

const (
	FILE = iota
	YOUTUBE
)

func TabGenerate(c echo.Context) error {
	reqType, err := strconv.Atoi(c.FormValue("type"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid type format"})
	}

	var audioURL string
	var wavData []byte

	switch reqType {
	case FILE:
		fileHeader, err := c.FormFile("audio_url")
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "No file uploaded"})
		}

		file, err := fileHeader.Open()
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Could not open file"})
		}
		defer file.Close()

		data, err := io.ReadAll(file)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Failed to read file data"})
		}
		wavData = data

		// TODO: add unique  file name generation
		audioURL = fileHeader.Filename

	case YOUTUBE:
		//TODO
		audioURL = c.FormValue("audio_url")
	default:
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid type"})
	}

	separationEnabled := c.FormValue("separation")

	logger.Debug(fmt.Sprintf("gotseparation value: %s", separationEnabled))

	if separationEnabled == "1" {
		logger.Log.Info("separating audio")
		audioData := separator.AudioFileData{
			FileName:   audioURL,
			AudioBytes: wavData,
		}

		separatedFiles, err := AudioSeparatorClient.SeparateAudio(context.Background(), &separator.SeparateRequest{
			AudioData: &audioData,
		})
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}
		otherStem, ok := separatedFiles.Stems["other"]
		if !ok {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Missing 'other' stem"})
		}

		audioURL = otherStem.FileName
		wavData = otherStem.AudioBytes
	}

	tabResp, err := TabGenClient.GenerateTab(context.Background(), &tab.TabRequest{
		Audio: &tab.AudioFileData{
			FileName: audioURL, AudioBytes: wavData,
		}},
	)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"tab": tabResp.Tab})
}
