package storage

import (
	"fmt"
	"io"
	"net/http"

	"api-gateway/internal/config"
)

func UploadFileToSupabaseStorage(
	bucketName string, fileName string,
	file io.Reader, contentType string,
) error {
	if bucketName == "" || fileName == "" {
		return fmt.Errorf("bucket and file name must be specified")
	}

	uploadURL := fmt.Sprintf(
		"%s/storage/v1/object/%s/%s",
		config.Load().SupabaseURL,
		bucketName,
		fileName,
	)

	req, err := http.NewRequest("PUT", uploadURL, file)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+config.Load().SupabaseKey)
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Cache-Control", "no-cache")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("upload failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("upload error (%d): %s", resp.StatusCode, string(body))
	}

	return nil
}
