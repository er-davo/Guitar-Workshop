package handlers

import (
	"api-gateway/internal/clients"
	"encoding/base64"
	"net/http"

	"github.com/labstack/echo"
)

func SeparateAudio(c echo.Context) error {
	audioFileName, audioFileData, err := parseAudioInput(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	stems, err := clients.SeparateAudio(c.Request().Context(), audioFileName, audioFileData)
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
