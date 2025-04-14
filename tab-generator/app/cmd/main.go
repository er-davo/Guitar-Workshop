package main

import (
	"context"
	"log"
	"net/http"

	"tabgen/internal/audioproto"
	"tabgen/internal/models"

	"github.com/labstack/echo/v4"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient(
		"audio-analyzer:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := audioproto.NewAudioAnalyzerClient(conn)

	e := echo.New()

	e.POST("/tab-generate", func(c echo.Context) error {
		var tabReq models.TabRequest
		if err := c.Bind(&tabReq); err != nil {
			return c.JSON(http.StatusBadRequest, "Invalid request")
		}

		audioReq, err := client.ProcessAudio(context.Background(), &audioproto.AudioRequest{
			AudioPath: tabReq.AudioURL,
		})
		if err != nil {
			return c.JSON(http.StatusInternalServerError, "Audio analysis failed")
		}

		tab := models.GenerateTab(audioReq.Notes)

		return c.JSON(http.StatusOK, models.TabResponse{
			Tab:    tab,
			Status: "success",
		})

	})

	e.Logger.Fatal(e.Start(":8080"))

}
