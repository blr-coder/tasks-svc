package events

import "context"

type IEventSender interface {
	Send(ctx context.Context, event *Event) error
}

type Event struct {
	Topic Topic  `json:"topic"`
	Name  Name   `json:"name"`
	Data  []byte `json:"value"`
}

type Name string

const (
	TaskCreated Name = "task_created"
	TaskDeleted Name = "task_deleted"
)

type Topic string

const (
	CUDTopic       Topic = "TASK_CUD_STREAMING_TOPIC"
	LifecycleTopic Topic = "TASK_LIFECYCLE_STREAMING_TOPIC"
)
