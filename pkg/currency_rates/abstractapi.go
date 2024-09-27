package currency_rates

import (
	"context"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"log"
	"net/http"
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
	ExchangeRates map[string]float64 `json:"currency_rates"`
}

func (c *Client) GetRates(ctx context.Context, from string, currenciesTo []string) (*ExchangeRatesInfo, error) {

	logrus.Info("GetRates for CURRENCY_FROM:", from)
	logrus.Info("GetRates for CURRENCIES_TO:", currenciesTo)

	/*url := fmt.Sprintf(
		abstractAPIFormat,
		c.apiToken,
		from,
		strings.Join(currenciesTo, ","),
	)
	request, err := http.NewRequest("GET", url, nil)*/
	//if err != nil {
	//	return nil, err
	//}
	/*response, err := c.http.Do(request)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = response.Body.Close()
		if err != nil {
			log.Println(err)
		}
	}()*/

	resp, err := http.Get("https://exchange-rates.abstractapi.com/v1/live/?api_key=c968d364db5b489f96b7bd85ef0da1e2&base=EUR&target=PLN")
	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(string(body))

	exchangeInfo := &ExchangeRatesInfo{}

	//res := json.NewDecoder(body).Decode(exchangeInfo)

	//return exchangeInfo, json.NewDecoder(response.Body).Decode(exchangeInfo)
	return exchangeInfo, nil
}
