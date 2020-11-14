package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type CustomClaims struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
}

const claimsContextKey = "claims"

// GetIDFromContext retrieves the authenticated user's id from the context
func GetIDFromContext(ctx *gin.Context) uint {
	claims := GetClaimsFromContext(ctx)
	return claims.ID
}

// GetClaimsFromContext extracts the CustomClaims from a gin Context
func GetClaimsFromContext(ctx *gin.Context) *CustomClaims {
	val, exists := ctx.Get(claimsContextKey)
	if !exists {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return nil
	}

	claims, ok := val.(*CustomClaims)
	if !ok {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return nil
	}

	if claims == nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return nil
	}

	return claims
}

// SetClaimsToContext sets the claims to a gin Context
func SetClaimsToContext(ctx *gin.Context, claims *CustomClaims) {
	ctx.Set(claimsContextKey, claims)
}
