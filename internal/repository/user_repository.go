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

	if tx := r.db.Save(user); tx.Error != nil {
		return nil, ErrorOther
	}

	return user, nil
}

func (r *repository) UserDelete(id uint) error {
	var user model.User
	return genericDelete(r, &user, id)
}

func (r *repository) UserGet(id uint) (*model.User, error) {
	var user model.User
	err := genericGet(r, &user, int(id), nil)
	return &user, err
}

func (r *repository) UserGetWithEmail(email string) (*model.User, error) {
	var user model.User
	err := genericGet(r, &user, -1, map[string]interface{}{"email": email})
	return &user, err
}
