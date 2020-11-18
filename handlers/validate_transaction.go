package handlers

import (
	"errors"
	"expense-api/model"
)

var (
	ErrorAmount = errors.New("cannot create new transaction with an amount of 0")
	ErrorType   = errors.New("invalid transaction Type; please use either 'income' or 'expense'")
)

func TransactionCreateRequestToModel(t *model.Transaction) (*model.Transaction, error) {
	if t.Amount == 0 {
		return nil, ErrorAmount
	}

	if t.Type != model.Income && t.Type != model.Expense {
		return nil, ErrorType
	}

	return &model.Transaction{
		Amount:      t.Amount,
		Type:        t.Type,
		Timestamp:   t.Timestamp,
		Description: t.Description,
	}, nil
}

func TransactionValidateUpdateBody(t *model.Transaction) error {
	if t.Type != "" {
		if t.Type != model.Income && t.Type != model.Expense {
			return ErrorType
		}
	}

	return nil
}
