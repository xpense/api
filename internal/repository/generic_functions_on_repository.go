package repository

import (
	"expense-api/internal/model"

	"gorm.io/gorm"
)

func genericCreate[
	M *model.User |
		*model.Wallet |
		*model.Transaction |
		*model.Party,
](r *repository, model M) (M, error) {
	if tx := r.db.Create(model); tx.Error != nil {
		if isUniqueConstaintViolationError(tx.Error) {
			return nil, ErrorUniqueConstaintViolation
		}
		return nil, ErrorOther
	}
	return model, nil
}

func genericGet[
	M *model.User |
		*model.Wallet |
		*model.Transaction |
		*model.Party,
](r *repository, model M, id uint) error {
	if tx := r.db.First(&model, id); tx.Error != nil {
		if tx.Error == gorm.ErrRecordNotFound {
			return ErrorRecordNotFound
		}
		return ErrorOther
	}
	return nil
}

func genericDelete[
	M *model.User |
		*model.Wallet |
		*model.Transaction |
		*model.Party,
](r *repository, model M, id uint) error {
	if err := genericGet(r, model, id); err != nil {
		return err
	}
	if tx := r.db.Delete(model); tx.Error != nil {
		return ErrorOther
	}
	return nil
}
