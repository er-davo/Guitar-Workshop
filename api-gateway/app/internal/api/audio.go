package api

import "time"

type AudioSepTask struct {
	ID string `json:"id,omitempty"`

	Status string `json:"status,omitempty"`

	SeparatedAudioSignedURLs map[string]string `json:"separated_audio_signed_urls,omitempty"`

	Error *string `json:"error,omitempty"`

	CreatedAt time.Time `json:"created_at"`
}

type StartAudioSepTask struct {
	ID string `json:"id,omitempty"`
}
