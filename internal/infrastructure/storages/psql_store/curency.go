package psql_store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/blr-coder/tasks-svc/internal/domain/errs"
	"github.com/blr-coder/tasks-svc/internal/domain/models"
	"github.com/jmoiron/sqlx"
)

type ICurrencyStorage interface {
	GetRateByEUR(ctx context.Context, currency models.Currency) (*models.CurrencyRate, error)
	ListCurrencyTickers(ctx context.Context) (models.CurrencyList, error)
	SetCurrencyRates(ctx context.Context, rates []models.CurrencyRate) error
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

func (cs *CurrencyPsqlStorage) ListCurrencyTickers(ctx context.Context) (models.CurrencyList, error) {
	var currencies models.CurrencyList

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

func (cs *CurrencyPsqlStorage) SetCurrencyRates(ctx context.Context, rates []models.CurrencyRate) error {
	// bulk/batch insert is supported by SQLX v. 1.3.0+ - https://github.com/jmoiron/sqlx/blob/master/README.md
	// The only restriction - no way to scan RETURNING * to the result's slice
	query := `
		INSERT INTO currency_rates (currency, rate_eur) 
		VALUES (:currency, :rate_eur)
		ON CONFLICT (currency)
		DO UPDATE
		SET
		    rate_eur = EXCLUDED.rate_eur,
		    updated_at = CURRENT_TIMESTAMP AT TIME ZONE 'UTC'
		    `

	_, err := cs.db.NamedExecContext(ctx, cs.db.Rebind(query), rates)
	if err != nil {
		return fmt.Errorf("set currency rates in storage err: %w", err)
	}

	return nil
}
