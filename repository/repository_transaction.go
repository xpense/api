package repository

import (
	"expense-api/model"
	"fmt"
	"time"
)

func (r *repository) TransactionCreate(timestamp time.Time, amount uint64, transactionType model.TransactionType) (*model.Transaction, error) {
	transaction := &model.Transaction{
		Amount:    amount,
		Type:      transactionType,
		Timestamp: timestamp,
	}

	if timestamp.IsZero() {
		transaction.Timestamp = time.Now()
	}

	if tx := r.db.Create(transaction); tx.Error != nil {
		return nil, tx.Error
	}

	return transaction, nil
}

func (r *repository) TransactionUpdate(id uint, timestamp time.Time, amount uint64, transactionType model.TransactionType) (*model.Transaction, error) {
	transaction, err := r.TransactionGet(id)
	if err != nil {
		return nil, err
	}

	if !timestamp.IsZero() {
		transaction.Timestamp = timestamp
	}

	if amount > 0 {
		transaction.Amount = amount
	}

	if transactionType == model.Income || transaction.Type == model.Expense {
		transaction.Type = transactionType
	}

	if tx := r.db.Save(transaction); tx.Error != nil {
		return nil, tx.Error
	}

	return transaction, nil
}

func (r *repository) TransactionGet(id uint) (*model.Transaction, error) {
	var transaction model.Transaction

	if tx := r.db.First(&transaction); tx.Error != nil {
		fmt.Printf("in get: %v\n", tx.Error)
		return nil, tx.Error
	}

	return &transaction, nil
}

func (r *repository) TransactionDelete(id uint) error {
	transaction, err := r.TransactionGet(id)
	if err != nil {
		return err
	}

	if tx := r.db.Delete(transaction); tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (r *repository) TransactionList() ([]*model.Transaction, error) {
	var transactions []*model.Transaction

	if tx := r.db.Find(&transactions); tx.Error != nil {
		return nil, tx.Error
	}

	return transactions, nil
}
