package banking

import (
	"context"
	"database/sql"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type PostgreSQLTransactionStorage struct {
	db *sql.DB
}

func NewPostgreSQLTransactionStorage(db *sql.DB) *PostgreSQLTransactionStorage {
	return &PostgreSQLTransactionStorage{db: db}
}

const setQuery = `
	INSERT INTO transactions (id, sender_id, receiver_id, money_value, money_currency)
	VALUES ($1, $2, $3, $4, $5);
`

func (s *PostgreSQLTransactionStorage) Set(ctx context.Context, transaction Transaction) error {
	_, err := s.db.ExecContext(
		ctx,
		setQuery,
		transaction.ID,
		transaction.SenderID,
		transaction.ReceiverID,
		transaction.Money.value,
		transaction.Money.Currency(),
	)

	return err
}

type MongoDBBalanceStorage struct {
	collection *mongo.Collection
}

func NewMongoDBBalanceStorage(db *mongo.Database) *MongoDBBalanceStorage {
	return &MongoDBBalanceStorage{
		collection: db.Collection("balances"),
	}
}

func (s *MongoDBBalanceStorage) Get(ctx context.Context, userID int32) (Balance, error) {
	type moneyRecord struct {
		Units    int32  `bson:"units"`
		Nanos    int32  `bson:"nanos"`
		Currency string `bson:"currency"`
	}

	type balanceRecord struct {
		UserID int32 `bson:"user_id"`
		Money  []moneyRecord
	}

	var record balanceRecord
	if err := s.collection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&record); err != nil {
		return Balance{}, err
	}

	balance := Balance{UserID: record.UserID}
	for _, money := range record.Money {
		m, err := NewMoney(money.Units, money.Nanos, Currency(money.Currency))
		if err != nil {
			return Balance{}, err
		}

		balance.Money = append(balance.Money, m)
	}

	return balance, nil
}
