package main

import (
	"github.com/IBM/sarama"
	"log"
	"tasks-svc/internal/config"
	"tasks-svc/internal/delivery/kafka"
)

func main() {
	err := runApp()
	if err != nil {
		log.Fatal(err)
	}
}

func runApp() error {
	log.Println("ASYNC")

	appConfig, err := config.NewAppConfig()
	if err != nil {
		return err
	}

	consumer, err := sarama.NewConsumer([]string{appConfig.KafkaConfig.Address}, sarama.NewConfig())
	if err != nil {
		return err
	}

	defer func() {
		if err := consumer.Close(); err != nil {
			log.Fatalln(err)
		}
	}()

	partitionConsumer, err := consumer.ConsumePartition(appConfig.KafkaConfig.Topic, int32(appConfig.KafkaConfig.Partition), sarama.OffsetOldest)
	if err != nil {
		return err
	}

	kc := kafka.NewConsumer(partitionConsumer)

	err = kc.Run()
	if err != nil {
		return err
	}

	return nil
}
