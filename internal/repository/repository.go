package repository

import (
	"expense-api/internal/model"

	"gorm.io/gorm"
)

type Repository interface {
	UserCreate(firstName, LastName, Email, Password, Salt string) (*model.User, error)
	UserUpdate(id uint, firstName, LastName, Email string) (*model.User, error)
	UserDelete(id uint) error
	UserGet(id uint) (*model.User, error)
	UserGetWithEmail(email string) (*model.User, error)

	WalletCreate(w *model.Wallet) error
	WalletUpdate(id uint, w *model.Wallet) (*model.Wallet, error)
	WalletGet(id uint) (*model.Wallet, error)
	WalletDelete(id uint) error
	WalletList(userID uint) ([]*model.Wallet, error)

	PartyCreate(w *model.Party) error
	PartyUpdate(id uint, w *model.Party) (*model.Party, error)
	PartyGet(id uint) (*model.Party, error)
	PartyDelete(id uint) error
	PartyList(userID uint) ([]*model.Party, error)

	TransactionCreate(t *model.Transaction) error
	TransactionUpdate(id uint, t *model.Transaction) (*model.Transaction, error)
	TransactionGet(id uint) (*model.Transaction, error)
	TransactionDelete(id uint) error
	TransactionList(userID uint) ([]*model.Transaction, error)
	TransactionListByWallet(userID, walletID uint) ([]*model.Transaction, error)
	TransactionListByParty(userID, partyID uint) ([]*model.Transaction, error)
}

type repository struct {
	db *gorm.DB
}

func New(db *gorm.DB) Repository {
	return &repository{db}
}
