package models

import (
	"encoding/json"
	"github.com/google/uuid"
	"time"
)

type Task struct {
	ID          int64      `json:"id" db:"id"`
	Title       string     `json:"title" db:"title"`
	Description string     `json:"description" db:"description"`
	CustomerID  uuid.UUID  `json:"customer_id" db:"customer_id"`
	ExecutorID  *uuid.UUID `json:"executor_id" db:"executor_id"`
	Status      TaskStatus `json:"status" db:"status"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
}

func (t *Task) ToJson() ([]byte, error) {
	return json.Marshal(t)
}

type CreateTask struct {
	Title       string
	Description string
	CustomerID  uuid.UUID
	ExecutorID  *uuid.UUID
}

type UpdateTask struct {
	ID          int64
	Title       *string
	Description *string
	CustomerID  *uuid.UUID
	ExecutorID  *uuid.UUID
	Status      *TaskStatus
}

type TasksFilter struct {
	CustomerID uuid.UUID
	ExecutorID uuid.UUID
	Status     TaskStatus
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

type TaskStatus string

const (
	PendingStatus    TaskStatus = "PENDING"
	ProcessingStatus TaskStatus = "PROCESSING"
	DoneStatus       TaskStatus = "DONE"
)
