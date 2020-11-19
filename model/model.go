package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type Model struct {
	ID        uint      `json:"id" gorm:"primarykey"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Transaction struct {
	Model
	Description string
	Timestamp   time.Time
	Amount      decimal.Decimal `gorm:"type:numeric"`
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
