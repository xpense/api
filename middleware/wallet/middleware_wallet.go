package wallet

import (
	"expense-api/middleware"
	auth_middleware "expense-api/middleware/auth"
	"expense-api/repository"
	"net/http"

	"github.com/gin-gonic/gin"
)

type WalletMiddleware interface {
	ValidateOwnership(*gin.Context)
}

type walletMiddleware struct {
	repo repository.Repository
}

func New(repo repository.Repository) WalletMiddleware {
	return &walletMiddleware{repo}
}

func (w *walletMiddleware) ValidateOwnership(ctx *gin.Context) {
	userID, err := auth_middleware.GetUserIDFromContext(ctx)
	if err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	id := middleware.GetIDParamFromContext(ctx)

	wModel, err := w.repo.WalletGet(id)
	if err != nil {
		if err == repository.ErrorRecordNotFound {
			ctx.AbortWithStatus(http.StatusNotFound)
			return
		}
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if wModel.UserID != userID {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	ctx.Next()
}
