package utils

import (
	"errors"
	"strings"
	"unicode"
)

// ErrorPasswordLength and the rest are password errors that are returned from the api
var (
	ErrorPasswordLength      = errors.New("password must be at least 8 characters")
	ErrorPasswordLowerChar   = errors.New("password must contain at least one lowercase char")
	ErrorPasswordUpperChar   = errors.New("password must contain at least one uppercase char")
	ErrorPasswordDigitChar   = errors.New("password must contain at least one digit")
	ErrorPasswordSpecialChar = errors.New("password must contain at least one special char")
)

var specialChars = " !\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~"

// IsPasswordStrong checks if provided password is strong enough
func IsPasswordStrong(password string) (bool, error) {
	if len(password) < 8 {
		return false, ErrorPasswordLength
	} else if !strings.ContainsAny(password, specialChars) {
		return false, ErrorPasswordSpecialChar
	}

	var hasLower bool
	var hasUpper bool
	var hasDigit bool

	for _, r := range password {
		if unicode.IsLower(r) {
			hasLower = true
		} else if unicode.IsUpper(r) {
			hasUpper = true
		} else if unicode.IsDigit(r) {
			hasDigit = true
		}
	}

	if !hasLower {
		return false, ErrorPasswordLowerChar
	} else if !hasUpper {
		return false, ErrorPasswordUpperChar
	} else if !hasDigit {
		return false, ErrorPasswordDigitChar
	}

	return true, nil
}
