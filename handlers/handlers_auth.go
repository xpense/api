package handlers

import (
	"expense-api/model"
	"expense-api/repository"
	"expense-api/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler interface {
	SignUp(ctx *gin.Context)
	Login(ctx *gin.Context)
}

type (
	LoginInfo struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	LoginToken struct {
		Token string `json:"token"`
	}
)

const (
	ErrMsgMissingPasswordOrEmail = "both email and password are required for login"
	ErrMsgNonExistentUser        = "user with this email does not exist"
	ErrMsgWrongPassword          = "wrong password"
)

func (h *handler) SignUp(ctx *gin.Context) {
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

	salt, err := h.hasher.GenerateSalt()
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	hashedPassword, err := h.hasher.HashPassword(userBody.Password, salt)
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	if _, err := h.repo.UserCreate(
		userBody.FirstName,
		userBody.LastName,
		userBody.Email,
		hashedPassword,
		salt,
	); err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.Status(http.StatusCreated)
}

func (h *handler) Login(ctx *gin.Context) {
	var loginInfo LoginInfo
	if err := ctx.Bind(&loginInfo); err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

	if loginInfo.Email == "" || loginInfo.Password == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": ErrMsgMissingPasswordOrEmail,
		})
		return
	}

	user, err := h.repo.UserGetWithEmail(loginInfo.Email)
	if err != nil {
		if err == repository.ErrorRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{
				"message": ErrMsgNonExistentUser,
			})
			return
		}
		ctx.Status(http.StatusInternalServerError)
		return
	}

	hashedPassword, err := h.hasher.HashPassword(loginInfo.Password, user.Salt)
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	if user.Password != hashedPassword {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": ErrMsgWrongPassword,
		})
		return
	}

	token, err := utils.CreateJWT(user.Email)
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	response := &LoginToken{token}
	ctx.JSON(http.StatusOK, response)
}
