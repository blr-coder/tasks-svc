package currency_rates_client

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/blr-coder/tasks-svc/internal/domain/models"
	"log"
	"net/http"
	"strings"
)

const abstractAPIFormat = "https://exchange-rates.abstractapi.com/v1/live/?api_key=%s&base=%s&target=%s"

type Client struct {
	http     *http.Client
	apiToken string
}

func (c *Client) Close() error {
	return nil
}

func NewClient(apiToken string) *Client {
	return &Client{
		http:     &http.Client{},
		apiToken: apiToken,
	}
}

type ExchangeRatesInfo struct {
	Base          string             `json:"base"`
	ExchangeRates map[string]float64 `json:"exchange_rates"`
}

func (c *Client) GetRates(ctx context.Context, from models.Currency, currenciesTo models.CurrencyList) (exchangeInfo *ExchangeRatesInfo, err error) {
	log.Println("GetRates for CURRENCY_FROM:", from)
	log.Println("GetRates for CURRENCIES_TO:", currenciesTo)

	url := fmt.Sprintf(
		abstractAPIFormat,
		c.apiToken,
		from.String(),
		strings.Join(currenciesTo.ToStrings(), ","),
	)

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	response, err := c.http.Do(request)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = response.Body.Close()
		if err != nil {
			log.Println(err)
		}
	}()

	err = json.NewDecoder(response.Body).Decode(&exchangeInfo)
	if err != nil {
		return nil, err
	}

	return exchangeInfo, nil
}
