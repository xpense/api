package handlers

import (
	"expense-api/internal/model"
	"expense-api/internal/utils"
	"time"
)

// Account is a user with an omitted 'password' field
type Account struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
}

func (a *Account) ValidateUpdateBody() error {
	if a.FirstName == "" && a.LastName == "" && a.Email == "" {
		return ErrorEmptyBody
	}

	if a.Email != "" && !utils.IsEmailValid(a.Email) {
		return ErrorEmail
	}

	return nil
}

// UserModelToAccountResponse cretes a user struct that doesn't expose the password of a user
func UserModelToAccountResponse(u *model.User) *Account {
	return &Account{
		ID:        u.ID,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Email:     u.Email,
	}
}
