package wallet

import (
	"expense-api/internal/middleware"
	auth_middleware "expense-api/internal/middleware/auth"
	"expense-api/internal/repository"
	"net/http"

	"github.com/gin-gonic/gin"
)

type WalletsMiddleware interface {
	ValidateOwnership(*gin.Context)
}

type walletsMiddleware struct {
	repo repository.Repository
}

func New(repo repository.Repository) WalletsMiddleware {
	return &walletsMiddleware{repo}
}

func (w *walletsMiddleware) ValidateOwnership(ctx *gin.Context) {
	userID, err := auth_middleware.GetUserIDFromContext(ctx)
	if err != nil {
		ctx.AbortWithStatus(http.StatusForbidden)
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
		ctx.AbortWithStatus(http.StatusForbidden)
		return
	}

	ctx.Next()
}
