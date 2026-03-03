package models

import (
	"time"
)

type Tab struct {
	ID           string     `json:"id,omitempty"`
	Name         string     `json:"name"`
	Path         string     `json:"path"`
	CreatedAt    time.Time  `json:"created_at"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty"`
	PresignedURL string     `json:"presigned_url"`
}

type TabCreate struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name"`
	Path string `json:"path"`
	Body string `json:"body"`
}

type TabGenTask struct {
	ID             string  `json:"id,omitempty"`
	AudioSepTaskID *string `json:"audio_sep_task_id,omitempty"`

	Status Status `json:"status"`

	InputAudioName string  `json:"input_audio_name"`
	ResultTabName  *string `json:"result_tab_name,omitempty"`

	Saved bool    `json:"saved"`
	Error *string `json:"error,omitempty"`

	CreatedAt time.Time `json:"created_at"`
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

type StartTabGenTask struct {
	ID             string  `json:"id,omitempty"`
	AudioSepTaskID *string `json:"audio_sep_task_id,omitempty"`
}

type PostGenTaskData struct {
	AudioFileName string

	Data []byte
	Size int64

	Separation bool
}

type TabGenResponse struct {
	Task *TabGenTask `json:"task"`
	Tab  *Tab        `json:"tab,omitempty"`
}
