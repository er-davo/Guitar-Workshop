package service

import "errors"

var (
	ErrMissingSeparatedStem = errors.New("audio separation result missing stem")
)
