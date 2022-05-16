package repository

import (
	"expense-api/internal/model"

	"gorm.io/gorm"
)

func genericCreate[M model.GormModel](r *repository, model M) (M, error) {
	if tx := r.db.Create(model); tx.Error != nil {
		if isUniqueConstaintViolationError(tx.Error) {
			return nil, ErrorUniqueConstaintViolation
		}
		return nil, ErrorOther
	}
	return model, nil
}

func genericGet[M model.GormModel](r *repository, model M, id int, query map[string]interface{}) error {
	tx := r.db.Where(query)
	if id >= 0 {
		tx = tx.First(&model, id)
	} else {
		tx = tx.First(&model)
	}

	if tx.Error != nil {
		if tx.Error == gorm.ErrRecordNotFound {
			return ErrorRecordNotFound
		}
		return ErrorOther
	}
	return nil
}

func genericDelete[M model.GormModel](r *repository, model M, id uint) error {
	if err := genericGet(r, model, int(id), nil); err != nil {
		return err
	}
	if tx := r.db.Delete(model); tx.Error != nil {
		return ErrorOther
	}
	return nil
}

func genericList[M model.GormModel](r *repository, models *[]M, query map[string]interface{}) error {
	if tx := r.db.Where(query).Find(models); tx.Error != nil {
		if tx.Error == gorm.ErrRecordNotFound {
			return ErrorRecordNotFound
		}
		return ErrorOther
	}
	return nil
}
