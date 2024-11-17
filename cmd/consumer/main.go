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

	db, err := sqlx.Open("postgres", appConfig.PostgresConnLink)
	if err != nil {
		return fmt.Errorf("connecting postgres: %w", err)
	}

	//r, err := kafka.NewReceiver(appConfig.KafkaConfig)
	r, err := kafka.NewGroupReceiver(appConfig.KafkaConfig)
	if err != nil {
		return err
	}

	// For test (test is OK)
	/*err = r.Receive(ctx)
	if err != nil {
		return err
	}*/

	ar := psql_store.NewTaskActionPsqlStorage(db)

	as := services.NewTaskActionService(ar)
	tah := handlers.NewTaskActionHandler(r, as)

	err = tah.Handle(ctx)
	if err != nil {
		return err
	}

	return nil
}
