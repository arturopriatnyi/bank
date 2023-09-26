package banking

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNew(t *testing.T) {
	transactions := NewMockTransactionStorage(t)
	balances := NewMockBalanceStorage(t)

	bank := New(transactions, balances)

	assert.Equal(t,
		&Bank{transactions: transactions, balances: balances},
		bank,
	)
}

func TestBank_MakeTransaction(t *testing.T) {
	var senderID, receiverID int32 = 1, 2
	money, err := NewMoney(5, 0, USD)
	assert.NoError(t, err)

	for name, tc := range map[string]struct {
		transactions                 func(*testing.T) TransactionStorage
		wantSenderID, wantReceiverID int32
		wantErr                      bool
	}{
		"TransactionIsMade": {
			transactions: func(t *testing.T) TransactionStorage {
				transactions := NewMockTransactionStorage(t)

				transactions.
					On(
						"Set",
						context.TODO(),
						mock.AnythingOfType("Transaction"),
					).
					Return(nil)

				return transactions
			},
			wantSenderID:   1,
			wantReceiverID: 2,
			wantErr:        false,
		},
		"TransactionStorageSetError": {
			transactions: func(t *testing.T) TransactionStorage {
				transactions := NewMockTransactionStorage(t)

				transactions.On(
					"Set",
					context.TODO(),
					mock.AnythingOfType("Transaction"),
				).
					Return(errors.New("test error"))

				return transactions
			},
			wantSenderID:   0,
			wantReceiverID: 0,
			wantErr:        true,
		},
	} {
		t.Run(name, func(t *testing.T) {
			bank := &Bank{transactions: tc.transactions(t)}

			ts, err := bank.MakeTransaction(context.TODO(), senderID, receiverID, money)

			assert.Equal(t, tc.wantErr, err != nil)
			assert.Equal(t, tc.wantSenderID, ts.SenderID)
			assert.Equal(t, tc.wantReceiverID, ts.ReceiverID)
		})
	}
}

func TestBank_Balance(t *testing.T) {
	var userID int32 = 1
	money, err := NewMoney(5, 0, USD)
	assert.NoError(t, err)
	balance := Balance{
		UserID: userID,
		Money:  []Money{money},
	}

	for name, tc := range map[string]struct {
		balances    func(*testing.T) BalanceStorage
		wantBalance Balance
		wantErr     bool
	}{
		"BalanceIsRetrieved": {
			balances: func(t *testing.T) BalanceStorage {
				balances := NewMockBalanceStorage(t)

				balances.
					On("Get", context.TODO(), userID).
					Return(balance, nil)

				return balances
			},
			wantBalance: balance,
			wantErr:     false,
		},
		"BalanceStorageGetError": {
			balances: func(t *testing.T) BalanceStorage {
				balances := NewMockBalanceStorage(t)

				balances.On(
					"Get", context.TODO(), userID).
					Return(Balance{}, errors.New("test error"))

				return balances
			},
			wantBalance: Balance{},
			wantErr:     true,
		},
	} {
		t.Run(name, func(t *testing.T) {
			bank := &Bank{balances: tc.balances(t)}

			balance, err := bank.Balance(context.TODO(), userID)

			assert.Equal(t, tc.wantErr, err != nil)
			assert.Equal(t, tc.wantBalance, balance)
		})
	}
}
