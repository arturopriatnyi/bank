package banking

type ExchangeRate struct {
	value        float64
	baseCurrency string
}

type ExchangeRatesClient interface {
}

type ExchangeRatesStorage interface{}

type ExchangeRatesLoader struct {
	client  ExchangeRatesClient
	storage ExchangeRatesStorage
}

func NewExchangeRatesLoader(client ExchangeRatesClient, storage ExchangeRatesStorage) *ExchangeRatesLoader {
	return &ExchangeRatesLoader{
		client:  client,
		storage: storage,
	}
}
