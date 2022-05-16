package repository

import (
	"expense-api/internal/model"
)

func (r *repository) PartyCreate(p *model.Party) error {
	return genericCreate(r, p)
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
	return genericGet[*model.Party](r, int(id), nil)
}

func (r *repository) PartyDelete(id uint) error {
	return genericDelete[*model.Party](r, id)
}

func (r *repository) PartyList(userID uint) ([]*model.Party, error) {
	var partys []*model.Party
	err := genericList(r, &partys, map[string]interface{}{"user_id": userID})
	return partys, err
}
