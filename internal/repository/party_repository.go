package repository

import (
	"expense-api/internal/model"

	"gorm.io/gorm"
)

func (r *repository) PartyCreate(p *model.Party) error {
	_, err := genericCreate(r, p)
	return err
}

func (r *repository) PartyUpdate(id uint, updated *model.Party) (*model.Party, error) {
	party, err := r.PartyGet(id)
	if err != nil {
		return nil, err
	}

	if updated.Name != "" {
		party.Name = updated.Name
	}

	if tx := r.db.Save(party); tx.Error != nil {
		if isUniqueConstaintViolationError(tx.Error) {
			return nil, ErrorUniqueConstaintViolation
		}
		return nil, ErrorOther
	}

	return party, nil
}

func (r *repository) PartyGet(id uint) (*model.Party, error) {
	var party model.Party
	err := genericGet(r, &party, id)
	return &party, err
}

func (r *repository) PartyDelete(id uint) error {
	party, err := r.PartyGet(id)
	if err != nil {
		return err
	}

	if tx := r.db.Delete(party); tx.Error != nil {
		return ErrorOther
	}

	return nil
}

func (r *repository) PartyList(userID uint) ([]*model.Party, error) {
	var partys []*model.Party

	if tx := r.db.Where("user_id = ?", userID).Find(&partys); tx.Error != nil {
		if tx.Error == gorm.ErrRecordNotFound {
			return nil, ErrorRecordNotFound
		}
		return nil, ErrorOther
	}

	return partys, nil
}
