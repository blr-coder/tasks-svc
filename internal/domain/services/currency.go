package services

import (
	"context"
	"fmt"
	"github.com/blr-coder/tasks-svc/internal/domain/models"
	"github.com/blr-coder/tasks-svc/internal/infrastructure/storages/psql_store"
	"github.com/blr-coder/tasks-svc/pkg/currency_rates"
	"github.com/davecgh/go-spew/spew"
)

type CurrencyService struct {
	currencyStorage     psql_store.ICurrencyStorage
	currencyRatesClient *currency_rates.Client
}

func NewCurrencyService(currencyStorage psql_store.ICurrencyStorage, currencyRatesClient *currency_rates.Client) *CurrencyService {
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

	spew.Dump(currencies)

	// GetRates
	rates, err := s.currencyRatesClient.GetRates(ctx, models.DefaultCurrency, currencies)
	if err != nil {
		return err
	}

	spew.Dump(rates)

	// UpdateRates

	return nil
}
