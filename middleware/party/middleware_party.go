package party

import (
	"expense-api/middleware"
	auth_middleware "expense-api/middleware/auth"
	"expense-api/repository"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PartyMiddleware interface {
	ValidateOwnership(*gin.Context)
}

type partyMiddleware struct {
	repo repository.Repository
}

func New(repo repository.Repository) PartyMiddleware {
	return &partyMiddleware{repo}
}

func (p *partyMiddleware) ValidateOwnership(ctx *gin.Context) {
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
