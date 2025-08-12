package handlers

import (
	"database/sql"
	"net/http"
	"regexp"
	"strings"

	"api-gateway/internal/logger"
	"api-gateway/internal/models"
	"api-gateway/internal/service"

	"github.com/labstack/echo"
	"go.uber.org/zap"
)

type TabHandler struct {
	service *service.TabService
}

func NewTabHandler(service *service.TabService) *TabHandler {
	return &TabHandler{
		service: service,
	}
}

func (h *TabHandler) TabGenerate(c echo.Context) error {
	audioFileName, audioFileData, err := parseAudioInput(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	separation := c.FormValue("separation") == "1"

	tab, err := h.service.GenerateTab(c.Request().Context(), audioFileName, audioFileData, separation)
	if err != nil {
		logger.Log.Error("error on generating tab", zap.Error(err))
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	logger.Log.Info("tab successfully generated")
	return c.JSON(http.StatusOK, map[string]string{"tab": tab.Body})
}

func (h *TabHandler) CreateTab(c echo.Context) error {
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

	err = h.service.CreateTab(c.Request().Context(), &tab)
	if err != nil {
		logger.Log.Error("failed to upload tab to storage", zap.String("path", tab.Path), zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "error on uploading tab to storage"})
	}

	logger.Log.Info("tab saved successfully", zap.String("id", tab.ID), zap.String("name", tab.Name))

	return c.JSON(http.StatusOK, map[string]string{"message": "tab saved", "id": tab.ID})
}

func (h *TabHandler) GetTab(c echo.Context) error {
	id := c.Param("id")

	tab, err := h.service.GetTabByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, tab)
}

func (h *TabHandler) DeleteTab(c echo.Context) error {
	id := c.Param("id")

	err := h.service.DeleteTab(c.Request().Context(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "tab not found"})
		} else {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *TabHandler) SearchTabs(c echo.Context) error {
	name := c.QueryParam("name")
	if name == "" {
		logger.Log.Info("missing name query param")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "name query param is required"})
	}

	tabs, err := h.service.FindTabsByNameLike(c.Request().Context(), name)
	if err != nil {
		logger.Log.Error("error on searching tabs", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	if len(tabs) == 0 {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "no tabs found"})
	}

	return c.JSON(http.StatusOK, tabs)
}

func (h *TabHandler) ViewTabPage(c echo.Context) error {
	return c.File("static/view.html")
}
