package main

import (
	"fmt"
	"github.com/blr-coder/tasks-svc/internal/config"
	grpc "github.com/blr-coder/tasks-svc/internal/delivery/grpc_api"
	"github.com/blr-coder/tasks-svc/internal/domain/services"
	"github.com/blr-coder/tasks-svc/internal/storage"
	"github.com/jmoiron/sqlx"
	"log"
	"net"
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

	// Test kafka
	/*consumer, err := sarama.NewConsumer([]string{appConfig.KafkaConfig.Address}, sarama.NewConfig())
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
	}*/
	// Test kafka

	db, err := sqlx.Open("postgres", appConfig.PostgresConnLink)
	if err != nil {
		return err
	}

	taskPsqlStorage := storage.NewTaskPsqlStorage(db)
	taskService := services.NewTaskService(taskPsqlStorage)
	taskGRPCServer := grpc.NewTaskServiceServer(taskService)

	grpcServer := grpc.NewGRPCServer(taskGRPCServer)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", appConfig.AppPort))
	if err != nil {
		return err
	}

	return grpcServer.Serve(listener)
}
