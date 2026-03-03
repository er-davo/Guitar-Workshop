package handlers

import (
	"fmt"
	"io"
	"mime/multipart"

	"github.com/labstack/echo"
)

func parseAudioInput(c echo.Context) (*multipart.FileHeader, []byte, error) {
	fileHeader, err := c.FormFile("audio_file")
	if err != nil {
		return nil, nil, fmt.Errorf("no file uploaded")
	}
	file, err := fileHeader.Open()
	if err != nil {
		return nil, nil, fmt.Errorf("could not open file")
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read file data")
	}

	return fileHeader, data, nil
}
