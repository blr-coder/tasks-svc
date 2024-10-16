package main

import (
	"context"
	"fmt"
	"github.com/blr-coder/tasks-svc/internal/config"
	grpc "github.com/blr-coder/tasks-svc/internal/delivery/grpc_api"
	"github.com/blr-coder/tasks-svc/internal/delivery/grpc_api/handlers"
	"github.com/blr-coder/tasks-svc/internal/domain/services"
	"github.com/blr-coder/tasks-svc/internal/infrastructure/queues/kafka"
	"github.com/blr-coder/tasks-svc/internal/infrastructure/storages/psql_store"
	"github.com/blr-coder/tasks-svc/internal/infrastructure/storages/transaction"
	"github.com/bufbuild/protovalidate-go"
	"github.com/jmoiron/sqlx"
	"log"
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
	log.Println("RUN TASK API")

	ctx := context.Background()

	ctx, cancel := signal.NotifyContext(ctx, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	appConfig, err := config.NewAppConfig()
	if err != nil {
		return err
	}

	db, err := sqlx.Open("postgres", appConfig.PostgresConnLink)
	if err != nil {
		return fmt.Errorf("connecting postgres: %w", err)
	}

	validator, err := protovalidate.New()
	if err != nil {
		return fmt.Errorf("creating protovalidator: %w", err)
	}

	taskPsqlStorage := psql_store.NewTaskPsqlStorage(db)
	currencyPsqlStorage := psql_store.NewCurrencyPsqlStorage(db)

	sender, err := kafka.NewKafkaSender(appConfig.KafkaConfig)
	if err != nil {
		return err
	}

	txManager := transaction.NewTransactionManager(db)

	taskService := services.NewTaskService(taskPsqlStorage, currencyPsqlStorage, sender, txManager)
	taskGRPCServer := handlers.NewTaskServiceServer(validator, taskService)

	grpcServer := grpc.NewGRPCServer(taskGRPCServer)

	go func() {
		err = grpcServer.Run(appConfig.AppPort)
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
