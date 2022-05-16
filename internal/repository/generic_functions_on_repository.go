package repository

import (
	"expense-api/internal/model"
)

func genericCreate[M model.GormModel](r *repository, model *M) error {
	if tx := r.db.Create(model); tx.Error != nil {
		return checkError(tx.Error)
	}
	return nil
}

func genericSave[M model.GormModel](r *repository, model *M) error {
	if tx := r.db.Save(model); tx.Error != nil {
		return checkError(tx.Error)
	}
	return nil
}

func genericGet[M model.GormModel](r *repository, query map[string]interface{}) (*M, error) {
	var model M
	if tx := r.db.Where(query).First(&model); tx.Error != nil {
		return nil, checkError(tx.Error)
	}
	return &model, nil
}

func genericDelete[M model.GormModel](r *repository, id uint) error {
	model, err := genericGet[M](r, map[string]interface{}{"id": id})
	if err != nil {
		return err
	}
	if tx := r.db.Delete(model); tx.Error != nil {
		return checkError(tx.Error)
	}
	return nil
}

func genericList[M model.GormModel](r *repository, query map[string]interface{}) ([]*M, error) {
	var models []*M
	if tx := r.db.Where(query).Find(&models); tx.Error != nil {
		return nil, checkError(tx.Error)
	}
	return models, nil
}
