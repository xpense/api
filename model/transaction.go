package model

import (
	"errors"
	"time"
)

var (
	ErrorAmount = errors.New("cannot create new transaction with an amount of 0")
	ErrorType   = errors.New("invalid transaction Type; please use either 'income' or 'expense'")
)

type TransactionType string

const (
	Income  TransactionType = "income"
	Expense TransactionType = "expense"
)

type Transaction struct {
	Model
	Timestamp time.Time       `json:"timestamp,omitempty"`
	Amount    uint64          `json:"amount"`
	Type      TransactionType `json:"type"`
}

func NewTransaction(timestamp time.Time, amount uint64, transactionType TransactionType) (*Transaction, error) {
	if amount == 0 {
		return nil, ErrorAmount
	}

	if transactionType != Income && transactionType != Expense {
		return nil, ErrorType
	}

	transaction := &Transaction{
		Amount:    amount,
		Type:      transactionType,
		Timestamp: timestamp,
	}

	if timestamp.IsZero() {
		transaction.Timestamp = time.Now()
	}

	return transaction, nil
}
