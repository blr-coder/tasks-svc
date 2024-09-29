package services

import (
	"context"
	"fmt"
	"github.com/blr-coder/tasks-svc/internal/infrastructure/storages/psql_store"
	"github.com/davecgh/go-spew/spew"
)

type CurrencyService struct {
	currencyStorage psql_store.ICurrencyStorage
}

func NewCurrencyService(currencyStorage psql_store.ICurrencyStorage) *CurrencyService {
	return &CurrencyService{
		currencyStorage: currencyStorage,
	}
}

func (s *CurrencyService) CheckRates(ctx context.Context) error {
	currencies, err := s.currencyStorage.ListCurrencyTitles(ctx)
	if err != nil {
		return fmt.Errorf("check rates err %w", err)
	}

	spew.Dump(currencies)

	return nil
}
