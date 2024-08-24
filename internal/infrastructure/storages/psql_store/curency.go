package psql_store

import (
	"context"
	"github.com/jmoiron/sqlx"
)

type ICurrencyStorage interface {
	GetRate(ctx context.Context, currency string) (float64, error)
}

type CurrencyPsqlStorage struct {
	db *sqlx.DB
}

func NewCurrencyPsqlStorage(database *sqlx.DB) *CurrencyPsqlStorage {
	return &CurrencyPsqlStorage{
		db: database,
	}
}

func (cs *CurrencyPsqlStorage) GetRate(ctx context.Context, currency string) (float64, error) {

	return 0, nil
}
