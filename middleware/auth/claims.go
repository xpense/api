package auth

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

// GetUserIDFromContext retrieves the authenticated user's id from the context
func GetUserIDFromContext(ctx *gin.Context) (uint, error) {
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
