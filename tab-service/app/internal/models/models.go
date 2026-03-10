package models

import "time"

type Tab struct {
	ID           string     `json:"id,omitempty"`
	Name         string     `json:"name"`
	Path         string     `json:"path"`
	CreatedAt    time.Time  `json:"created_at"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty"`
	PresignedURL string     `json:"presigned_url"`
}
