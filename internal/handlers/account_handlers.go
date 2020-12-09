package handlers

import (
	"expense-api/internal/middleware/auth"
	"expense-api/internal/repository"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AccountHandler interface {
	GetAccount(ctx *gin.Context)
	UpdateAccount(ctx *gin.Context)
	DeleteAccount(ctx *gin.Context)
}

func (h *handler) UpdateAccount(ctx *gin.Context) {
	id, err := auth.GetUserIDFromContext(ctx)
	if err != nil {
		ctx.AbortWithStatus(http.StatusForbidden)
		return
	}

	var accountBody Account
	if err := ctx.Bind(&accountBody); err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

	if err := accountBody.ValidateUpdateBody(); err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	userModel, err := h.repo.UserUpdate(
		id,
		accountBody.FirstName,
		accountBody.LastName,
		accountBody.Email,
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
		ctx.AbortWithStatus(http.StatusForbidden)
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
		ctx.AbortWithStatus(http.StatusForbidden)
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
