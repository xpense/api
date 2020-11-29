package repository

import (
	"expense-api/internal/model"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		model.User{},
		model.Wallet{},
		model.Party{},
		model.Transaction{},
	)
}
