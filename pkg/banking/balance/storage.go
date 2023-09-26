package balance

import (
	"context"
	"database/sql"
	"go.mongodb.org/mongo-driver/mongo/options"

	"bank/pkg/banking"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type PostgreSQLTransactionStorage struct {
	db *sql.DB
}

func NewPostgreSQLTransactionStorage(db *sql.DB) *PostgreSQLTransactionStorage {
	return &PostgreSQLTransactionStorage{db: db}
}

const getBalancesByCurrencyQuery = `
SELECT user_id, money_currency, SUM(balance) AS total_balance
FROM (
    SELECT sender_id AS user_id, money_currency, -money_value AS balance FROM transactions
    UNION ALL
    SELECT receiver_id AS user_id, money_currency, money_value AS balance FROM transactions
) AS user_balances
GROUP BY user_id, money_currency;
`

func (s *PostgreSQLTransactionStorage) GetBalances(ctx context.Context) ([]banking.Balance, error) {
	rows, err := s.db.QueryContext(ctx, getBalancesByCurrencyQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var (
		userID   int32
		currency string
		balance  int32
	)

	balancesMap := make(map[int32][]banking.Money)
	for rows.Next() {
		if err := rows.Scan(&userID, &currency, &balance); err != nil {
			return nil, err
		}

		money := banking.NewMoneyFromValue(balance, banking.Currency(currency))
		balancesMap[userID] = append(balancesMap[userID], money)
	}
	if rows.Err() != nil {
		return nil, err
	}

	var balances []banking.Balance
	for userID, money := range balancesMap {
		balances = append(balances, banking.Balance{
			UserID: userID,
			Money:  money,
		})

	}

	return balances, nil
}

type MongoDBBalancesStorage struct {
	collection *mongo.Collection
}

func NewMongoDBBalancesStorage(db *mongo.Database) *MongoDBBalancesStorage {
	return &MongoDBBalancesStorage{
		collection: db.Collection("balances"),
	}
}

func (s *MongoDBBalancesStorage) Set(ctx context.Context, balance banking.Balance) error {
	filter := bson.M{"user_id": balance.UserID}

	var money []bson.M
	for _, m := range balance.Money {
		money = append(money, bson.M{
			"units":    m.Units(),
			"nanos":    m.Nanos(),
			"currency": m.Currency(),
		})

	}
	update := bson.M{
		"$set": bson.M{
			"user_id": balance.UserID,
			"money":   money,
		},
	}

	opts := options.Update().SetUpsert(true)
	_, err := s.collection.UpdateOne(ctx, filter, update, opts)

	return err
}
