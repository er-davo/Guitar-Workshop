package main

import (
	"context"
	"log"
	"net/http"

	audiopb "tabgen/internal/audioproto"
	"tabgen/internal/logger"
	"tabgen/internal/models"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	defer logger.Log.Sync()

	conn, err := grpc.NewClient(
		"audio-analyzer:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := audiopb.NewAudioAnalyzerClient(conn)

	e := echo.New()

	e.POST("/tab-generate", func(c echo.Context) error {
		var tabReq models.TabRequest
		if err := c.Bind(&tabReq); err != nil {
			return c.JSON(http.StatusBadRequest, "Invalid request")
		}

		audioResp, err := client.ProcessAudio(context.Background(), &audiopb.AudioRequest{
			AudioPath: tabReq.AudioURL,
		})
		if err != nil {
			logger.Log.Error("error on process audio", zap.Error(err))
			return c.JSON(http.StatusInternalServerError, "Audio analysis failed")
		}

		tab, err := models.GenerateTab(audioResp)

		if err != nil {
			return c.JSON(http.StatusInternalServerError, "Tab generarion failed")
		}

		return c.JSON(http.StatusOK, models.TabResponse{
			Tab:    tab,
			Status: "success",
		})

	})

	e.Logger.Fatal(e.Start(":8080"))

}
