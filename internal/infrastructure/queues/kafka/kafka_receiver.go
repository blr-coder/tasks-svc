package kafka

import (
	"context"
	"github.com/IBM/sarama"
)

type Receiver struct {
	kafkaConsumer sarama.PartitionConsumer
}

func (r *Receiver) Receive(ctx context.Context) error {

	return nil
}
