package exchange

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"bank/pkg/banking"

	"go.uber.org/zap"
)

const ratesAPIURL = "http://api.exchangeratesapi.io/v1/latest?access_key="

type Rates map[banking.Currency]banking.Money

type RatesStorage interface {
	Set(ctx context.Context, currency banking.Currency, money banking.Money) error
}

type RatesLoader struct {
	l          *zap.Logger
	httpClient *http.Client
	apiKey     string
	storage    RatesStorage
}

func NewRatesLoader(l *zap.Logger, apiKey string, storage RatesStorage) *RatesLoader {
	return &RatesLoader{
		l:          l,
		httpClient: http.DefaultClient,
		apiKey:     apiKey,
		storage:    storage,
	}
}

func (c *RatesLoader) Start(ctx context.Context) {
	t := time.NewTicker(time.Minute)

	for {
		rates, err := c.load(ctx)
		if err != nil {
			c.l.Error("exchange rates loading error", zap.Error(err))
		}

		for cur, price := range rates {
			if err = c.storage.Set(ctx, cur, price); err != nil {
				c.l.Error("exchange rates storing error", zap.Error(err))
			}
		}

		c.l.Info("exchange rates loaded")

		select {
		case <-t.C:
		case <-ctx.Done():
			c.l.Info("exchange rates loader shut down")

			return
		}

	}
}

func (c *RatesLoader) load(ctx context.Context) (Rates, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, ratesAPIURL+c.apiKey, nil)
	if err != nil {
		return nil, err
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	type responseBody struct {
		Success bool                         `json:"success"`
		Base    banking.Currency             `json:"base"`
		Rates   map[banking.Currency]float64 `json:"rates"`
	}
	var body responseBody
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		return nil, err
	}

	if body.Success != true {
		return nil, errors.New("invalid API request")
	}

	rates := make(Rates)
	for currency := range body.Rates {
		currencyRateToUSD := (1 / body.Rates[banking.USD]) * body.Rates[currency]

		rates[currency], err = banking.NewMoney(
			int32(currencyRateToUSD),
			int32(currencyRateToUSD*100)%100,
			banking.USD,
		)
		if err != nil {
			return nil, err
		}
	}

	return rates, nil
}
