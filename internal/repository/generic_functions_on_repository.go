package repository

import (
	"expense-api/internal/model"

	"gorm.io/gorm"
)

func genericCreate[M model.GormModel](r *repository, model *M) error {
	if tx := r.db.Create(model); tx.Error != nil {
		if isUniqueConstaintViolationError(tx.Error) {
			return ErrorUniqueConstaintViolation
		}
		return ErrorOther
	}
	return nil
}

func genericSave[M model.GormModel](r *repository, model *M) error {
	if tx := r.db.Save(model); tx.Error != nil {
		if isUniqueConstaintViolationError(tx.Error) {
			return ErrorUniqueConstaintViolation
		}
		return ErrorOther
	}
	return nil
}

func genericGet[M model.GormModel](r *repository, id int, query map[string]interface{}) (*M, error) {
	var model M
	var tx *gorm.DB
	if id >= 0 {
		tx = r.db.Where(query).First(&model, id)
	} else {
		tx = r.db.Where(query).First(&model)
	}

	if tx.Error != nil {
		if tx.Error == gorm.ErrRecordNotFound {
			return nil, ErrorRecordNotFound
		}
		return nil, ErrorOther
	}
	return &model, nil
}

func genericDelete[M model.GormModel](r *repository, id uint) error {
	model, err := genericGet[M](r, int(id), nil)
	if err != nil {
		return err
	}
	if tx := r.db.Delete(model); tx.Error != nil {
		return ErrorOther
	}
	return nil
}

func genericList[M model.GormModel](r *repository, query map[string]interface{}) ([]*M, error) {
	var models []*M
	if tx := r.db.Where(query).Find(&models); tx.Error != nil {
		if tx.Error == gorm.ErrRecordNotFound {
			return nil, ErrorRecordNotFound
		}
		return nil, ErrorOther
	}
	return models, nil
}
