package main

import (
	"context"
	"github.com/blr-coder/tasks-svc/internal/config"
	"github.com/blr-coder/tasks-svc/internal/infrastructure/queues/kafka"
	"log"
)

func main() {
	err := runTaskActionConsumer()
	if err != nil {
		log.Fatal(err)
	}
}

func runTaskActionConsumer() error {
	log.Println("RUN TASK ACTION CONSUMER")

	ctx := context.Background()

	appConfig, err := config.NewAppConfig()
	if err != nil {
		return err
	}

	r, err := kafka.NewReceiver(appConfig.KafkaConfig)
	if err != nil {
		return err
	}

	err = r.Receive(ctx)
	if err != nil {
		return err
	}

	return nil
}
