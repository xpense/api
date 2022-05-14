package repository

import (
	"expense-api/internal/model"
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

func (r *repository) TransactionCreate(t *model.Transaction) error {
	if t.Timestamp.IsZero() {
		t.Timestamp = time.Now()
	}
	_, err := genericCreate(r, t)
	return err
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
	var transaction model.Transaction

	if tx := r.db.First(&transaction, id); tx.Error != nil {
		if tx.Error == gorm.ErrRecordNotFound {
			return nil, ErrorRecordNotFound
		}
		return nil, ErrorOther
	}

	return &transaction, nil
}

func (r *repository) TransactionDelete(id uint) error {
	transaction, err := r.TransactionGet(id)
	if err != nil {
		return err
	}

	if tx := r.db.Delete(transaction); tx.Error != nil {
		return ErrorOther
	}

	return nil
}

func (r *repository) TransactionList(userID uint) ([]*model.Transaction, error) {
	query := map[string]interface{}{
		"user_id": userID,
	}

	return r.transactionList(query)
}

func (r *repository) TransactionListByWallet(userID, walletID uint) ([]*model.Transaction, error) {
	query := map[string]interface{}{
		"user_id":   userID,
		"wallet_id": walletID,
	}

	return r.transactionList(query)
}

func (r *repository) TransactionListByParty(userID, partyID uint) ([]*model.Transaction, error) {
	query := map[string]interface{}{
		"user_id":  userID,
		"party_id": partyID,
	}

	return r.transactionList(query)
}

func (r *repository) transactionList(query map[string]interface{}) ([]*model.Transaction, error) {
	var transactions []*model.Transaction

	if tx := r.db.Where(query).Find(&transactions); tx.Error != nil {
		if tx.Error == gorm.ErrRecordNotFound {
			return nil, ErrorRecordNotFound
		}
		return nil, ErrorOther
	}

	return transactions, nil
}
