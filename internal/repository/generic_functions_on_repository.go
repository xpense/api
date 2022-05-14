package repository

import "expense-api/internal/model"

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
