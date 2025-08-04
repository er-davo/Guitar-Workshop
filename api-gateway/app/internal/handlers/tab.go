package handlers

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"api-gateway/internal/clients"
	"api-gateway/internal/logger"
	"api-gateway/internal/models"

	"github.com/labstack/echo"
	"go.uber.org/zap"
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

	if separationEnabled == "1" {
		logger.Log.Info("separating audio")

		separatedFiles, err := clients.SeparateAudio(context.Background(), audioFileName, audioFileData)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}
		otherStem, ok := separatedFiles["other"]
		if !ok {
			logger.Log.Error("missing 'other' stem after audio separation")
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Missing 'other' stem"})
		}

		audioFileName = otherStem.FileName
		audioFileData = otherStem.AudioBytes
	}

	tabResp, err := clients.GenerateTab(context.Background(), audioFileName, audioFileData)
	if err != nil {
		logger.Log.Error("error on generating tab", zap.Error(err))
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	logger.Log.Info("tab successfully generated")
	return c.JSON(http.StatusOK, map[string]string{"tab": tabResp.Tab})
}

func SaveTab(c echo.Context) error {
	tab := models.Tab{}
	err := c.Bind(&tab)
	if err != nil {
		logger.Log.Warn("invalid tab payload", zap.Error(err))
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	if strings.TrimSpace(tab.Name) == "" {
		logger.Log.Warn("tab name is empty")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Tab name cannot be empty"})
	}
	if strings.TrimSpace(tab.Body) == "" {
		logger.Log.Warn("tab body is empty")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Tab body cannot be empty"})
	}

	re := regexp.MustCompile(`[^\w\d_-]+`)
	sanitizedName := re.ReplaceAllString(tab.Name, "_")
	tab.Path = "generated/" + sanitizedName + ".txt"

	err = tabRepo.Create(c.Request().Context(), &tab)
	if err != nil {
		if err.Error() == "The resource already exists" {
			logger.Log.Warn("tab already exists", zap.String("name", tab.Name))
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Tab with that name already exists"})
		}
		logger.Log.Error("failed to create tab in database", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error on creating tab in database: " + err.Error()})
	}

	_, err = clients.Supabase.Storage.UploadFile(
		"tabs",
		tab.Path,
		bytes.NewReader([]byte(tab.Body)),
	)
	if err != nil {
		logger.Log.Error("failed to upload tab to storage", zap.String("path", tab.Path), zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "error on uploading tab to storage"})
	}

	logger.Log.Info("tab saved successfully", zap.String("id", tab.ID), zap.String("name", tab.Name))

	return c.JSON(http.StatusOK, map[string]string{"message": "tab saved", "id": fmt.Sprint(tab.ID)})
}

func SearchTabs(c echo.Context) error {
	name := c.QueryParam("name")
	if name == "" {
		logger.Log.Info("missing name query param")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "name query param is required"})
	}

	tabs, err := tabRepo.FindByNameLike(c.Request().Context(), name)
	if err != nil {
		logger.Log.Error("error on searching tabs", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	if len(tabs) == 0 {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "no tabs found"})
	}

	return c.JSON(http.StatusOK, tabs)
}

func GetTab(c echo.Context) error {
	id := c.Param("id")

	tab, err := tabRepo.GetByID(c.Request().Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "tab not found"})
		}
		logger.Log.Info(err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	buf, err := clients.Supabase.Storage.DownloadFile("tabs", tab.Path)
	if err != nil {
		logger.Log.Error("error on downloading file from storage", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "can not load tab from storage: " + err.Error()})
	}

	tab.Body = string(buf)

	logger.Log.Info("returning tab " + fmt.Sprintf("%+v", tab))

	return c.JSON(http.StatusOK, tab)
}

func ViewTabPage(c echo.Context) error {
	// id := c.Param("id")

	// tab, err := tabRepo.GetByID(c.Request().Context(), id)
	// if err != nil {
	// 	return c.String(http.StatusNotFound, "tab not found")
	// }

	// return c.Render(http.StatusOK, "static/view.html", map[string]any{
	// 	"Tab": tab,
	// })

	return c.File("static/view.html")
}
