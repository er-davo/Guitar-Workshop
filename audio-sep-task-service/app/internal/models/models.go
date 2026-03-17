package models

import (
	"time"
)

type Status string

func (s Status) String() string {
	return string(s)
}

var (
	StatusCreated    Status = "CREATED"
	StatusPending    Status = "PENDING"
	StatusProcessing Status = "PROCESSING"
	StatusDone       Status = "DONE"
	StatusFailed     Status = "FAILED"
)

type AudioSepTask struct {
	ID string

	Status Status

	InputAudioName   string
	SeparatedDirName string

	SeparatedAudioSignedURLs map[string]string

	Error *string

	CreatedAt time.Time
}

func (a *AudioSepTask) AudioFileObjectName() string {
	return a.ID + "/" + a.InputAudioName
}

func (a *AudioSepTask) SeparatedPrefix() string {
	return a.SeparatedDirName
}

type StartAudioSepTask struct {
	ID string `json:"id,omitempty"`
}

type AudioSepTaskData struct {
	AudioFileName string

	Data []byte
	Size int64
}
