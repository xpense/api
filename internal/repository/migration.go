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

func Cleanup(db *gorm.DB) error {
	migrator := db.Migrator()

	if migrator.HasTable(model.User{}) {
		if err := migrator.DropTable(model.User{}); err != nil {
			return err
		}
	}

	if migrator.HasTable(model.Wallet{}) {
		if err := migrator.DropTable(model.Wallet{}); err != nil {
			return err
		}
	}

	if migrator.HasTable(model.Party{}) {
		if err := migrator.DropTable(model.Party{}); err != nil {
			return err
		}
	}

	if migrator.HasTable(model.Transaction{}) {
		if err := migrator.DropTable(model.Transaction{}); err != nil {
			return err
		}
	}

	return nil
}
