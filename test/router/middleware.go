package router

import (
	"errors"
	auth_middleware "expense-api/internal/middleware/auth"
	"expense-api/test/spies"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

// UnauthorizedTestCases runs two test cases for authentication middleware functionality
func UnauthorizedTestCases(
	missingTokenReq *http.Request,
	invalidTokenReq *http.Request,
	r *gin.Engine,
	jwtServiceSpy *spies.JWTServiceSpy,
) func(*testing.T) {
	return func(t *testing.T) {
		t.Run("Missing token", func(t *testing.T) {
			missingTokenReq.Header.Del("Authorization")

			res := httptest.NewRecorder()
			r.ServeHTTP(res, missingTokenReq)

			wantErrorMessage := auth_middleware.ErrMsgMalformedToken

			AssertStatusCode(t, res, http.StatusUnauthorized)
			AssertErrorMessage(t, res, wantErrorMessage)
		})

		t.Run("Invalid token", func(t *testing.T) {
			token := "invalid-token"
			jwtServiceSpy.On("ValidateJWT", token).Return(nil, errors.New("dummy error")).Once()

			res := httptest.NewRecorder()
			r.ServeHTTP(res, invalidTokenReq)

			AssertStatusCode(t, res, http.StatusUnauthorized)
		})
	}
}
