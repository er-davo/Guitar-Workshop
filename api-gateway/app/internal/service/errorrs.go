package service

import "errors"

var (
	ErrTaskNotCompleted    = errors.New("task not completed")
	ErrNoResultTabName     = errors.New("no result tab name")
	ErrNoSeparedAudioFiles = errors.New("no separated audio files")
)
