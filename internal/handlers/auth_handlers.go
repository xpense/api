package handlers

import (
	"expense-api/internal/repository"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler interface {
	SignUp(ctx *gin.Context)
	Login(ctx *gin.Context)
}

func (h *handler) SignUp(ctx *gin.Context) {
	var signUpInfo SignUpInfo
	if err := ctx.Bind(&signUpInfo); err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

	if err := signUpInfo.Validate(); err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	salt, err := h.hasher.GenerateSalt()
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	hashedPassword, err := h.hasher.HashPassword(signUpInfo.Password, salt)
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	if _, err := h.repo.UserCreate(
		signUpInfo.FirstName,
		signUpInfo.LastName,
		signUpInfo.Email,
		hashedPassword,
		salt,
	); err != nil {
		if err == repository.ErrorUniqueConstaintViolation {
			ctx.JSON(http.StatusConflict, ErrorEmailConflict)
			return
		}

		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.Status(http.StatusCreated)

	// ctx.Status(http.StatusCreated)
}

func (h *handler) Login(ctx *gin.Context) {
	var loginInfo LoginInfo
	if err := ctx.Bind(&loginInfo); err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

	if loginInfo.Email == "" || loginInfo.Password == "" {
		ctx.JSON(http.StatusBadRequest, ErrorMissingPasswordOrEmail)
		return
	}

	user, err := h.repo.UserGetWithEmail(loginInfo.Email)
	if err != nil {
		if err == repository.ErrorRecordNotFound {
			ctx.JSON(http.StatusNotFound, ErrorNonExistentUser)
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
		ctx.JSON(http.StatusBadRequest, ErrorWrongPassword)
		return
	}

	token, err := h.jwtService.CreateJWT(user.ID, user.Email)
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	response := &LoginToken{Token: token}
	ctx.JSON(http.StatusOK, response)
}
