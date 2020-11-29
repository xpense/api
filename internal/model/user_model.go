package model

import (
	"errors"
	"expense-api/internal/utils"
)

var (
	ErrorName      = errors.New("first and/or last name missing")
	ErrorEmail     = errors.New("invalid email address")
	ErrorEmptyBody = errors.New("empty body")
)

func UserValidateUpdateBody(firstName, lastName, email string) error {
	if firstName == "" && lastName == "" && email == "" {
		return ErrorEmptyBody
	}

	if email != "" && !utils.IsEmailValid(email) {
		return ErrorEmail
	}

	return nil
}
