package models

import "github.com/google/uuid"

type Task struct {
	ID          int64
	Title       string
	Description string
	CustomerID  uuid.UUID
	ExecutorID  uuid.UUID
	Status      Status
}

type CreateTask struct {
	Title       string
	Description string
	CustomerID  uuid.UUID
	ExecutorID  uuid.UUID
}

type Status string

const (
	PendingStatus    Status = "PENDING"
	ProcessingStatus Status = "PROCESSING"
	DoneStatus       Status = "DONE"
)
