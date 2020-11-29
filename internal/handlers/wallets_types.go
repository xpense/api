package handlers

import (
	"expense-api/internal/model"
	"time"
)

// Wallet is a list of transactions belonging to an account
type Wallet struct {
	ID          uint      `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
}

func WalletModelToResponse(w *model.Wallet) *Wallet {
	return &Wallet{
		ID:          w.ID,
		CreatedAt:   w.CreatedAt,
		UpdatedAt:   w.UpdatedAt,
		Name:        w.Name,
		Description: w.Description,
	}
}

func WalletRequestToModel(w *Wallet, userID uint) *model.Wallet {
	return &model.Wallet{
		Name:        w.Name,
		Description: w.Description,
		UserID:      userID,
	}
}
