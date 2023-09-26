package exchange

import (
	"context"

	"bank/pkg/banking"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDBRatesStorage struct {
	collection *mongo.Collection
}

func NewMongoDBRatesStorage(db *mongo.Database) *MongoDBRatesStorage {
	return &MongoDBRatesStorage{
		collection: db.Collection("exchange_rates"),
	}
}

func (s *MongoDBRatesStorage) Set(ctx context.Context, currency banking.Currency, rate banking.Money) error {
	type exchangeRateRecord struct {
		Currency     string `bson:"currency"`
		RateUnits    int32  `bson:"rate_units"`
		RateNanos    int32  `bson:"rate_nanos"`
		RateCurrency string `bson:"rate_currency"`
	}

	opts := options.Update().SetUpsert(true)
	filter := bson.M{"currency": currency}
	update := bson.M{"$set": &exchangeRateRecord{
		Currency:     string(currency),
		RateUnits:    rate.Units(),
		RateNanos:    rate.Nanos(),
		RateCurrency: string(rate.Currency()),
	}}

	_, err := s.collection.UpdateOne(ctx, filter, update, opts)

	return err
}
