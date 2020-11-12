package handlers

import (
	"expense-api/model"
	"expense-api/repository"
	"expense-api/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SignUp(repo repository.Repository, hasher utils.PasswordHasher) func(*gin.Context) {
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

type LoginInfo struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Login(repo repository.Repository, hasher utils.PasswordHasher) func(*gin.Context) {
	return func(ctx *gin.Context) {
		var loginInfo LoginInfo
		if err := ctx.Bind(&loginInfo); err != nil {
			ctx.Status(http.StatusBadRequest)
			return
		}

		if loginInfo.Email == "" || loginInfo.Password == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "both email and password are required for login",
			})
			return
		}

		user, err := repo.UserGetWithEmail(loginInfo.Email)
		if err != nil {
			if err == repository.ErrorRecordNotFound {
				ctx.JSON(http.StatusNotFound, gin.H{
					"message": "user with this email does not exist",
				})
				return
			}
			ctx.Status(http.StatusInternalServerError)
			return
		}

		hashedPassword, err := hasher.HashPassword(loginInfo.Password, user.Salt)
		if err != nil {
			ctx.Status(http.StatusInternalServerError)
			return
		}

		if user.Password != hashedPassword {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "wrong password",
			})
			return
		}

		token, err := utils.CreateJWT(user.Email)

		if err != nil {
			ctx.Status(http.StatusInternalServerError)
			return
		}

		response := struct {
			Token string `json:"token"`
		}{token}

		ctx.JSON(http.StatusOK, response)
	}
}
