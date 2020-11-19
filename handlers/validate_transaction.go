package handlers

import (
	"errors"
	"expense-api/model"
)

var (
	ErrorAmount = errors.New("cannot create new transaction with an amount of 0")
)

func TransactionCreateRequestToModel(t *model.Transaction) (*model.Transaction, error) {
	if t.Amount == 0 {
		return nil, ErrorAmount
	}

	return &model.Transaction{
		Amount:      t.Amount,
		Timestamp:   t.Timestamp,
		Description: t.Description,
	}, nil
}
