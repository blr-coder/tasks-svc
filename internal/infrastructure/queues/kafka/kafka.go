package kafka

import (
	"context"
	"github.com/IBM/sarama"
	"github.com/blr-coder/tasks-svc/internal/config"
	"github.com/blr-coder/tasks-svc/internal/domain/models"
	"github.com/blr-coder/tasks-svc/internal/events"
	"log"
	"time"
)

type Sender struct {
	kafkaProducer sarama.SyncProducer
}

func NewKafkaSender(config config.KafkaConfig) (*Sender, error) {
	newConfig := sarama.NewConfig()
	newConfig.Producer.RequiredAcks = sarama.WaitForAll
	newConfig.Producer.Partitioner = sarama.NewHashPartitioner
	newConfig.Producer.Return.Successes = true

	syncProducer, err := sarama.NewSyncProducer([]string{config.Address}, newConfig)
	if err != nil {
		return nil, err
	}
	return &Sender{
		kafkaProducer: syncProducer,
	}, nil
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

	log.Println("SENT")
	return nil
}

func (s *Sender) SendTaskCreated(ctx context.Context, task *models.Task) {
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

	log.Println("TaskCreated SENT")
	return
}

func (s *Sender) SendTaskUpdated(ctx context.Context, task *models.Task) {
	log.Println("TaskUpdated")

	jTask, err := task.ToJson()
	if err != nil {
		log.Println("EEEERRRRROOOORRRR", err)
		return
	}

	_, _, err = s.kafkaProducer.SendMessage(&sarama.ProducerMessage{
		Topic:     string(events.CUDTopic),
		Key:       sarama.StringEncoder(events.TaskUpdated),
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

	log.Println("TaskUpdated SENT")
	return
}

func (s *Sender) SendTaskDeleted(ctx context.Context, task *models.Task) {
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

	log.Println("TaskDeleted SENT")
	return
}
