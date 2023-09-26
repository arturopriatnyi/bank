package banking

import (
	"errors"
	"math"
)

type Currency string

var (
	USD Currency = "USD"
	EUR Currency = "EUR"
	UAH Currency = "UAH"

	currencies = map[Currency]struct{}{
		USD: {},
		EUR: {},
		UAH: {},
	}
)

type Money struct {
	value    int32
	currency Currency
}

func NewMoney(units, nanos int32, currency Currency) (Money, error) {
	if nanos < 0 {
		return Money{}, errors.New("negative nanos")
	}
	if nanos > 99 {
		return Money{}, errors.New("too big nanos")
	}

	if _, ok := currencies[currency]; !ok {
		return Money{}, errors.New("invalid currency")
	}

	return Money{
		value:    units*100 + nanos,
		currency: currency,
	}, nil
}

func NewMoneyFromValue(v int32, currency Currency) Money {
	return Money{
		value:    v,
		currency: currency,
	}
}

func (m *Money) Units() int32 {
	return m.value / 100
}

func (m *Money) Nanos() int32 {
	return int32(math.Abs(float64(m.value))) % 100
}

func (m *Money) Currency() Currency {
	return m.currency
}

type Balance struct {
	UserID int32
	Money  []Money
}
