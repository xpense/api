package party

import (
	"expense-api/internal/middleware"
	auth_middleware "expense-api/internal/middleware/auth"
	"expense-api/internal/repository"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PartiesMiddleware interface {
	ValidateOwnership(*gin.Context)
}

type partiesMiddleware struct {
	repo repository.Repository
}

func New(repo repository.Repository) PartiesMiddleware {
	return &partiesMiddleware{repo}
}

func (p *partiesMiddleware) ValidateOwnership(ctx *gin.Context) {
	userID, err := auth_middleware.GetUserIDFromContext(ctx)
	if err != nil {
		ctx.AbortWithStatus(http.StatusForbidden)
		return
	}

	id := middleware.GetIDParamFromContext(ctx)

	wModel, err := p.repo.PartyGet(id)
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
