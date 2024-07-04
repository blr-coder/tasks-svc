package models

type Currency string

const (
	CurrencyEUR Currency = "EUR"
	CurrencyUSD Currency = "USD"
	CurrencyPLN Currency = "PLN"
)

func (c Currency) String() string {
	return string(c)
}
