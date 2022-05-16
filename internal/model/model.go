package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type GormModel interface {
	*User | *Wallet | *Transaction | *Party
}

type Model struct {
	ID        uint      `json:"id" gorm:"primarykey"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type User struct {
	Model
	FirstName string `json:"first_name" gorm:"not null;"`
	LastName  string `json:"last_name" gorm:"not null;"`
	Email     string `json:"email" gorm:"type:varchar(255);unique;not null;"`
	Password  string `json:"password" gorm:"not null;"`
	Salt      string `json:"salt" gorm:"not null;"`
}

type Wallet struct {
	Model
	Name        string `json:"name" gorm:"uniqueIndex:idx_userid_wallet_name;not null;"`
	Description string `json:"description"`
	UserID      uint   `json:"user_id" gorm:"uniqueIndex:idx_userid_wallet_name;not null;"`
	User        User   `json:"user" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type Party struct {
	Model
	Name   string `json:"name" gorm:"uniqueIndex:idx_userid_party_name;not null;"`
	UserID uint   `json:"user_id" gorm:"uniqueIndex:idx_userid_party_name;not null;"`
	User   User   `json:"user" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type Transaction struct {
	Model
	Description string          `json:"description"`
	Timestamp   time.Time       `json:"timestamp"`
	Amount      decimal.Decimal `json:"amount" gorm:"type:numeric"`
	UserID      uint            `json:"user_id" gorm:"not null;"`
	User        User            `json:"user" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	WalletID    uint            `json:"wallet_id" gorm:"not null;"`
	Wallet      Wallet          `json:"wallet" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	PartyID     uint            `json:"party_id" gorm:"not null;"`
	Party       Party           `json:"party" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
