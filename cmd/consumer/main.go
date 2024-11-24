package main

import (
	"context"
	"fmt"
	"github.com/blr-coder/tasks-svc/internal/config"
	"github.com/blr-coder/tasks-svc/internal/delivery/kafka/handlers"
	"github.com/blr-coder/tasks-svc/internal/domain/services"
	"github.com/blr-coder/tasks-svc/internal/infrastructure/queues/kafka"
	"github.com/blr-coder/tasks-svc/internal/infrastructure/storages/psql_store"
	"github.com/jmoiron/sqlx"
	"log"
	"os/signal"
	"syscall"
)

func main() {
	err := startTaskActionConsumer()
	if err != nil {
		log.Fatal(err)
	}
}

func startTaskActionConsumer() error {
	log.Println("Starting Task Action Consumer")

	ctx := context.Background()

	ctx, cancel := signal.NotifyContext(ctx, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	appConfig, err := config.NewAppConfig()
	if err != nil {
		return fmt.Errorf("failed to load app config: %w", err)
	}

	db, err := sqlx.Open("postgres", appConfig.PostgresConnLink)
	if err != nil {
		return fmt.Errorf("failed to connect to PostgreSQL: %w", err)
	}
	defer db.Close()

	//kafkaReceiver, err := kafka.NewReceiver(appConfig.KafkaConfig)
	kafkaReceiver, err := kafka.NewGroupReceiver(appConfig.KafkaConfig)
	if err != nil {
		return err
	}

	// For test (test is OK)
	/*err = r.Receive(ctx)
	if err != nil {
		return err
	}*/

	taskActionStorage := psql_store.NewTaskActionPsqlStorage(db)
	taskActionService := services.NewTaskActionService(taskActionStorage)
	taskActionHandler := handlers.NewTaskActionHandler(kafkaReceiver, taskActionService)

	if err = taskActionHandler.Handle(ctx); err != nil {
		return fmt.Errorf("failed to handle Kafka messages: %w", err)
	}

	log.Println("Task Action Consumer stopped gracefully")
	return nil
}
