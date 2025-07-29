package handlers

import (
	"context"
	"fmt"
	"net/http"

	"api-gateway/internal/clients"
	"api-gateway/internal/logger"

	"github.com/labstack/echo"
)

const (
	FILE = iota
	YOUTUBE
)

func TabGenerate(c echo.Context) error {
	audioFileName, audioFileData, err := parseAudioInput(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	separationEnabled := c.FormValue("separation")

	logger.Debug(fmt.Sprintf("gotseparation value: %s", separationEnabled))

	if separationEnabled == "1" {
		logger.Log.Info("separating audio")

		separatedFiles, err := clients.SeparateAudio(context.Background(), audioFileName, audioFileData)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}
		otherStem, ok := separatedFiles["other"]
		if !ok {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Missing 'other' stem"})
		}

		audioFileName = otherStem.FileName
		audioFileData = otherStem.AudioBytes
	}

	tabResp, err := clients.GenerateTab(context.Background(), audioFileName, audioFileData)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"tab": tabResp.Tab})
}

func SaveTab(c echo.Context) error {

	return c.JSON(http.StatusOK, map[string]string{})
}
