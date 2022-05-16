package repository

import (
	"expense-api/internal/model"
)

func (r *repository) UserCreate(firstName, lastName, email, password, salt string) (*model.User, error) {
	user := model.User{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Password:  password,
		Salt:      salt,
	}
	err := genericCreate(r, &user)
	return &user, err
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

	err = genericSave(r, user)
	return user, err
}

func (r *repository) UserDelete(id uint) error {
	return genericDelete[model.User](r, id)
}

func (r *repository) UserGet(id uint) (*model.User, error) {
	return genericGet[model.User](r, map[string]interface{}{"id": id})
}

func (r *repository) UserGetWithEmail(email string) (*model.User, error) {
	return genericGet[model.User](r, map[string]interface{}{"email": email})
}
