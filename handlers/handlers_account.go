package handlers

import (
	"expense-api/middleware/auth"
	"expense-api/model"
	"expense-api/repository"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type AccountHandler interface {
	GetAccount(ctx *gin.Context)
	UpdateAccount(ctx *gin.Context)
	DeleteAccount(ctx *gin.Context)
}

// Account is a user with an omitted 'password' field
type Account struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Email     string    `json:"email"`
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

func (h *handler) UpdateAccount(ctx *gin.Context) {
	id, err := auth.GetUserIDFromContext(ctx)
	if err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	var userBody model.User
	if err := ctx.Bind(&userBody); err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

	if err := model.UserValidateUpdateBody(
		userBody.FirstName,
		userBody.LastName,
		userBody.Email,
	); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	userModel, err := h.repo.UserUpdate(
		id,
		userBody.FirstName,
		userBody.LastName,
		userBody.Email,
	)
	if err != nil {
		if err == repository.ErrorRecordNotFound {
			ctx.Status(http.StatusNotFound)
			return
		}

		ctx.Status(http.StatusInternalServerError)
		return
	}

	accountResponse := UserModelToAccountResponse(userModel)
	ctx.JSON(http.StatusOK, accountResponse)
}

func (h *handler) DeleteAccount(ctx *gin.Context) {
	id, err := auth.GetUserIDFromContext(ctx)
	if err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if err := h.repo.UserDelete(id); err != nil {
		if err == repository.ErrorRecordNotFound {
			ctx.Status(http.StatusNotFound)
			return
		}

		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (h *handler) GetAccount(ctx *gin.Context) {
	id, err := auth.GetUserIDFromContext(ctx)
	if err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	userModel, err := h.repo.UserGet(id)
	if err != nil {
		if err == repository.ErrorRecordNotFound {
			ctx.Status(http.StatusNotFound)
			return
		}

		ctx.Status(http.StatusInternalServerError)
		return
	}

	accountResponse := UserModelToAccountResponse(userModel)
	ctx.JSON(http.StatusOK, accountResponse)
}
