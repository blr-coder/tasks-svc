package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/blr-coder/tasks-svc/internal/config"
	"github.com/blr-coder/tasks-svc/internal/infrastructure/storages/psql_store"
	"github.com/blr-coder/tasks-svc/pkg/currency_rates"
	"github.com/davecgh/go-spew/spew"
	"github.com/jmoiron/sqlx"
	"log"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	err := runCheckRates()
	if err != nil {
		log.Fatal(err)
	}
}

func runCheckRates() error {
	log.Println("RUN WORKER!")

	ctx := context.Background()

	ctx, cancel := signal.NotifyContext(ctx, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	appConfig, err := config.NewAppConfig()
	if err != nil {
		return err
	}

	spew.Dump(appConfig)

	db, err := sqlx.Open("postgres", appConfig.PostgresConnLink)
	if err != nil {
		return fmt.Errorf("connecting postgres: %w", err)
	}

	currencyPsqlStorage := psql_store.NewCurrencyPsqlStorage(db)

	w := NewCurrencyWorker(appConfig.AbstractAPIConfig, currencyPsqlStorage)

	err = w.Run(ctx)
	if err != nil {
		//TODO: Handle err
		return err
	}

	return nil
}

type CurrencyRatesClient interface {
	Close() error
	GetRates(ctx context.Context, from string, currenciesTo []string) (*currency_rates.ExchangeRatesInfo, error)
}

type CurrencyWorker struct {
	tickDuration        time.Duration
	currencyStorage     psql_store.ICurrencyStorage
	currencyRatesClient CurrencyRatesClient
}

func NewCurrencyWorker(apiConfig config.AbstractAPIConfig, currencyStorage psql_store.ICurrencyStorage) *CurrencyWorker {
	currencyRatesClient := currency_rates.NewClient(apiConfig.APIKey)

	return &CurrencyWorker{
		tickDuration:        apiConfig.TickerInterval,
		currencyStorage:     currencyStorage,
		currencyRatesClient: currencyRatesClient,
	}
}

func (w *CurrencyWorker) Run(ctx context.Context) error {
	ticker := time.NewTicker(w.tickDuration)
	defer ticker.Stop()

	for true {
		select {
		case <-ctx.Done():
			// TODO: XZ its OK?
			return errors.New("context die by timeout")
		case <-ticker.C:
			fmt.Println("TICKER TICK", time.Now().Format(time.RFC3339))
			log.Println("RUN WORKER!!!!")

			rates, err := w.currencyRatesClient.GetRates(ctx, "EUR", []string{"PLN"})
			if err != nil {
				return err
			}

			spew.Dump(rates)
		}
	}

	return nil
}
