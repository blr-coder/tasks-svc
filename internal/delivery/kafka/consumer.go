package kafka

import (
	"encoding/json"
	"github.com/IBM/sarama"
	"github.com/davecgh/go-spew/spew"
	"log"
)

type Consumer struct {
	partitionConsumer sarama.PartitionConsumer
}

func NewConsumer(partitionConsumer sarama.PartitionConsumer) *Consumer {
	return &Consumer{partitionConsumer: partitionConsumer}
}

func (kc *Consumer) Run() error {
	var event Task

	i := 0
	for ; ; i++ {
		msg := <-kc.partitionConsumer.Messages()

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

type Task struct {
	ID          int64
	Title       string
	Description string
}
