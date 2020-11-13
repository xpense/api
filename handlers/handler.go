package handlers

import (
	"expense-api/repository"
	"expense-api/utils"
)

type Handler interface {
	TransactionsHandler
	UserHandler
	AuthHandler
}

type handler struct {
	repo   repository.Repository
	hasher utils.PasswordHasher
}

func New(repo repository.Repository, hasher utils.PasswordHasher) Handler {
	return &handler{repo, hasher}
}
