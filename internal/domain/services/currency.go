package services

import (
	"context"
	"fmt"
	"github.com/blr-coder/tasks-svc/internal/domain/models"
	"github.com/blr-coder/tasks-svc/internal/infrastructure/storages/psql_store"
	"github.com/blr-coder/tasks-svc/pkg/currency_rates_client"
	"github.com/davecgh/go-spew/spew"
	"log"
)

type CurrencyService struct {
	currencyStorage     psql_store.ICurrencyStorage
	currencyRatesClient *currency_rates_client.Client
}

func NewCurrencyService(currencyStorage psql_store.ICurrencyStorage, currencyRatesClient *currency_rates_client.Client) *CurrencyService {
	return &CurrencyService{
		currencyStorage:     currencyStorage,
		currencyRatesClient: currencyRatesClient,
	}
}

func (s *CurrencyService) CheckRates(ctx context.Context) error {
	currencies, err := s.currencyStorage.ListCurrencyTickers(ctx)
	if err != nil {
		return fmt.Errorf("check rates err %w", err)
	}

	log.Println("<<<<<<<<<<<<<<<<<<<<<<<<<<<<<  DB CURRENCIES >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
	spew.Dump(currencies)
	log.Println("<<<<<<<<<<<<<<<<<<<<<<<<<<<<<  DB CURRENCIES >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")

	// Get rates from API
	rates, err := s.currencyRatesClient.GetRates(ctx, models.DefaultCurrency, currencies)
	if err != nil {
		return fmt.Errorf("get rates from ext api err %w", err)
	}

	log.Println("<<<<<<<<<<<<<<<<<<<<<<<<<<<<< CLIENT RATES >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
	spew.Dump(rates)
	log.Println("<<<<<<<<<<<<<<<<<<<<<<<<<<<<< CLIENT RATES >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")

	// Update Rates in DB
	err = s.currencyStorage.SetCurrencyRates(ctx, clientRatesToDBRates(rates))
	if err != nil {
		return fmt.Errorf("set rates err %w", err)
	}

	log.Println("rates were checked successfully")

	return nil
}

func clientRatesToDBRates(clientRates *currency_rates_client.ExchangeRatesInfo) []models.CurrencyRate {
	dbRates := make([]models.CurrencyRate, 0, len(clientRates.ExchangeRates))

	for k, v := range clientRates.ExchangeRates {
		dbRates = append(dbRates, models.CurrencyRate{
			Currency: k,
			RateEur:  v,
		})
	}

	return dbRates
}
