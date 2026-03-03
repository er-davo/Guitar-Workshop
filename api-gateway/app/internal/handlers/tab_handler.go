package handlers

import (
	"errors"
	"net/http"
	"strings"

	"api-gateway/internal/api"
	"api-gateway/internal/repository"
	"api-gateway/internal/service"

	"github.com/labstack/echo"
	"go.uber.org/zap"
)

type TabHandler struct {
	service *service.TabService
	log     *zap.Logger
}

func NewTabHandler(service *service.TabService, log *zap.Logger) *TabHandler {
	return &TabHandler{
		service: service,
		log:     log,
	}
}

// // POST /tab
// func (h *TabHandler) CreateTab(c echo.Context) error {
// 	tab := api.TabCreate{}
// 	if err := c.Bind(&tab); err != nil {
// 		h.log.Warn("invalid tab payload", zap.Error(err))
// 		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
// 	}

// 	if strings.TrimSpace(tab.Name) == "" || strings.TrimSpace(tab.Body) == "" {
// 		h.log.Warn("empty tab name or body", zap.String("tab_name", tab.Name))
// 		return c.JSON(http.StatusBadRequest, map[string]string{"error": "tab name and body cannot be empty"})
// 	}

// 	re := regexp.MustCompile(`[^\w\d_-]+`)
// 	sanitizedName := re.ReplaceAllString(tab.Name, "_")
// 	tab.Name = sanitizedName

// 	if err := h.service.CreateTab(c.Request().Context(), &models.TabCreate{
// 		Name: tab.Name,
// 		Body: tab.Body,
// 	}); err != nil {
// 		h.log.Error("failed to create tab",
// 			zap.String("tab_name", tab.Name),
// 			zap.String("sanitized_name", tab.Name),
// 			zap.Error(err),
// 		)
// 		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to create tab"})
// 	}

// 	h.log.Info("tab created successfully",
// 		zap.String("tab_id", tab.ID),
// 		zap.String("tab_name", tab.Name),
// 		zap.String("sanitized_name", tab.Name),
// 	)

// 	return c.JSON(http.StatusOK, map[string]string{
// 		"id":      tab.ID,
// 		"message": "tab saved",
// 	})
// }

// GET /tab/:id
func (h *TabHandler) GetTab(c echo.Context) error {
	id := c.Param("id")
	tab, err := h.service.GetTabByID(c.Request().Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			h.log.Warn("tab not found", zap.String("tab_id", id))
			return c.JSON(http.StatusNotFound, map[string]string{"error": "tab not found"})
		}
		h.log.Error("failed to fetch tab", zap.String("tab_id", id), zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to fetch tab"})
	}

	h.log.Info("tab fetched", zap.String("tab_id", tab.ID), zap.String("tab_name", tab.Name))
	return c.JSON(http.StatusOK, api.Tab{
		ID:           tab.ID,
		Name:         tab.Name,
		PresignedURL: tab.PresignedURL,
	})
}

// DELETE /tab/:id
func (h *TabHandler) DeleteTab(c echo.Context) error {
	id := c.Param("id")
	err := h.service.DeleteTab(c.Request().Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			h.log.Warn("tab not found on delete", zap.String("tab_id", id))
			return c.JSON(http.StatusNotFound, map[string]string{"error": "tab not found"})
		}
		h.log.Error("failed to delete tab", zap.String("tab_id", id), zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to delete tab"})
	}

	h.log.Info("tab deleted successfully", zap.String("tab_id", id))
	return c.NoContent(http.StatusNoContent)
}

// GET /tab/search?name=...
func (h *TabHandler) SearchTabs(c echo.Context) error {
	name := c.QueryParam("name")
	if strings.TrimSpace(name) == "" {
		h.log.Warn("search query empty")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "name query param is required"})
	}

	tabs, err := h.service.FindTabsByNameLike(c.Request().Context(), name)
	if err != nil {
		h.log.Error("failed to search tabs", zap.String("query", name), zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to search tabs"})
	}

	if len(tabs) == 0 {
		h.log.Warn("no tabs found for search", zap.String("query", name))
		return c.JSON(http.StatusNotFound, map[string]string{"error": "no tabs found"})
	}

	h.log.Info("tabs found for search", zap.String("query", name), zap.Int("count", len(tabs)))
	return c.JSON(http.StatusOK, tabs)
}

// Register routes
func (h *TabHandler) RegisterRoutes(e *echo.Echo) {
	tabGroup := e.Group("/tab")

	// tabGroup.POST("/", h.CreateTab)
	tabGroup.GET("/search", h.SearchTabs)
	tabGroup.GET("/:id", h.GetTab)
	tabGroup.DELETE("/:id", h.DeleteTab)
}
