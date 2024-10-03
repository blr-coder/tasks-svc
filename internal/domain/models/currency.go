package models

import "time"

type Currency string

const (
	CurrencyEUR Currency = "EUR"
	CurrencyUSD Currency = "USD"
	CurrencyPLN Currency = "PLN"
)

var DefaultCurrency = CurrencyEUR

func (c Currency) String() string {
	return string(c)
}

type CurrencyList []*Currency

func (l CurrencyList) ToStrings() []string {
	strings := make([]string, len(l))

	for i := range strings {
		strings[i] = l[i].String()
	}

	return strings
}

type CurrencyRate struct {
	Currency  string    `json:"currency" db:"currency"`
	RateEur   float64   `json:"rate_eur" db:"rate_eur"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
