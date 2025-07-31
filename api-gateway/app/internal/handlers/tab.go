package handlers

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"api-gateway/internal/clients"
	"api-gateway/internal/logger"
	"api-gateway/internal/models"

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
	tab := models.Tab{}
	err := c.Bind(&tab)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	tab.Path = "generated/" + tab.Name + ".txt"

	tabRepo.Create(c.Request().Context(), &tab)

	clients.Supabase.Storage.UploadFile(
		"tabs",
		tab.Path,
		bytes.NewReader([]byte(tab.Body)),
	)

	return c.JSON(http.StatusOK, map[string]string{})
}

func SearchTabs(c echo.Context) error {
	name := c.QueryParam("name")
	if name == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "name query param is required"})
	}

	tabs, err := tabRepo.FindByNameLike(c.Request().Context(), name)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	if len(tabs) == 0 {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "no tabs found"})
	}

	return c.JSON(http.StatusOK, tabs)
}

func GetTab(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}

	tab, err := tabRepo.GetByID(c.Request().Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "tab not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, tab)
}
