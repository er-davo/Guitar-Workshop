package models

import "fmt"

type TabRequest struct {
	AudioURL string `json:"audio_url"`
}

type TabResponse struct {
	Tab    string `json:"tab"`
	Status string `json:"status"`
}

func GenerateTab(notes []string) string {
	return fmt.Sprintf(
		"e|-%s-\nB|-%s-\nG|-%s-\nD|--\nA|--\nE|--",
		notes[0], notes[1], notes[2],
	)
}
