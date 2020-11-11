package repository

import (
	"expense-api/model"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(model.Transaction{}, model.User{})
}