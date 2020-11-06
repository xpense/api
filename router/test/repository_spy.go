package test

import (
	"expense-api/model"
	"expense-api/repository"
	"sort"
	"time"
)

type RepositorySpy struct {
	transactions map[uint]*model.Transaction
}

func NewRepositorySpy() *RepositorySpy {
	return &RepositorySpy{
		transactions: map[uint]*model.Transaction{},
	}
}

func (r *RepositorySpy) transactionSlice() []*model.Transaction {
	res := make([]*model.Transaction, 0, len(r.transactions))

	for _, t := range r.transactions {
		res = append(res, t)
	}

	sort.Slice(res, func(i, j int) bool {
		return res[i].ID < res[j].ID
	})

	return res
}

func (r *RepositorySpy) TransactionCreate(timestamp time.Time, amount uint64, transactionType model.TransactionType) (*model.Transaction, error) {
	id := uint(len(r.transactions)) + 1

	t := &model.Transaction{}

	t.ID = id
	if !timestamp.IsZero() {
		t.Timestamp = timestamp.Round(0)
	}

	if amount > 0 {
		t.Amount = amount
	}

	if transactionType == model.Income || t.Type == model.Expense {
		t.Type = transactionType
	}

	r.transactions[id] = t

	return t, nil
}

func (r *RepositorySpy) TransactionUpdate(id uint, timestamp time.Time, amount uint64, transactionType model.TransactionType) (*model.Transaction, error) {
	t, err := r.TransactionGet(id)
	if err != nil {
		return nil, err
	}

	if !timestamp.IsZero() {
		t.Timestamp = timestamp
	}

	if amount > 0 {
		t.Amount = amount
	}

	if transactionType == model.Income || t.Type == model.Expense {
		t.Type = transactionType
	}

	return t, nil
}

func (r *RepositorySpy) TransactionGet(id uint) (*model.Transaction, error) {
	t, ok := r.transactions[id]
	if !ok {
		return nil, repository.ErrorRecordNotFound
	}

	return t, nil
}

func (r *RepositorySpy) TransactionDelete(id uint) error {
	if _, err := r.TransactionGet(id); err != nil {
		return err
	}

	delete(r.transactions, id)
	return nil
}

func (r *RepositorySpy) TransactionList() ([]*model.Transaction, error) {
	transactions := make([]*model.Transaction, 0, len(r.transactions))

	for _, t := range r.transactions {
		transactions = append(transactions, t)
	}

	return transactions, nil
}

// ---------------- USER -----------------

func (r *RepositorySpy) UserCreate(firstName, LastName, Email, Password, Salt string) (*model.User, error) {
	return nil, nil
}

func (r *RepositorySpy) UserUpdate(id uint, firstName, LastName, Email string) (*model.User, error) {
	return nil, nil
}

func (r *RepositorySpy) UserDelete(id uint) error {
	return nil
}

func (r *RepositorySpy) UserGet(id uint) (*model.User, error) {
	return nil, nil
}
