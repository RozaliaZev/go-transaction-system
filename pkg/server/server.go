package server

import (
	"context"
	"go-tr-syst/internal/db"
	ts "go-tr-syst/pkg/gen-tr-syst"
	"strconv"
)

type GRPCServer struct {
	ts.UnimplementedTransactionServiceServer
}

func (s GRPCServer) CreateTransaction(ctx context.Context, req *ts.CreateTransactionRequest) (*ts.CreateTransactionResponse, error) {
	err := db.CreateTable()
	if err != nil {
		return nil, err
	}

	id, err := db.WriteTransactionData(req)
	if err != nil {
		return nil, err
	}

	return &ts.CreateTransactionResponse{Id: strconv.FormatInt(id, 10)}, err
}

func (s GRPCServer) GetTransaction(ctx context.Context, req *ts.GetTransactionRequest) (*ts.GetTransactionResponse, error) {
	var tr db.Transaction
	tr, err := db.GetTransactionByID(req)
	if err != nil {
		return nil, err
	}

	return &ts.GetTransactionResponse{
		Amount:               tr.Amount,
		Currency:             tr.Currency,
		SenderCardNumber:     tr.SenderCardNumber,
		RecipientPhoneNumber: tr.RecipientPhoneNumber,
		DateTime:             tr.DateTime,
	}, nil
}
