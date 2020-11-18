package transaction

import (
	"expense-api/middleware"
	auth_middleware "expense-api/middleware/auth"
	"expense-api/repository"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TransactionMiddleware interface {
	ValidateOwnership(*gin.Context)
}

type transactionMiddleware struct {
	repo repository.Repository
}

func New(repo repository.Repository) TransactionMiddleware {
	return &transactionMiddleware{repo}
}

func (t *transactionMiddleware) ValidateOwnership(ctx *gin.Context) {
	userID, err := auth_middleware.GetUserIDFromContext(ctx)
	if err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	tID := middleware.GetIDParamFromContext(ctx)

	tModel, err := t.repo.TransactionGet(tID)
	if err != nil {
		if err == repository.ErrorRecordNotFound {
			ctx.AbortWithStatus(http.StatusNotFound)
			return
		}
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if tModel.UserID != userID {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	ctx.Next()
}
