package psql_store

import (
	"context"
	"database/sql"
	"errors"
	"github.com/blr-coder/tasks-svc/internal/domain/errs"
	"github.com/blr-coder/tasks-svc/internal/domain/models"
	"github.com/jmoiron/sqlx"
)

type ICurrencyStorage interface {
	GetRateByEUR(ctx context.Context, currency models.Currency) (*models.CurrencyRate, error)
	ListCurrencyTitles(ctx context.Context) ([]*models.Currency, error)
}

type CurrencyPsqlStorage struct {
	db *sqlx.DB
}

func NewCurrencyPsqlStorage(database *sqlx.DB) *CurrencyPsqlStorage {
	return &CurrencyPsqlStorage{
		db: database,
	}
}

func (cs *CurrencyPsqlStorage) GetRateByEUR(ctx context.Context, currency models.Currency) (*models.CurrencyRate, error) {
	rate := &models.CurrencyRate{}

	query := `
		SELECT
			*
		FROM currency_rates
		WHERE
			currency = $1
	`

	err := cs.db.GetContext(ctx, rate, query, currency.String())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.NewDomainNotFoundError().WithParam("currency", currency.String())
		}

		return nil, err
	}

	return rate, nil
}

func (cs *CurrencyPsqlStorage) ListCurrencyTitles(ctx context.Context) ([]*models.Currency, error) {
	var currencies []*models.Currency

	query := `
		SELECT
			currency
		FROM currency_rates
	`

	err := cs.db.SelectContext(ctx, &currencies, query)
	if err != nil {
		return nil, err
	}

	return currencies, nil
}
