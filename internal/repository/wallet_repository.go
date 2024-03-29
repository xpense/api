package repository

import (
	"expense-api/internal/model"
)

func (r *repository) WalletCreate(w *model.Wallet) error {
	return genericCreate(r, w)
}

func (r *repository) WalletUpdate(id uint, updated *model.Wallet) (*model.Wallet, error) {
	wallet, err := r.WalletGet(id)
	if err != nil {
		return nil, err
	}

	if updated.Name != "" {
		wallet.Name = updated.Name
	}

	if updated.Description != "" {
		wallet.Description = updated.Description
	}

	err = genericSave(r, wallet)
	return wallet, err
}

func (r *repository) WalletGet(id uint) (*model.Wallet, error) {
	return genericGet[model.Wallet](r, map[string]interface{}{"id": id})
}

func (r *repository) WalletDelete(id uint) error {
	return genericDelete[model.Wallet](r, id)
}

func (r *repository) WalletList(userID uint) ([]*model.Wallet, error) {
	return genericList[model.Wallet](r, map[string]interface{}{"user_id": userID})
}
