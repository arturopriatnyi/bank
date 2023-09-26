package grpc

import (
	bankv1 "bank/gen/proto/bank/v1"
	"bank/pkg/banking"
	"context"

	"bank/gen/proto/bank/v1/bankv1connect"
	"connectrpc.com/connect"
)

type BankHandler struct {
	bankv1connect.UnimplementedBankServiceHandler

	bank *banking.Bank
}

func NewBankHandler(bank *banking.Bank) *BankHandler {
	return &BankHandler{bank: bank}
}

func (h *BankHandler) MakeTransaction(ctx context.Context, req *connect.Request[bankv1.MakeTransactionRequest]) (*connect.Response[bankv1.MakeTransactionResponse], error) {
	money, err := banking.NewMoney(req.Msg.Money.Units, req.Msg.Money.Nanos, banking.Currency(req.Msg.Money.Currency))
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	transaction, err := h.bank.MakeTransaction(ctx, req.Msg.SenderId, req.Msg.ReceiverId, money)

	return connect.NewResponse(&bankv1.MakeTransactionResponse{
		Transaction: &bankv1.Transaction{
			Id:         transaction.ID,
			SenderId:   transaction.SenderID,
			ReceiverId: transaction.ReceiverID,
			Money: &bankv1.Money{
				Units:    transaction.Money.Units(),
				Nanos:    transaction.Money.Nanos(),
				Currency: string(transaction.Money.Currency()),
			},
		},
	}), err
}

func (h *BankHandler) Balance(ctx context.Context, req *connect.Request[bankv1.BalanceRequest]) (*connect.Response[bankv1.BalanceResponse], error) {
	balance, err := h.bank.Balance(ctx, req.Msg.UserId)
	if err != nil {
		return nil, connect.NewError(connect.CodeNotFound, err)
	}

	var balances []*bankv1.Money
	for _, m := range balance.Money {
		balances = append(balances, &bankv1.Money{
			Units:    m.Units(),
			Nanos:    m.Nanos(),
			Currency: string(m.Currency()),
		})

	}

	return connect.NewResponse(&bankv1.BalanceResponse{
		Balances: balances,
	}), nil
}
