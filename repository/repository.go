package repository

import "expense-api/model"

type DB struct {
	TotalTransactions uint64                          `json:"total_transactions"`
	UserTransactions  map[string][]*model.Transaction `json:"user_transactions"`
}

func NewDB() DB {
	return DB{
		TotalTransactions: 0,
		UserTransactions:  map[string][]*model.Transaction{},
	}
}

func (db *DB) AddTransaction(user string, newTransaction *model.Transaction) {
	userTransactions, ok := db.UserTransactions[user]

	if !ok {
		db.UserTransactions[user] = []*model.Transaction{}
	}

	db.UserTransactions[user] = append(userTransactions, newTransaction)
	db.TotalTransactions++
}
