package repository

import (
	"expense-api/model"

	"gorm.io/gorm"
)

func (r *repository) WalletCreate(w *model.Wallet) error {
	if tx := r.db.Create(w); tx.Error != nil {
		if isUniqueConstaintViolationError(tx.Error) {
			return ErrorUniqueConstaintViolation
		}
		return ErrorOther
	}

	return nil
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

	if tx := r.db.Save(wallet); tx.Error != nil {
		if isUniqueConstaintViolationError(tx.Error) {
			return nil, ErrorUniqueConstaintViolation
		}
		return nil, ErrorOther
	}

	return wallet, nil
}

func (r *repository) WalletGet(id uint) (*model.Wallet, error) {
	var wallet model.Wallet

	if tx := r.db.First(&wallet, id); tx.Error != nil {
		if tx.Error == gorm.ErrRecordNotFound {
			return nil, ErrorRecordNotFound
		}
		return nil, ErrorOther
	}

	return &wallet, nil
}

func (r *repository) WalletDelete(id uint) error {
	wallet, err := r.WalletGet(id)
	if err != nil {
		return err
	}

	if tx := r.db.Delete(wallet); tx.Error != nil {
		return ErrorOther
	}

	return nil
}

func (r *repository) WalletList(userID uint) ([]*model.Wallet, error) {
	var wallets []*model.Wallet

	if tx := r.db.Where("user_id = ?", userID).Find(&wallets); tx.Error != nil {
		if tx.Error == gorm.ErrRecordNotFound {
			return nil, ErrorRecordNotFound
		}
		return nil, ErrorOther
	}

	return wallets, nil
}
