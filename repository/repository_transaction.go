package repository

import (
	"expense-api/model"
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
		return nil, ErrorOther
	}

	return transaction, nil
}

func (r *repository) TransactionUpdate(id uint, timestamp time.Time, amount uint64, tType model.TransactionType) (*model.Transaction, error) {
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

	if tType == model.Income || tType == model.Expense {
		transaction.Type = tType
	}

	if tx := r.db.Save(transaction); tx.Error != nil {
		return nil, ErrorOther
	}

	return transaction, nil
}

func (r *repository) TransactionGet(id uint) (*model.Transaction, error) {
	var transaction model.Transaction

	if tx := r.db.First(&transaction); tx.Error != nil {
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

func (r *repository) TransactionList() ([]*model.Transaction, error) {
	var transactions []*model.Transaction

	if tx := r.db.Find(&transactions); tx.Error != nil {
		return nil, ErrorOther
	}

	return transactions, nil
}
