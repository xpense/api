package handlers

import (
	"expense-api/middleware/auth"
	"expense-api/repository"
	"expense-api/utils"
)

type Handler interface {
	TransactionsHandler
	AccountHandler
	AuthHandler
}

type handler struct {
	repo       repository.Repository
	jwtService auth.JWTService
	hasher     utils.PasswordHasher
}

func New(
	repo repository.Repository,
	jwtService auth.JWTService,
	hasher utils.PasswordHasher,
) Handler {
	return &handler{repo, jwtService, hasher}
}
