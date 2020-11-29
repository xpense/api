package handlers

import (
	"expense-api/internal/model"
	"time"

	"github.com/shopspring/decimal"
)

// Transaction is a transaction with an omitted user
type Transaction struct {
	ID          uint            `json:"id"`
	WalletID    uint            `json:"wallet_id"`
	PartyID     uint            `json:"party_id"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	Timestamp   time.Time       `json:"timestamp"`
	Amount      decimal.Decimal `json:"amount"`
	Description string          `json:"description"`
}

func TransactionModelToResponse(t *model.Transaction) *Transaction {
	return &Transaction{
		ID:          t.ID,
		WalletID:    t.WalletID,
		PartyID:     t.PartyID,
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
		Timestamp:   t.Timestamp,
		Amount:      t.Amount,
		Description: t.Description,
	}
}

func TransactionRequestToModel(t *Transaction, userID uint) *model.Transaction {
	return &model.Transaction{
		Amount:      t.Amount,
		Timestamp:   t.Timestamp,
		Description: t.Description,
		WalletID:    t.WalletID,
		PartyID:     t.PartyID,
		UserID:      userID,
	}
}
