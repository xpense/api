package handlers

import "errors"

var (
	ErrorName                   = errors.New("first and/or last name missing")
	ErrorEmail                  = errors.New("invalid email address")
	ErrorEmptyBody              = errors.New("empty body")
	ErrorMissingPasswordOrEmail = errors.New("both email and password are required for login")
	ErrorNonExistentUser        = errors.New("user with this email does not exist")
	ErrorEmailConflict          = errors.New("user with this email already exists")
	ErrorWrongPassword          = errors.New("wrong password")
)
