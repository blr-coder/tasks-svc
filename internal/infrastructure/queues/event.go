package queues

import "context"

type IQueueEventSender interface {
	Send(ctx context.Context, key, value []byte) error
}
