package models

type Status string

const (
	StatusCreated              Status = "CREATED"
	StatusPending              Status = "PENDING"
	StatusProcessing           Status = "PROCESSING"
	StatusWaitingForSeparation Status = "WAITING_FOR_SEPARATION"
	StatusDone                 Status = "DONE"
	StatusFailed               Status = "FAILED"
)
