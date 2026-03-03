package api

import (
	"time"
)

type Tab struct {
	ID           string    `json:"id,omitempty"`
	Name         string    `json:"name"`
	CreatedAt    time.Time `json:"created_at"`
	PresignedURL string    `json:"presigned_url"`
}

type TabCreate struct {
	Name string `json:"name"`
}

type TabGenTask struct {
	ID             string  `json:"id,omitempty"`
	AudioSepTaskID *string `json:"audio_sep_task_id,omitempty"`

	Status string `json:"status"`

	Error *string `json:"error,omitempty"`

	CreatedAt time.Time `json:"created_at"`
}

type TabGenRequest struct {
	Separation bool `json:"separation"`
}

type TabGenTaskResponse struct {
	Task *TabGenTask `json:"task"`
}

type TabGenResponse struct {
	Task *TabGenTask `json:"task"`
	Tab  *Tab        `json:"tab,omitempty"`
}
