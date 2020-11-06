package repository

import (
	"errors"
	"expense-api/model"
	"time"

	"gorm.io/gorm"
)

type Repository interface {
	TransactionCreate(timestamp time.Time, amount uint64, transactionType model.TransactionType) (*model.Transaction, error)
	TransactionUpdate(id uint, timestamp time.Time, amount uint64, transactionType model.TransactionType) (*model.Transaction, error)
	TransactionGet(id uint) (*model.Transaction, error)
	TransactionDelete(id uint) error
	TransactionList() ([]*model.Transaction, error)
	UserCreate(firstName, LastName, Email, Password, Salt string) (*model.User, error)
	UserUpdate(id uint, firstName, LastName, Email string) (*model.User, error)
	UserDelete(id uint) error
	UserGet(id uint) (*model.User, error)
}

var (
	ErrorRecordNotFound = errors.New("resource not found")
	ErrorOther          = errors.New("an error occurred")
)

type repository struct {
	db *gorm.DB
}

func New(db *gorm.DB) Repository {
	return &repository{db}
}
