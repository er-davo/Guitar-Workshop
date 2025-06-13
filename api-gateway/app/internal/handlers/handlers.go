package handlers

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"slices"
	"strconv"

	"api-gateway/internal/convertor"
	"api-gateway/internal/proto/audioproc"
	"api-gateway/internal/proto/separator"
	"api-gateway/internal/proto/tab"
	"api-gateway/internal/storage"

	"github.com/labstack/echo"
)

var testFiles = []string{
	"nothing-else-matters.wav",
	"chords.wav",
}

const (
	FILE = iota
	YOUTUBE
)

func TabGenerate(c echo.Context) error {
	reqType, err := strconv.Atoi(c.FormValue("type"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid type format"})
	}

	var audioURL string
	var wavData []byte

	switch reqType {
	case FILE:
		fileHeader, err := c.FormFile("audio_url")
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "No file uploaded"})
		}
		file, err := fileHeader.Open()
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Could not open file"})
		}
		defer file.Close()

		contentType := fileHeader.Header.Get("Content-Type")

		if contentType == "audio/mpeg" { // .mp3
			wavData, err = convertor.MP3ToWAVInMemory(file)
			if err != nil {
				return c.JSON(http.StatusBadRequest, map[string]string{"error": "Could not convert from .mp3 to .wav"})
			}
		} else if contentType != "audio/wav" && contentType != "audio/x-wav" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "usupported file format"})
		} else {
			wavData, err = io.ReadAll(file)
			if err != nil {
				return c.JSON(http.StatusBadRequest, map[string]string{"error": "Could not read file"})
			}
		}

		if slices.Contains(testFiles, fileHeader.Filename) {
			audioURL = fileHeader.Filename
			break
		}

		// TODO: add unique  file name generation
		audioURL = fileHeader.Filename

		err = storage.UploadFileToSupabaseStorage(
			"audio-bucket",
			audioURL,
			file,
			fileHeader.Header.Get("Content-Type"),
		)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": fmt.Sprintf("Upload failed: %v", err),
			})
		}

	case YOUTUBE:
		//TODO
		audioURL = c.FormValue("audio_url")
	default:
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid type"})
	}

	audioData := separator.AudioFileData{
		FileName:   audioURL,
		AudioBytes: wavData,
	}

	separatedFiles, err := AudioSeparatorClient.SeparateAudio(context.Background(), &separator.SeparateRequest{
		AudioData: &audioData,
	})
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	otherStem, ok := separatedFiles.Stems["other"]
	if !ok {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Missing 'other' stem"})
	}

	procConfig := audioproc.AudioProcessingConfig{
		Threshold:   0.01,
		Margin:      300,
		HighPass:    80.0,
		UseBandpass: true,
		BandLow:     80.0,
		BandHigh:    1300.0,
		FadeSamples: 2048,
		SampleRate:  44100,
	}

	procAudio, err := AudioProcessorClient.ProcessAudio(context.Background(), &audioproc.ProcessAudioRequest{
		WavData:  otherStem.AudioBytes,
		FileName: audioURL,
		Config:   &procConfig,
	})
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	chunkConfig := audioproc.ChunkingConfig{
		SampleRate:          44100,
		Threshold:           0.01,
		ChunkMinDurationSec: 8,
		ChunkMaxDurationSec: 15,
		OverlapDurationSec:  1,
	}

	audioChunks, err := AudioProcessorClient.SplitIntoChunks(context.Background(), &audioproc.SplitAudioRequest{
		WavData:  procAudio.WavData,
		FileName: "tempname",
		Config:   &chunkConfig,
	})
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	chunks := make([]*tab.AudioChunk, len(audioChunks.Chunks))
	for i, chunk := range audioChunks.Chunks {
		chunks[i] = &tab.AudioChunk{
			StartTime: chunk.StartTime,
			AudioData: chunk.AudioData,
		}
	}

	tabResp, err := TabGenClient.GenerateTab(context.Background(), &tab.TabRequest{Chunks: chunks})
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, tabResp.Tab)
}
