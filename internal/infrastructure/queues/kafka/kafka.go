package kafka

import (
	"context"
	"github.com/IBM/sarama"
	"github.com/blr-coder/tasks-svc/internal/infrastructure/queues"
	"log"
	"time"
)

type Sender struct {
	kafkaProducer sarama.SyncProducer
}

func NewKafkaSender(producer sarama.SyncProducer) *Sender {
	return &Sender{
		kafkaProducer: producer,
	}
}

func (s *Sender) Send(ctx context.Context, event *queues.Event) error {
	log.Println("Send")
	_, _, err := s.kafkaProducer.SendMessage(&sarama.ProducerMessage{
		Topic:     event.Topic,
		Key:       sarama.StringEncoder(event.Name),
		Value:     sarama.StringEncoder(event.Data),
		Headers:   nil,
		Metadata:  nil,
		Offset:    -1,
		Partition: 0,
		Timestamp: time.Now(),
	})
	if err != nil {
		log.Println("EEEERRRRROOOORRRR", err)
	}
	return nil
}
