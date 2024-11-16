package kafka

import (
	"context"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/blr-coder/tasks-svc/internal/config"
	"log"
)

const groupID = "task-consumer-group"

type GroupReceiver struct {
	kafkaGroupConsumer   sarama.ConsumerGroup
	consumerGroupHandler ConsumerGroupHandler
}

type ConsumerGroupHandler struct {
	actionFunc func(ctx context.Context, event any) error
}

func (handler ConsumerGroupHandler) RunAction() error {

	log.Println("RUNNING NEEDED ACTION...!!!")

	return nil
}

func (handler ConsumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (handler ConsumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (handler ConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {

	for msg := range claim.Messages() {

		fmt.Println("[Message Recieved] ", " timeStamp:", msg.Timestamp.Format("2006-01-02 15:04:05"), "consumerId:", " - topic:", msg.Topic, " - key:", string(msg.Key), " - msgValue:", string(msg.Value), " - partition:", msg.Partition, " - offset:", msg.Offset)

		if err := handler.RunAction(); err != nil {
			log.Printf("action function error for message (offset %d): %v", msg.Offset, err)
		}

		session.MarkMessage(msg, "")
	}
	return nil
}

func NewGroupReceiver(config config.KafkaConfig) (*GroupReceiver, error) {
	saramaConfig := sarama.NewConfig()

	saramaConfig.Consumer.Offsets.Initial = sarama.OffsetNewest

	groupConsumer, err := sarama.NewConsumerGroup([]string{config.Address}, groupID, saramaConfig)
	if err != nil {
		_ = groupConsumer.Close()
		return nil, err
	}

	return &GroupReceiver{kafkaGroupConsumer: groupConsumer}, nil
}

func (gr *GroupReceiver) ReceiveWithAction(ctx context.Context, actionFunc func(ctx context.Context, event any) error) error {
	defer func() {
		if err := gr.kafkaGroupConsumer.Close(); err != nil {
			log.Printf("error closing receiver: %v", err)
		}
	}()

	// TODO: I think it's not good decision
	handler := ConsumerGroupHandler{
		actionFunc: actionFunc,
	}

	for {
		err := gr.kafkaGroupConsumer.Consume(ctx, []string{"topic"}, handler)
		if err != nil {
			log.Printf("Error consuming messages: %v", err)
		}
	}
}
