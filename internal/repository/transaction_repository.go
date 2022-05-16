package repository

import (
	"expense-api/internal/model"
	"time"

	"github.com/shopspring/decimal"
)

func (r *repository) TransactionCreate(t *model.Transaction) error {
	if t.Timestamp.IsZero() {
		t.Timestamp = time.Now()
	}
	return genericCreate(r, t)
}

func (r *repository) TransactionUpdate(id uint, updated *model.Transaction) (*model.Transaction, error) {
	transaction, err := r.TransactionGet(id)
	if err != nil {
		return nil, err
	}

	if !updated.Timestamp.IsZero() {
		transaction.Timestamp = updated.Timestamp
	}

	if updated.Amount.Cmp(decimal.Zero) != 0 {
		transaction.Amount = updated.Amount
	}

	if updated.Description != "" {
		transaction.Description = updated.Description
	}

	if updated.WalletID != 0 {
		transaction.WalletID = updated.WalletID
	}

	if tx := r.db.Save(transaction); tx.Error != nil {
		return nil, ErrorOther
	}

	return transaction, nil
}

func (r *repository) TransactionGet(id uint) (*model.Transaction, error) {
	return genericGet[*model.Transaction](r, int(id), nil)
}

func (r *repository) TransactionDelete(id uint) error {
	return genericDelete[*model.Transaction](r, id)
}

func (r *repository) TransactionList(userID uint) ([]*model.Transaction, error) {
	return r.transactionList(map[string]interface{}{"user_id": userID})
}

func (r *repository) TransactionListByWallet(userID, walletID uint) ([]*model.Transaction, error) {
	return r.transactionList(map[string]interface{}{
		"user_id":   userID,
		"wallet_id": walletID,
	})
}

func (r *repository) TransactionListByParty(userID, partyID uint) ([]*model.Transaction, error) {
	return r.transactionList(map[string]interface{}{
		"user_id":  userID,
		"party_id": partyID,
	})
}

func (r *repository) transactionList(query map[string]interface{}) ([]*model.Transaction, error) {
	var transactions []*model.Transaction
	err := genericList(r, &transactions, query)
	return transactions, err
}
