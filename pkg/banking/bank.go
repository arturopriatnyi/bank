//go:generate mockery --name=TransactionStorage --inpackage --case=underscore --testonly
//go:generate mockery --name=BalanceStorage --inpackage --case=underscore --testonly
package banking

import (
	"context"
)

type TransactionStorage interface {
	Set(ctx context.Context, t Transaction) error
}

type BalanceStorage interface {
	Get(ctx context.Context, userID int32) (Balance, error)
}

type Bank struct {
	transactions TransactionStorage
	balances     BalanceStorage
}

func New(transactions TransactionStorage, balances BalanceStorage) *Bank {
	return &Bank{
		transactions: transactions,
		balances:     balances,
	}
}

func (b *Bank) MakeTransaction(ctx context.Context, senderID, receiverID int32, money Money) (Transaction, error) {
	transaction := NewTransaction(senderID, receiverID, money)

	if err := b.transactions.Set(ctx, transaction); err != nil {
		return Transaction{}, err
	}

	return transaction, nil
}

func (b *Bank) Balance(ctx context.Context, userID int32) (Balance, error) {
	balance, err := b.balances.Get(ctx, userID)
	if err != nil {
		return Balance{}, err

	}

	return balance, nil
}
