package handlers

import "expense-api/internal/utils"

type (
	LoginInfo struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	LoginToken struct {
		Token string `json:"token"`
	}
)

type SignUpInfo struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

func (s *SignUpInfo) Validate() *ErrorMessage {
	if s.FirstName == "" || s.LastName == "" {
		return ErrorName
	}

	if !utils.IsEmailValid(s.Email) {
		return ErrorEmail
	}

	if _, err := utils.IsPasswordStrong(s.Password); err != nil {
		return &ErrorMessage{
			Message: err.Error(),
		}
	}

	return nil
}
