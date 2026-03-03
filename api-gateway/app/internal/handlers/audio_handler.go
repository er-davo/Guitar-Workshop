package handlers

import (
	"io"
	"net/http"
	"time"

	"api-gateway/internal/api"
	"api-gateway/internal/models"
	"api-gateway/internal/service"

	"github.com/labstack/echo"
	"go.uber.org/zap"
)

type AudioSepTaskHandler struct {
	service *service.AudioSepTaskService
	log     *zap.Logger
}

func NewAudioSepTaskHandler(service *service.AudioSepTaskService, log *zap.Logger) *AudioSepTaskHandler {
	return &AudioSepTaskHandler{
		service: service,
		log:     log,
	}
}

// POST /audio/separation
func (h *AudioSepTaskHandler) PostAudioSeparation(c echo.Context) error {
	start := time.Now()

	fileHeader, err := c.FormFile("audio_file")
	if err != nil {
		h.log.Warn("audio_file missing")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "audio_file is required"})
	}

	file, err := fileHeader.Open()
	if err != nil {
		h.log.Error("failed to open audio_file", zap.String("filename", fileHeader.Filename), zap.Error(err))
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "failed to open audio_file"})
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		h.log.Error("failed to read audio_file", zap.String("filename", fileHeader.Filename), zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to read audio_file"})
	}

	h.log.Info("creating audio separation task",
		zap.String("filename", fileHeader.Filename),
		zap.Int64("size_bytes", int64(len(data))),
	)

	task, err := h.service.PostAudioSepTask(
		c.Request().Context(),
		models.AudioSepTaskData{
			AudioFileName: fileHeader.Filename,
			Data:          data,
			Size:          int64(len(data)),
		},
	)
	if err != nil {
		h.log.Error("failed to create audio separation task",
			zap.String("filename", fileHeader.Filename),
			zap.Error(err),
		)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to create audio separation task"})
	}

	h.log.Info("audio separation task created",
		zap.String("task_id", task.ID),
		zap.String("filename", fileHeader.Filename),
		zap.Duration("duration", time.Since(start)),
	)

	return c.JSON(http.StatusCreated, api.AudioSepTask{
		ID:                       task.ID,
		Status:                   string(task.Status),
		SeparatedAudioSignedURLs: task.SeparatedAudioSignedURLs,
		Error:                    task.Error,
		CreatedAt:                task.CreatedAt,
	})
}

// GET /audio/separation/:id
func (h *AudioSepTaskHandler) GetAudioSeparation(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		h.log.Warn("empty task id on get")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "id is required"})
	}

	task, err := h.service.Get(c.Request().Context(), id)
	if err != nil {
		h.log.Error("failed to fetch audio separation task", zap.String("task_id", id), zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to get audio separation task"})
	}

	h.log.Info("audio separation task fetched", zap.String("task_id", task.ID))
	return c.JSON(http.StatusOK, api.AudioSepTask{
		ID:                       task.ID,
		Status:                   string(task.Status),
		SeparatedAudioSignedURLs: task.SeparatedAudioSignedURLs,
		Error:                    task.Error,
		CreatedAt:                task.CreatedAt,
	})
}

// Register routes
func (h *AudioSepTaskHandler) RegisterRoutes(e *echo.Echo) {
	group := e.Group("/audio/separation")

	group.POST("", h.PostAudioSeparation)
	group.GET("/:id", h.GetAudioSeparation)
}
