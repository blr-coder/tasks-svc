package events

import (
	"context"
	"github.com/blr-coder/tasks-svc/internal/domain/models"
)

type IEventSender interface {
	Send(ctx context.Context, event *Event) error
	SendTaskCreated(ctx context.Context, task *models.Task) error
	SendTaskUpdated(ctx context.Context, task *models.Task) error
	SendTaskDeleted(ctx context.Context, taskID int64) error
}

type Event struct {
	Topic EventTopic `json:"topic"`
	Name  EventName  `json:"name"`
	Data  []byte     `json:"value"`
}

type EventName string

const (
	TaskCreated EventName = "task_created"
	TaskUpdated EventName = "task_updated"
	TaskDeleted EventName = "task_deleted"
)

type EventTopic string

const (
	CUDTopic       EventTopic = "TASK_CUD_STREAMING_TOPIC"
	LifecycleTopic EventTopic = "TASK_LIFECYCLE_STREAMING_TOPIC"
)
