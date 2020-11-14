package utils

import "github.com/gin-gonic/gin"

type CustomClaims struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
}

const claimsContextKey = "claims"

// GetClaimsFromContext extracts the CustomClaims from a gin Context
func GetClaimsFromContext(ctx *gin.Context) *CustomClaims {
	val, exists := ctx.Get(claimsContextKey)
	if !exists {
		return nil
	}

	claims, ok := val.(*CustomClaims)
	if !ok {
		return nil
	}

	return claims
}

// SetClaimsToContext sets the claims to a gin Context
func SetClaimsToContext(ctx *gin.Context, claims *CustomClaims) {
	ctx.Set(claimsContextKey, claims)
}
