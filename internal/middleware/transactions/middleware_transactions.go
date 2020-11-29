package transaction

import (
	"expense-api/internal/middleware"
	auth_middleware "expense-api/internal/middleware/auth"
	"expense-api/internal/repository"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TransactionsMiddleware interface {
	ValidateOwnership(*gin.Context)
}

type transactionsMiddleware struct {
	repo repository.Repository
}

func New(repo repository.Repository) TransactionsMiddleware {
	return &transactionsMiddleware{repo}
}

func (t *transactionsMiddleware) ValidateOwnership(ctx *gin.Context) {
	userID, err := auth_middleware.GetUserIDFromContext(ctx)
	if err != nil {
		ctx.AbortWithStatus(http.StatusForbidden)
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
		ctx.AbortWithStatus(http.StatusForbidden)
		return
	}

	ctx.Next()
}
