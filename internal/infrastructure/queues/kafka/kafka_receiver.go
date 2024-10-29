package kafka

import (
	"context"
	"encoding/json"
	"github.com/IBM/sarama"
	"github.com/blr-coder/tasks-svc/internal/config"
	"github.com/davecgh/go-spew/spew"
	"log"
)

type Receiver struct {
	kafkaConsumer sarama.PartitionConsumer
}

func NewReceiver(config config.KafkaConfig) (*Receiver, error) {
	consumer, err := sarama.NewConsumer([]string{config.Address}, sarama.NewConfig())
	if err != nil {
		return nil, err
	}

	partitionConsumer, err := consumer.ConsumePartition(config.ConsumerTopic, int32(config.Partition), sarama.OffsetOldest)
	if err != nil {
		_ = consumer.Close() // Закрываем consumer в случае ошибки
		return nil, err
	}

	return &Receiver{kafkaConsumer: partitionConsumer}, nil
}

func (r *Receiver) Receive(ctx context.Context) error {
	var event TaskAction

	i := 0
	for ; ; i++ {
		msg := <-r.kafkaConsumer.Messages()

		err := json.Unmarshal(msg.Value, &event)
		if err != nil {
			return err
		}
		spew.Dump(event)

		if string(msg.Key) == "THE END" {
			i++
			break
		}
	}
	log.Printf("Finished. Received %d messages.\n", i)

	return nil
}

// TODO: Remove it after testing

type TaskAction struct {
	ID     int64
	TaskID int64
	Type   ActionType
	URL    string
}

type ActionType string

const (
	commit       ActionType = "COMMIT"
	mergeRequest ActionType = "MERGE_REQUEST"
	merge        ActionType = "MERGE"
)
