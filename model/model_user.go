package model

import (
	"errors"
	"expense-api/utils"
)

var (
	ErrorName  = errors.New("first and/or last name missing")
	ErrorEmail = errors.New("invalid email address")
)

func UserValidateInfo(firstName, lastName, email string) error {
	if firstName == "" || lastName == "" {
		return ErrorName
	}

	if !utils.IsEmailValid(email) {
		return ErrorEmail
	}

	return nil
}

func UserValidateCreateBody(firstName, lastName, email, password string) error {
	if err := UserValidateInfo(firstName, lastName, email); err != nil {
		return err
	}

	_, err := utils.IsPasswordStrong(password)

	return err
}
