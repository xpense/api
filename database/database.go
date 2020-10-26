package database

import "expense-api/transaction"

type DB struct {
	TotalTransactions uint64                                `json:"total_transactions"`
	UserTransactions  map[string][]*transaction.Transaction `json:"user_transactions"`
}

func NewDB() DB {
	return DB{
		TotalTransactions: 0,
		UserTransactions:  map[string][]*transaction.Transaction{},
	}
}

func (db *DB) AddTransaction(user string, newTransaction *transaction.Transaction) {
	userTransactions, ok := db.UserTransactions[user]

	if !ok {
		db.UserTransactions[user] = []*transaction.Transaction{}
	}

	db.UserTransactions[user] = append(userTransactions, newTransaction)
	db.TotalTransactions++
}
