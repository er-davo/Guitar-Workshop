package models

import (
	"time"
)

type AudioSepTask struct {
	ID string `json:"id,omitempty"`

	Status Status `json:"status,omitempty"`

	InputAudioName   string
	SeparatedDirName string

	SeparatedAudioSignedURLs map[string]string `json:"separated_audio_signed_urls,omitempty"`

	Error *string `json:"error,omitempty"`

	CreatedAt time.Time `json:"created_at"`
}

func (a *AudioSepTask) AudioFileObjectName() string {
	return a.ID + "/" + a.InputAudioName
}

func (a *AudioSepTask) SeparatedPrefix() string {
	return a.SeparatedDirName
}

type AudioSepTaskCompletedEvent struct {
	ID string `json:"id"`
}
