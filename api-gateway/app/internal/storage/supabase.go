package storage

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

var (
	supabaseURL string
	supabaseKey string
)

func init() {
	supabaseURL = os.Getenv("SUPABASE_URL")
	supabaseKey = os.Getenv("ACCESS_KEY")
}

func UploadFileToSupabaseStorage(
	bucketName string, fileName string,
	file io.Reader, contentType string,
) error {
	uploadURL := fmt.Sprintf(
		"%s/storage/v1/object/%s/%s",
		supabaseURL,
		bucketName,
		fileName,
	)

	req, err := http.NewRequest("PUT", uploadURL, file)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+supabaseKey)
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
