package model

import (
	"errors"
	"regexp"
	"unicode"
)

var (
	ErrorName                = errors.New("first and/or last name missing")
	ErrorEmail               = errors.New("invalid email address")
	ErrorPasswordLength      = errors.New("password must be at least 8 characters")
	ErrorPasswordLowerChar   = errors.New("password must contain at least one lowercase char")
	ErrorPasswordUpperChar   = errors.New("password must contain at least one uppercase char")
	ErrorPasswordDigitChar   = errors.New("password must contain at least one digit")
	ErrorPasswordSpecialChar = errors.New("password must contain at least one special char")
)

func UserValidateInfo(firstName, lastName, email string) error {
	if firstName == "" || lastName == "" {
		return ErrorName
	}

	if !isEmailValid(email) {
		return ErrorEmail
	}

	return nil
}

func UserValidateCreateBody(firstName, lastName, email, password string) error {
	if err := UserValidateInfo(firstName, lastName, email); err != nil {
		return err
	}

	return isPasswordStrong(password)
}

var emailRegex = regexp.MustCompile(
	"^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$",
)

// isEmailValid checks if the email provided passes the required structure and length.
func isEmailValid(e string) bool {
	if len(e) < 3 || len(e) > 254 {
		return false
	}
	return emailRegex.MatchString(e)
}

func isPasswordStrong(p string) error {
	if len(p) < 8 {
		return ErrorPasswordLength
	}

	containsLower := false
	containsUpper := false
	containsDigit := false
	containsSpecial := false

	for _, r := range p {
		if unicode.IsLower(r) {
			containsLower = true
		} else if unicode.IsUpper(r) {
			containsUpper = true
		} else if unicode.IsDigit(r) {
			containsDigit = true
		} else if unicode.IsSymbol(r) {
			containsSpecial = true
		}
	}

	if !containsLower {
		return ErrorPasswordLowerChar
	} else if !containsUpper {
		return ErrorPasswordUpperChar
	} else if !containsDigit {
		return ErrorPasswordDigitChar
	} else if !containsSpecial {
		return ErrorPasswordSpecialChar
	}

	return nil
}
