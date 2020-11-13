package handlers

import (
	"expense-api/model"
	"expense-api/repository"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type UserHandler interface {
	GetUser(ctx *gin.Context)
	UpdateUserInfo(ctx *gin.Context)
	DeleteUser(ctx *gin.Context)
}

// User is a user with an omitted 'password' field
type User struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Email     string    `json:"email"`
}

// UserModelToResponse cretes a user struct that doesn't expose the password of a user
func UserModelToResponse(u *model.User) *User {
	return &User{
		ID:        u.ID,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Email:     u.Email,
	}
}

func (h *handler) UpdateUserInfo(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		ctx.Status(http.StatusBadRequest)
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

	userModel, err := h.repo.UserUpdate(uint(id), userBody.FirstName, userBody.LastName, userBody.Email)
	if err != nil {
		if err == repository.ErrorRecordNotFound {
			ctx.Status(http.StatusNotFound)
			return
		}

		ctx.Status(http.StatusInternalServerError)
		return
	}

	userResponse := UserModelToResponse(userModel)
	ctx.JSON(http.StatusOK, userResponse)
}

func (h *handler) DeleteUser(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		ctx.Status(http.StatusBadRequest)
		return
	}

	if err := h.repo.UserDelete(uint(id)); err != nil {
		if err == repository.ErrorRecordNotFound {
			ctx.Status(http.StatusNotFound)
			return
		}

		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (h *handler) GetUser(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		ctx.Status(http.StatusBadRequest)
		return
	}

	userModel, err := h.repo.UserGet(uint(id))
	if err != nil {
		if err == repository.ErrorRecordNotFound {
			ctx.Status(http.StatusNotFound)
			return
		}

		ctx.Status(http.StatusInternalServerError)
		return
	}

	userResponse := UserModelToResponse(userModel)
	ctx.JSON(http.StatusOK, userResponse)
}
