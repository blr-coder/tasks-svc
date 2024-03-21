package kafka

import (
	"context"
	"github.com/IBM/sarama"
	"github.com/blr-coder/tasks-svc/internal/domain/models"
	"github.com/blr-coder/tasks-svc/internal/events"
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

func (s *Sender) Send(ctx context.Context, event *events.Event) error {
	log.Println("Send")
	_, _, err := s.kafkaProducer.SendMessage(&sarama.ProducerMessage{
		Topic:     string(event.Topic),
		Key:       sarama.StringEncoder(event.Name),
		Value:     sarama.ByteEncoder(event.Data),
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

func (s *Sender) TaskCreated(ctx context.Context, task *models.Task) {
	log.Println("TaskCreatedSend")

	jTask, err := task.ToJson()
	if err != nil {
		log.Println("EEEERRRRROOOORRRR", err)
		return
	}

	_, _, err = s.kafkaProducer.SendMessage(&sarama.ProducerMessage{
		Topic:     string(events.CUDTopic),
		Key:       sarama.StringEncoder(events.TaskCreated),
		Value:     sarama.ByteEncoder(jTask),
		Headers:   nil,
		Metadata:  nil,
		Offset:    -1,
		Partition: 0,
		Timestamp: time.Now(),
	})
	if err != nil {
		log.Println("EEEERRRRROOOORRRR", err)
	}

	return
}

func (s *Sender) TaskDeleted(ctx context.Context, task *models.Task) {
	log.Println("TaskDeleted")

	jTask, err := task.ToJson()
	if err != nil {
		log.Println("EEEERRRRROOOORRRR", err)
		return
	}

	_, _, err = s.kafkaProducer.SendMessage(&sarama.ProducerMessage{
		Topic:     string(events.CUDTopic),
		Key:       sarama.StringEncoder(events.TaskDeleted),
		Value:     sarama.ByteEncoder(jTask),
		Headers:   nil,
		Metadata:  nil,
		Offset:    -1,
		Partition: 0,
		Timestamp: time.Now(),
	})
	if err != nil {
		log.Println("EEEERRRRROOOORRRR", err)
	}

	return
}
