package main

import (
	"context"
	"fmt"
	"github.com/blr-coder/tasks-svc/internal/config"
	grpc "github.com/blr-coder/tasks-svc/internal/delivery/grpc_api"
	"github.com/blr-coder/tasks-svc/internal/domain/services"
	"github.com/blr-coder/tasks-svc/internal/infrastructure/queues/kafka"
	"github.com/blr-coder/tasks-svc/internal/infrastructure/storages/psql_store"
	"github.com/bufbuild/protovalidate-go"
	"github.com/jmoiron/sqlx"
	"log"
	"net"
	"os/signal"
	"syscall"
)

func main() {
	err := runApp()
	if err != nil {
		log.Fatal(err)
	}
}

func runApp() error {
	log.Println("ASYNC")

	ctx := context.Background()

	ctx, cancel := signal.NotifyContext(ctx, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

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
		return fmt.Errorf("connecting postgres: %w", err)
	}

	validator, err := protovalidate.New()
	if err != nil {
		return fmt.Errorf("creating protovalidator: %w", err)
	}

	taskPsqlStorage := psql_store.NewTaskPsqlStorage(db)

	sender, err := kafka.NewKafkaSender(appConfig.KafkaConfig)
	if err != nil {
		return err
	}

	taskService := services.NewTaskService(taskPsqlStorage, sender)
	taskGRPCServer := grpc.NewTaskServiceServer(validator, taskService)

	grpcServer := grpc.NewGRPCServer(taskGRPCServer)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", appConfig.AppPort))
	if err != nil {
		return err
	}

	go func() {
		err = grpcServer.GRPCServer.Serve(listener)
		if err != nil {
			cancel()
			return
		}
	}()

	// Graceful shutdown
	<-ctx.Done()

	grpcServer.GRPCServer.GracefulStop()
	log.Println("\nGracefully stopped")

	return err
}
