package model

import (
	"time"
)

type Model struct {
	ID        uint      `json:"id" gorm:"primarykey"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type TransactionType string

const (
	Income  TransactionType = "income"
	Expense TransactionType = "expense"
)

type Transaction struct {
	Model
	Description string
	Timestamp   time.Time
	Amount      uint64
	Type        TransactionType
	UserID      uint
	User        User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type User struct {
	Model
	FirstName string `json:"firstName" gorm:"not null;"`
	LastName  string `json:"lastName" gorm:"not null;"`
	Email     string `json:"email" gorm:"type:varchar(255);unique;not null;"`
	Password  string `json:"password" gorm:"not null;"`
	Salt      string `json:"salt" gorm:"not null;"`
}
