package models

type Tab struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name"`
	Path string `json:"path"`
	Body string `json:"body"`
}
