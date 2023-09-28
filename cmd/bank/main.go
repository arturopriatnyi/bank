package main

import (
	"bank/cmd/bank/config"
	"bank/cmd/bank/grpc"
	"bank/gen/proto/bank/v1/bankv1connect"
	"bank/pkg/banking"
	"bank/pkg/banking/balance"
	"bank/pkg/banking/exchange"
	"context"
	"database/sql"
	"fmt"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/sethvargo/go-envconfig"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

func main() {
	ctx := context.Background()

	l, err := zap.NewProduction()
	if err != nil {
		log.Panic(err)
	}
	defer l.Sync()

	l.Info("starting Bank service")

	l.Info("reading config")
	var cfg config.Config
	if err = envconfig.Process(ctx, &cfg); err != nil {
		l.Fatal("reading config error", zap.Error(err))
	}

	mgo, err := mongo.Connect(
		ctx,
		options.
			Client().
			ApplyURI(cfg.MongoDB.URL).
			SetServerAPIOptions(options.ServerAPI(options.ServerAPIVersion1)),
	)
	if err != nil {
		l.Fatal("mongo connection error", zap.Error(err))
	}

	if err = mgo.Database("admin").RunCommand(ctx, bson.D{{"ping", 1}}).Decode(&bson.M{}); err != nil {
		l.Fatal("mongo ping error", zap.Error(err))
	}
	l.Info("connected to MongoDB")

	pg, err := sql.Open(
		"postgres",
		fmt.Sprintf(
			"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			cfg.PostgreSQL.Host,
			cfg.PostgreSQL.Port,
			cfg.PostgreSQL.User,
			cfg.PostgreSQL.Password,
			cfg.PostgreSQL.Database,
			cfg.PostgreSQL.SSLMode,
		),
	)
	if err != nil {
		l.Fatal("postgresql connection error", zap.Error(err))
	}
	l.Info("connected to PostgreSQL")

	driver, err := postgres.WithInstance(pg, &postgres.Config{})
	if err != nil {
		l.Fatal("migration driver error", zap.Error(err))
	}
	migration, err := migrate.NewWithDatabaseInstance(
		"file://tools/postgresql/migrations",
		"postgres",
		driver,
	)
	if err != nil {
		l.Fatal("migration error", zap.Error(err))
	}
	if err := migration.Up(); err != nil && err != migrate.ErrNoChange {
		l.Fatal("migration error")
	}

	ratesLoader := exchange.NewRatesLoader(
		l,
		cfg.ExchangeRatesAPI.AccessKey,
		exchange.NewMongoDBRatesStorage(
			mgo.Database(cfg.MongoDB.Database),
		),
	)
	go ratesLoader.Start(ctx)
	l.Info("starting exchange rates loader")

	balancesLoader := balance.NewLoader(
		l,
		balance.NewPostgreSQLTransactionStorage(pg),
		balance.NewMongoDBBalancesStorage(
			mgo.Database(cfg.MongoDB.Database),
		),
	)
	go balancesLoader.Start(ctx)
	l.Info("starting balances loader")

	mux := http.NewServeMux()
	path, handler := bankv1connect.NewBankServiceHandler(
		grpc.NewBankHandler(
			banking.New(
				banking.NewPostgreSQLTransactionStorage(pg),
				banking.NewMongoDBBalanceStorage(
					mgo.Database(cfg.MongoDB.Database),
				),
			),
		),
	)
	mux.Handle(path, handler)

	server := &http.Server{
		Addr:    cfg.HTTPServer.Addr,
		Handler: h2c.NewHandler(mux, &http2.Server{}),
	}

	l.Info("starting http server", zap.String("addr", cfg.HTTPServer.Addr))
	go func() {
		if err = server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			l.Fatal("http server error", zap.Error(err))
		}
	}()

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt)
	<-done

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err = server.Shutdown(ctx); err != nil {
		l.Error("http server shut down error", zap.Error(err))
	}
	l.Info("http server shut down")

	if err = pg.Close(); err != nil {
		l.Error("postgresql closing error", zap.Error(err))
	}
	l.Info("postgresql closed")

	if err = mgo.Disconnect(ctx); err != nil {
		l.Error("mongodb closing error", zap.Error(err))
	}
	l.Info("mongodb closed")
}
