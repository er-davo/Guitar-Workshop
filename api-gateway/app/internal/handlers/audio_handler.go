package handlers

import (
	"api-gateway/internal/service"
	"encoding/base64"
	"net/http"

	"github.com/labstack/echo"
)

type AudioHandler struct {
	service *service.AudioService
}

func NewAudioHandler(service *service.AudioService) *AudioHandler {
	return &AudioHandler{
		service: service,
	}
}

func (h *AudioHandler) SeparateAudio(c echo.Context) error {
	audioFileName, audioFileData, err := parseAudioInput(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	stems, err := h.service.SeparateAudio(c.Request().Context(), audioFileName, audioFileData)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "разделение не удалось: " + err.Error()})
	}

	resp := make(map[string]string, len(stems))
	for name, data := range stems {
		encoded := base64.StdEncoding.EncodeToString(data.AudioBytes)
		resp[name] = "data:audio/wav;base64," + encoded
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"stems": resp,
	})
}
