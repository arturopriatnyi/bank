// Code generated by protoc-gen-connect-go. DO NOT EDIT.
//
// Source: proto/bank/v1/bank.proto

package bankv1connect

import (
	v1 "bank/gen/proto/bank/v1"
	connect "connectrpc.com/connect"
	context "context"
	errors "errors"
	http "net/http"
	strings "strings"
)

// This is a compile-time assertion to ensure that this generated file and the connect package are
// compatible. If you get a compiler error that this constant is not defined, this code was
// generated with a version of connect newer than the one compiled into your binary. You can fix the
// problem by either regenerating this code with an older version of connect or updating the connect
// version compiled into your binary.
const _ = connect.IsAtLeastVersion0_1_0

const (
	// BankServiceName is the fully-qualified name of the BankService service.
	BankServiceName = "bank.v1.BankService"
)

// These constants are the fully-qualified names of the RPCs defined in this package. They're
// exposed at runtime as Spec.Procedure and as the final two segments of the HTTP route.
//
// Note that these are different from the fully-qualified method names used by
// google.golang.org/protobuf/reflect/protoreflect. To convert from these constants to
// reflection-formatted method names, remove the leading slash and convert the remaining slash to a
// period.
const (
	// BankServiceMakeTransactionProcedure is the fully-qualified name of the BankService's
	// MakeTransaction RPC.
	BankServiceMakeTransactionProcedure = "/bank.v1.BankService/MakeTransaction"
	// BankServiceBalanceProcedure is the fully-qualified name of the BankService's Balance RPC.
	BankServiceBalanceProcedure = "/bank.v1.BankService/Balance"
)

// BankServiceClient is a client for the bank.v1.BankService service.
type BankServiceClient interface {
	MakeTransaction(context.Context, *connect.Request[v1.MakeTransactionRequest]) (*connect.Response[v1.MakeTransactionResponse], error)
	Balance(context.Context, *connect.Request[v1.BalanceRequest]) (*connect.Response[v1.BalanceResponse], error)
}

// NewBankServiceClient constructs a client for the bank.v1.BankService service. By default, it uses
// the Connect protocol with the binary Protobuf Codec, asks for gzipped responses, and sends
// uncompressed requests. To use the gRPC or gRPC-Web protocols, supply the connect.WithGRPC() or
// connect.WithGRPCWeb() options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewBankServiceClient(httpClient connect.HTTPClient, baseURL string, opts ...connect.ClientOption) BankServiceClient {
	baseURL = strings.TrimRight(baseURL, "/")
	return &bankServiceClient{
		makeTransaction: connect.NewClient[v1.MakeTransactionRequest, v1.MakeTransactionResponse](
			httpClient,
			baseURL+BankServiceMakeTransactionProcedure,
			opts...,
		),
		balance: connect.NewClient[v1.BalanceRequest, v1.BalanceResponse](
			httpClient,
			baseURL+BankServiceBalanceProcedure,
			opts...,
		),
	}
}

// bankServiceClient implements BankServiceClient.
type bankServiceClient struct {
	makeTransaction *connect.Client[v1.MakeTransactionRequest, v1.MakeTransactionResponse]
	balance         *connect.Client[v1.BalanceRequest, v1.BalanceResponse]
}

// MakeTransaction calls bank.v1.BankService.MakeTransaction.
func (c *bankServiceClient) MakeTransaction(ctx context.Context, req *connect.Request[v1.MakeTransactionRequest]) (*connect.Response[v1.MakeTransactionResponse], error) {
	return c.makeTransaction.CallUnary(ctx, req)
}

// Balance calls bank.v1.BankService.Balance.
func (c *bankServiceClient) Balance(ctx context.Context, req *connect.Request[v1.BalanceRequest]) (*connect.Response[v1.BalanceResponse], error) {
	return c.balance.CallUnary(ctx, req)
}

// BankServiceHandler is an implementation of the bank.v1.BankService service.
type BankServiceHandler interface {
	MakeTransaction(context.Context, *connect.Request[v1.MakeTransactionRequest]) (*connect.Response[v1.MakeTransactionResponse], error)
	Balance(context.Context, *connect.Request[v1.BalanceRequest]) (*connect.Response[v1.BalanceResponse], error)
}

// NewBankServiceHandler builds an HTTP handler from the service implementation. It returns the path
// on which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewBankServiceHandler(svc BankServiceHandler, opts ...connect.HandlerOption) (string, http.Handler) {
	bankServiceMakeTransactionHandler := connect.NewUnaryHandler(
		BankServiceMakeTransactionProcedure,
		svc.MakeTransaction,
		opts...,
	)
	bankServiceBalanceHandler := connect.NewUnaryHandler(
		BankServiceBalanceProcedure,
		svc.Balance,
		opts...,
	)
	return "/bank.v1.BankService/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case BankServiceMakeTransactionProcedure:
			bankServiceMakeTransactionHandler.ServeHTTP(w, r)
		case BankServiceBalanceProcedure:
			bankServiceBalanceHandler.ServeHTTP(w, r)
		default:
			http.NotFound(w, r)
		}
	})
}

// UnimplementedBankServiceHandler returns CodeUnimplemented from all methods.
type UnimplementedBankServiceHandler struct{}

func (UnimplementedBankServiceHandler) MakeTransaction(context.Context, *connect.Request[v1.MakeTransactionRequest]) (*connect.Response[v1.MakeTransactionResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("bank.v1.BankService.MakeTransaction is not implemented"))
}

func (UnimplementedBankServiceHandler) Balance(context.Context, *connect.Request[v1.BalanceRequest]) (*connect.Response[v1.BalanceResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("bank.v1.BankService.Balance is not implemented"))
}