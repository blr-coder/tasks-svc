package models

import (
	"encoding/json"
	"github.com/google/uuid"
	"time"
)

type Task struct {
	ID          int64
	Title       string
	Description string
	CustomerID  uuid.UUID
	ExecutorID  uuid.UUID
	Status      Status
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (t *Task) ToJson() ([]byte, error) {
	return json.Marshal(t)
}

type CreateTask struct {
	Title       string
	Description string
	CustomerID  uuid.UUID
	ExecutorID  uuid.UUID
}

type TasksFilter struct {
	CustomerID uuid.UUID
	ExecutorID uuid.UUID
	Status     Status
	Search     string
	Sorting    Sorting
	Limiting   Limiting
}

type Sorting struct {
	SortBy, SortOrder string
}

type Limiting struct {
	Limit, Offset int32
}

type Status string

const (
	PendingStatus    Status = "PENDING"
	ProcessingStatus Status = "PROCESSING"
	DoneStatus       Status = "DONE"
)
