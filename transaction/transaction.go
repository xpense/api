package transaction

import (
	"time"
)

var (
	ErrorMessageID     = "Cannot create new Transaction with a non positive id."
	ErrorMessageAmount = "Cannot create new Transaction with an amount of 0."
	ErrorMessageType   = "Invalid Transaction Type. Please use either 'Income' or 'Expense'."
)

type Error struct {
	ID              uint64
	Amount          uint64
	TransactionType Type
}

func NewError(id, amount uint64, transactionType Type) *Error {
	return &Error{id, amount, transactionType}
}

func (t *Error) Error() string {
	switch true {
	case t.ID == 0:
		return ErrorMessageID
	case t.Amount == 0:
		return ErrorMessageAmount
	default:
		return ErrorMessageType
	}
}

type Type uint8

const (
	Income Type = iota
	Expense
)

type Transaction struct {
	ID        uint64    `json:"id,omiempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	Date      time.Time `json:"date, omitempty"`
	Amount    uint64    `json:"amount"`
	Type      Type      `json:"type"`
}

func New(id uint64, date *time.Time, amount uint64, transactionType Type) (*Transaction, error) {
	switch true {
	case id == 0:
		return nil, NewError(id, amount, transactionType)
	case amount == 0:
		return nil, NewError(id, amount, transactionType)
	case transactionType != Income && transactionType != Expense:
		return nil, NewError(id, amount, transactionType)
	}

	now := time.Now()
	transaction := &Transaction{
		ID:        id,
		CreatedAt: now,
		UpdatedAt: now,
		Date:      now,
		Amount:    amount,
		Type:      transactionType,
	}

	if date != nil {
		transaction.Date = *date
	}

	return transaction, nil
}
