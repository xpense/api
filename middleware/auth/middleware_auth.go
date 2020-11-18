package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const ErrMsgMalformedToken = "Malformed token"

type AuthMiddleware interface {
	IsAuthenticated(*gin.Context)
}

type authMiddleware struct {
	jwtService JWTService
}

func New(jwtService JWTService) AuthMiddleware {
	return &authMiddleware{jwtService}
}

func (a *authMiddleware) IsAuthenticated(ctx *gin.Context) {
	authHeader := strings.Split(ctx.GetHeader("Authorization"), "Bearer ")
	if len(authHeader) != 2 {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": ErrMsgMalformedToken,
		})
		return
	}

	claims, err := a.jwtService.ValidateJWT(authHeader[1])
	if err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	ctx.Set(claimsContextKey, claims)
	ctx.Next()
}
