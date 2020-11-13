package middleware

import (
	"expense-api/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type AuthMiddleware interface {
	Handler(*gin.Context)
}

type authMiddleware struct {
	jwtService utils.JWTService
}

func NewAuthMiddleware(jwtService utils.JWTService) AuthMiddleware {
	return &authMiddleware{jwtService}
}

func (a *authMiddleware) Handler(ctx *gin.Context) {
	authHeader := strings.Split(ctx.GetHeader("Authorization"), "Bearer ")
	if len(authHeader) != 2 {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "Malformed token",
		})
		return
	}

	claims, err := a.jwtService.ValidateJWT(authHeader[1])
	if err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	utils.SetClaimsToContext(ctx, claims)
	ctx.Next()
}
