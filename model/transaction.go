package model

import (
	"errors"
	"time"
)

var (
	ErrorAmount = errors.New("cannot create new transaction with an amount of 0")
	ErrorType   = errors.New("invalid transaction Type; please use either 'income' or 'expense'")
)

func TransactionValidateCreateBody(timestamp time.Time, amount uint64, transactionType TransactionType) error {
	if amount == 0 {
		return ErrorAmount
	}

	if transactionType != Income && transactionType != Expense {
		return ErrorType
	}

	return nil
}
