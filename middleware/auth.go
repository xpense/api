package middleware

import (
	"expense-api/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const claimsKey = "claims"

func Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := strings.Split(ctx.GetHeader("Authorization"), "Bearer ")
		if len(authHeader) != 2 {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "Malformed token",
			})
			return
		}

		claims, err := utils.ValidateJWT(authHeader[1])
		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		ctx.Set(claimsKey, claims)
		ctx.Next()
	}
}

func GetClaimsFromContext(ctx *gin.Context) *utils.CustomClaims {
	val, exists := ctx.Get(claimsKey)
	if !exists {
		return nil
	}

	claims, ok := val.(*utils.CustomClaims)
	if !ok {
		return nil
	}

	return claims
}
