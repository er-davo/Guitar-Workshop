package handlers

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"api-gateway/internal/clients"
	"api-gateway/internal/logger"

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

		separatedFiles, err := clients.SeparateAudio(context.Background(), audioURL, wavData)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}
		otherStem, ok := separatedFiles["other"]
		if !ok {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Missing 'other' stem"})
		}

		audioURL = otherStem.FileName
		wavData = otherStem.AudioBytes
	}

	tabResp, err := clients.GenerateTab(context.Background(), audioURL, wavData)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"tab": tabResp.Tab})
}

func SeparateAudio(c echo.Context) error {
	fileHeader, err := c.FormFile("audio_file")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "не удалось прочитать файл"})
	}

	file, err := fileHeader.Open()
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ошибка открытия файла"})
	}
	defer file.Close()

	audioBytes, err := io.ReadAll(file)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "не удалось прочитать байты файла"})
	}

	stems, err := clients.SeparateAudio(c.Request().Context(), fileHeader.Filename, audioBytes)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "разделение не удалось: " + err.Error()})
	}

	resp := make(map[string]string)
	for name, data := range stems {
		encoded := base64.StdEncoding.EncodeToString(data.AudioBytes)
		resp[name] = "data:audio/wav;base64," + encoded
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"stems": resp,
	})
}
