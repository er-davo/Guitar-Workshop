package models

import (
	"time"
)

// type Tab struct {
// 	ID           string     `json:"id,omitempty"`
// 	Name         string     `json:"name"`
// 	Path         string     `json:"path"`
// 	CreatedAt    time.Time  `json:"created_at"`
// 	DeletedAt    *time.Time `json:"deleted_at,omitempty"`
// 	PresignedURL string     `json:"presigned_url"`
// }

// type TabCreate struct {
// 	ID   string `json:"id,omitempty"`
// 	Name string `json:"name"`
// 	Path string `json:"path"`
// 	Body string `json:"body"`
// }

type TabGenTask struct {
	ID             string
	AudioSepTaskID *string

	Status Status

	InputAudioName string
	ResultTabName  *string

	Saved bool
	Error *string

	CreatedAt time.Time
}

func (t *TabGenTask) AudioFileObjectName() string {
	return t.ID + "/" + t.InputAudioName
}

func (t *TabGenTask) ResultTabObjectName() string {
	if t.ResultTabName == nil {
		return t.ID
	}
	return t.ID + "/" + *t.ResultTabName
}

type TabGenTaskStartEvent struct {
	ID string `json:"id,omitempty"`
}
