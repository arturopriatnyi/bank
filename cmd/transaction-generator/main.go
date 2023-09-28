package main

import (
	"bank/pkg/banking"
	"context"
	"log"
	"math/rand"
	"net/http"

	"bank/cmd/bank/config"
	bankv1 "bank/gen/proto/bank/v1"
	"bank/gen/proto/bank/v1/bankv1connect"

	"connectrpc.com/connect"
	"github.com/sethvargo/go-envconfig"
	"go.uber.org/zap"
)

func main() {
	ctx := context.Background()

	l, err := zap.NewProduction()
	if err != nil {
		log.Panic(err)
	}
	defer l.Sync()

	l.Info("starting Transaction Generator")

	l.Info("reading config")
	var cfg config.Config
	if err = envconfig.Process(ctx, &cfg); err != nil {
		l.Fatal("reading config error", zap.Error(err))
	}

	client := bankv1connect.NewBankServiceClient(http.DefaultClient, "http://"+cfg.HTTPServer.Addr)

	currencies := []banking.Currency{banking.USD, banking.EUR, banking.UAH}
	for i := 0; i < 100; i++ {
		senderID := 1 + rand.Int31n(10)
		receiverID := 1 + rand.Int31n(10)

		units := rand.Int31n(101)
		nanos := rand.Int31n(100)
		currency := currencies[rand.Int31n(3)]

		_, err = client.MakeTransaction(
			ctx,
			connect.NewRequest(
				&bankv1.MakeTransactionRequest{
					SenderId:   senderID,
					ReceiverId: receiverID,
					Money: &bankv1.Money{
						Units:    units,
						Nanos:    nanos,
						Currency: string(currency),
					},
				},
			),
		)
		if err != nil {
			l.Error("making transaction error", zap.Error(err))
		}
	}

	l.Info("generated 100 transactions")
}
