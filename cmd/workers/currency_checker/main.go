package main

import (
	"context"
	"fmt"
	"github.com/blr-coder/tasks-svc/internal/config"
	"github.com/blr-coder/tasks-svc/internal/domain/services"
	"github.com/blr-coder/tasks-svc/internal/infrastructure/storages/psql_store"
	"github.com/blr-coder/tasks-svc/pkg/currency_rates_client"
	"github.com/blr-coder/tasks-svc/pkg/worker"
	"github.com/davecgh/go-spew/spew"
	"github.com/jmoiron/sqlx"
	"log"
	"os/signal"
	"syscall"
)

func main() {
	err := runCheckRates()
	if err != nil {
		log.Fatal(err)
	}
}

func runCheckRates() error {
	log.Println("RUN CURRENCY CHECKER!")

	_, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	appConfig, err := config.NewAppConfig()
	if err != nil {
		return err
	}

	log.Println("<<<<<<<<<<<<<<<<<<<<<<<<<<<<<  APP CONFIG >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
	spew.Dump(appConfig)
	log.Println("<<<<<<<<<<<<<<<<<<<<<<<<<<<<<  APP CONFIG >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")

	db, err := sqlx.Open("postgres", appConfig.PostgresConnLink)
	if err != nil {
		return fmt.Errorf("connecting postgres: %w", err)
	}

	currencyService := services.NewCurrencyService(
		psql_store.NewCurrencyPsqlStorage(db),
		currency_rates_client.NewClient(appConfig.AbstractAPIConfig.APIKey),
	)

	w := worker.NewWorker(appConfig.AbstractAPIConfig.TickerInterval)

	w.Run(currencyService.CheckRates)

	return nil
}
