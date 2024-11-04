package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/blr-coder/tasks-svc/internal/config"
	"github.com/blr-coder/tasks-svc/internal/domain/models"
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

func (r *Receiver) ReceiveWithAction(ctx context.Context, actionFunc func(ctx context.Context, event any) error) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case msg, ok := <-r.kafkaConsumer.Messages():
			if !ok {
				return nil // Канал закрыт, завершаем без ошибки
			}

			var event any
			if err := json.Unmarshal(msg.Value, &event); err != nil {
				log.Println(fmt.Errorf("failed to unmarshal message: %w", err))
				continue
			}

			if err := actionFunc(ctx, event); err != nil {
				log.Println(fmt.Errorf("action function error: %w", err))
			}
		}
	}
}

func (r *Receiver) Receive(ctx context.Context) error {
	var event models.TaskAction

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
