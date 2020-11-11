package handlers

import (
	"expense-api/model"
	"expense-api/repository"
	"expense-api/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type User struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Email     string    `json:"email"`
}

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

func CreateUser(repo repository.Repository, hasher utils.PasswordHasher) func(*gin.Context) {
	return func(ctx *gin.Context) {
		var userBody model.User
		if err := ctx.Bind(&userBody); err != nil {
			ctx.Status(http.StatusBadRequest)
			return
		}

		if err := model.UserValidateCreateBody(
			userBody.FirstName,
			userBody.LastName,
			userBody.Email,
			userBody.Password,
		); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}

		salt, err := hasher.GenerateSalt()
		if err != nil {
			ctx.Status(http.StatusInternalServerError)
			return
		}

		hashedPassword, err := hasher.HashPassword(userBody.Password, salt)
		if err != nil {
			ctx.Status(http.StatusInternalServerError)
			return
		}

		userModel, err := repo.UserCreate(
			userBody.FirstName,
			userBody.LastName,
			userBody.Email,
			hashedPassword,
			salt,
		)
		if err != nil {
			ctx.Status(http.StatusInternalServerError)
			return
		}

		userResponse := UserModelToResponse(userModel)
		ctx.JSON(http.StatusCreated, userResponse)
	}
}

func UpdateUserInfo(repo repository.Repository) func(*gin.Context) {
	return func(ctx *gin.Context) {
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

		userModel, err := repo.UserUpdate(uint(id), userBody.FirstName, userBody.LastName, userBody.Email)
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
}

func DeleteUser(repo repository.Repository) func(*gin.Context) {
	return func(ctx *gin.Context) {
		idStr := ctx.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil || id <= 0 {
			ctx.Status(http.StatusBadRequest)
			return
		}

		if err := repo.UserDelete(uint(id)); err != nil {
			if err == repository.ErrorRecordNotFound {
				ctx.Status(http.StatusNotFound)
				return
			}

			ctx.Status(http.StatusInternalServerError)
			return
		}

		ctx.Status(http.StatusNoContent)
	}
}

func GetUser(repo repository.Repository) func(*gin.Context) {
	return func(ctx *gin.Context) {
		idStr := ctx.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil || id <= 0 {
			ctx.Status(http.StatusBadRequest)
			return
		}

		userModel, err := repo.UserGet(uint(id))
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
}
