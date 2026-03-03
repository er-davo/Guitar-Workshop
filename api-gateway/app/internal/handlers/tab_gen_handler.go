package handlers

import (
	"errors"
	"io"
	"net/http"
	"time"

	"api-gateway/internal/api"
	"api-gateway/internal/models"
	"api-gateway/internal/service"

	"github.com/labstack/echo"
	"go.uber.org/zap"
)

type GenTaskHandler struct {
	service *service.GenTaskService
	log     *zap.Logger
}

func NewGenTaskHandler(service *service.GenTaskService, log *zap.Logger) *GenTaskHandler {
	return &GenTaskHandler{
		service: service,
		log:     log,
	}
}

// POST /generation
func (h *GenTaskHandler) PostGeneration(c echo.Context) error {
	start := time.Now()

	fileHeader, err := c.FormFile("audio_file")
	if err != nil {
		h.log.Warn("audio_file missing")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "audio_file is required"})
	}

	file, err := fileHeader.Open()
	if err != nil {
		h.log.Error("failed to open uploaded file", zap.String("filename", fileHeader.Filename), zap.Error(err))
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "failed to open file"})
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		h.log.Error("failed to read uploaded file", zap.String("filename", fileHeader.Filename), zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to read file"})
	}

	tgReq := &api.TabGenRequest{}
	if err := c.Bind(tgReq); err != nil {
		h.log.Error("failed to bind tab gen request", zap.Error(err))
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "failed to bind tab gen request"})
	}

	h.log.Info("creating generation task",
		zap.String("filename", fileHeader.Filename),
		zap.Bool("separation", tgReq.Separation),
	)

	task, err := h.service.PostGenTask(
		c.Request().Context(),
		models.PostGenTaskData{
			AudioFileName: fileHeader.Filename,
			Data:          data,
			Size:          int64(len(data)),
			Separation:    tgReq.Separation,
		},
	)
	if err != nil {
		h.log.Error("failed to create generation task", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to create generation task"})
	}

	h.log.Info("generation task created",
		zap.String("task_id", task.ID),
		zap.Duration("duration", time.Since(start)),
	)

	return c.JSON(http.StatusCreated, api.TabGenResponse{
		Task: &api.TabGenTask{
			ID:        task.ID,
			Status:    string(task.Status),
			Error:     task.Error,
			CreatedAt: task.CreatedAt,
		},
	})
}

// GET /generation/:id
func (h *GenTaskHandler) GetGenTask(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		h.log.Warn("empty task id")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "id is required"})
	}

	taskTab, err := h.service.Get(c.Request().Context(), id)
	if err != nil {
		h.log.Error("failed to get generation task", zap.String("task_id", id), zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to get generation task"})
	}

	h.log.Info("generation task fetched",
		zap.String("task_id", id),
	)

	resp := api.TabGenResponse{
		Task: &api.TabGenTask{
			ID:        taskTab.Task.ID,
			Status:    string(taskTab.Task.Status),
			Error:     taskTab.Task.Error,
			CreatedAt: taskTab.Task.CreatedAt,
		},
	}

	if taskTab.Tab != nil {
		resp.Tab = &api.Tab{
			ID:           taskTab.Tab.ID,
			Name:         taskTab.Tab.Name,
			CreatedAt:    taskTab.Tab.CreatedAt,
			PresignedURL: taskTab.Tab.PresignedURL,
		}
	}

	return c.JSON(http.StatusOK, resp)
}

// POST /generation/save/:id
func (h *GenTaskHandler) SaveTab(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		h.log.Warn("empty task id")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "id is required"})
	}

	tc := &api.TabCreate{}
	if err := c.Bind(tc); err != nil {
		h.log.Error("failed to bind tab create", zap.Error(err))
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "failed to bind tab create"})
	}

	err := h.service.SaveTab(c.Request().Context(), id, tc.Name)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrTaskNotCompleted):
			h.log.Warn("attempt to save unfinished task", zap.String("task_id", id))
			return c.JSON(http.StatusConflict, map[string]string{"error": "task is not completed"})
		case errors.Is(err, service.ErrNoResultTabName):
			h.log.Warn("attempt to save task without result", zap.String("task_id", id))
			return c.JSON(http.StatusConflict, map[string]string{"error": "task has no result tab"})
		default:
			h.log.Error("failed to save tab from task", zap.String("task_id", id), zap.Error(err))
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to save tab"})
		}
	}

	h.log.Info("tab saved from task", zap.String("task_id", id))

	return c.NoContent(http.StatusNoContent)
}

// Register routes
func (h *GenTaskHandler) RegisterRoutes(e *echo.Echo) {
	group := e.Group("/generation")

	group.POST("", h.PostGeneration)
	group.GET("/:id", h.GetGenTask)
	group.POST("/save/:id", h.SaveTab)
}
