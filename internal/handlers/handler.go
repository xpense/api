package handlers

import (
	"expense-api/internal/middleware/auth"
	"expense-api/internal/repository"
	"expense-api/internal/utils"
)

type Handler interface {
	AuthHandler
	AccountHandler
	TransactionsHandler
	WalletsHandler
	PartiesHandler
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
