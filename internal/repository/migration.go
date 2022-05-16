package repository

import (
	"expense-api/internal/model"

	"gorm.io/gorm"
)

var models = []interface{}{
	model.User{},
	model.Wallet{},
	model.Party{},
	model.Transaction{},
}

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(models...)
}

func Cleanup(db *gorm.DB) error {
	migrator := db.Migrator()

	for _, model := range models {
		if migrator.HasTable(model) {
			if err := migrator.DropTable(model); err != nil {
				return err
			}
		}
	}

	return nil
}
