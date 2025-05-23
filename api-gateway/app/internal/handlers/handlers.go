package handlers

import (
	"context"
	"fmt"
	"net/http"
	"slices"
	"strconv"

	"api-gateway/internal/logger"
	"api-gateway/internal/storage"
	tabpb "api-gateway/internal/tabproto"

	"github.com/labstack/echo"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var testFiles = []string{
	"nothing-else-matters.wav",
	"chords.wav",
}

func TabGenerate(c echo.Context) error {
	conn, err := grpc.NewClient(
		"tab-generator:50052",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		logger.Log.Fatal("gRPC connection failed", zap.Error(err))
	}
	defer conn.Close()

	client := tabpb.NewTabGenerateClient(conn)

	intType, err := strconv.Atoi(c.FormValue("type"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid type format"})
	}
	reqType := tabpb.RequestType(intType)

	var audioURL string

	switch reqType {
	case tabpb.RequestType_FILE:
		file, err := c.FormFile("audio_url")
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "No file uploaded"})
		}

		if slices.Contains(testFiles, file.Filename) {
			audioURL = file.Filename
			break
		}

		// TODO: add unique  file name generation
		audioURL = file.Filename

		src, err := file.Open()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to open file"})
		}
		defer src.Close()

		err = storage.UploadFileToSupabaseStorage(
			"audio-bucket",
			audioURL,
			src,
			file.Header.Get("Content-Type"),
		)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": fmt.Sprintf("Upload failed: %v", err),
			})
		}

	case tabpb.RequestType_YOUTUBE:
		audioURL = c.FormValue("audio_url")
	default:
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid type"})
	}

	tabResp, err := client.GenerateTab(context.Background(), &tabpb.TabRequest{
		AudioUrl: audioURL,
		Type:     reqType,
	})

	if err != nil {
		logger.Log.Error("error on process audio", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, "Tab generation failed")
	}

	return c.JSON(http.StatusOK, tabResp)
}
