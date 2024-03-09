package queues

import "context"

type IQueueEventSender interface {
	Send(ctx context.Context, event *Event) error
}

type Event struct {
	Topic string `json:"topic"`
	Name  string `json:"name"`
	Data  []byte `json:"value"`
}
