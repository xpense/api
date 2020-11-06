package repository

import (
	"expense-api/model"

	"gorm.io/gorm"
)

func (r *repository) UserCreate(firstName, lastName, email, password, salt string) (*model.User, error) {
	user := &model.User{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Password:  password,
		Salt:      salt,
	}

	if tx := r.db.Create(user); tx.Error != nil {
		return nil, ErrorOther
	}

	return user, nil
}

func (r *repository) UserUpdate(id uint, firstName, lastName, email string) (*model.User, error) {
	user, err := r.UserGet(id)
	if err != nil {
		return nil, err
	}

	if user.FirstName != firstName {
		user.FirstName = firstName
	}

	if user.LastName != lastName {
		user.LastName = lastName
	}

	if user.Email != email {
		user.Email = email
	}

	if tx := r.db.Save(user); tx.Error != nil {
		return nil, ErrorOther
	}

	return user, nil
}

func (r *repository) UserDelete(id uint) error {
	user, err := r.UserGet(id)
	if err != nil {
		return err
	}

	if tx := r.db.Delete(user); tx.Error != nil {
		return ErrorOther
	}

	return nil
}

func (r *repository) UserGet(id uint) (*model.User, error) {
	var user model.User

	if tx := r.db.First(&user, id); tx.Error != nil {
		if tx.Error == gorm.ErrRecordNotFound {
			return nil, ErrorRecordNotFound
		}
		return nil, ErrorOther
	}

	return &user, nil
}
