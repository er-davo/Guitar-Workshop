package handlers

import (
	"fmt"
	"io"

	"github.com/labstack/echo"
)

func parseAudioInput(c echo.Context) (string, []byte, error) {
	fileHeader, err := c.FormFile("audio_file")
	if err != nil {
		return "", nil, fmt.Errorf("no file uploaded")
	}
	file, err := fileHeader.Open()
	if err != nil {
		return "", nil, fmt.Errorf("could not open file")
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return "", nil, fmt.Errorf("failed to read file data")
	}

	return fileHeader.Filename, data, nil
}
