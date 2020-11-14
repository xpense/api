package utils

import (
	"errors"

	"github.com/gin-gonic/gin"
)

type CustomClaims struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
}

var errNilCustomClaims = errors.New("custom claims not set")

const claimsContextKey = "claims"

// GetIDFromContext retrieves the authenticated user's id from the context
func GetIDFromContext(ctx *gin.Context) (uint, error) {
	claims, err := GetClaimsFromContext(ctx)
	if err != nil {
		return 0, err
	}

	return claims.ID, nil
}

// GetClaimsFromContext extracts the CustomClaims from a gin Context
func GetClaimsFromContext(ctx *gin.Context) (*CustomClaims, error) {
	val, exists := ctx.Get(claimsContextKey)
	if !exists {
		return nil, errNilCustomClaims
	}

	claims, ok := val.(*CustomClaims)
	if !ok || claims == nil {
		return nil, errNilCustomClaims
	}

	return claims, nil
}

// SetClaimsToContext sets the claims to a gin Context
func SetClaimsToContext(ctx *gin.Context, claims *CustomClaims) {
	ctx.Set(claimsContextKey, claims)
}
