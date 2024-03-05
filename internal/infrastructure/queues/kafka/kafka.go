package kafka

import (
	"context"
	"github.com/IBM/sarama"
	"log"
)

type Sender struct {
	kafkaProducer sarama.SyncProducer
}

func NewKafkaSender(producer sarama.SyncProducer) *Sender {
	return &Sender{
		kafkaProducer: producer,
	}
}

func (s Sender) Send(ctx context.Context, key, value []byte) error {
	log.Println("Send")
	return nil
}
