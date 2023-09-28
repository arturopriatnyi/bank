package balance

import (
	"context"
	"time"

	"bank/pkg/banking"

	"go.uber.org/zap"
)

type TransactionStorage interface {
	GetBalances(ctx context.Context) ([]banking.Balance, error)
}

type BalancesStorage interface {
	Set(ctx context.Context, balance banking.Balance) error
}

type Loader struct {
	l            *zap.Logger
	transactions TransactionStorage
	balances     BalancesStorage
}

func NewLoader(l *zap.Logger, transactions TransactionStorage, balances BalancesStorage) *Loader {
	return &Loader{
		l:            l,
		transactions: transactions,
		balances:     balances,
	}
}

func (l *Loader) Start(ctx context.Context) {
	t := time.NewTicker(time.Minute)

	for {
		balances, err := l.transactions.GetBalances(ctx)
		if err != nil {
			l.l.Error("balances loading error", zap.Error(err))
		}

		for _, balance := range balances {
			if err = l.balances.Set(ctx, balance); err != nil {
				l.l.Error("balances storing error", zap.Error(err))
			}
		}

		l.l.Info("balances loaded")

		select {
		case <-t.C:
		case <-ctx.Done():
			l.l.Info("balances loader shut down")

			return
		}
	}
}
