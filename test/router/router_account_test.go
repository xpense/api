package router

import (
	"expense-api/internal/handlers"
	auth_middleware "expense-api/internal/middleware/auth"
	"expense-api/internal/model"
	"expense-api/internal/repository"
	"expense-api/internal/router"
	"expense-api/test/spies"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetAccount(t *testing.T) {
	repoSpy := &spies.RepositorySpy{}
	jwtServiceSpy := &spies.JWTServiceSpy{}
	hasherSpy := &spies.PasswordHasherSpy{}

	r := router.Setup(repoSpy, jwtServiceSpy, hasherSpy, router.TestConfig)

	t.Run("Missing/Invalid authorization token cases", func(t *testing.T) {
		token := "invalid-token"

		missingTokenReq := NewGetAccountRequest(token)
		invalidTokenReq := NewGetAccountRequest(token)

		unauthorizedTestCases := UnauthorizedTestCases(missingTokenReq, invalidTokenReq, r, jwtServiceSpy)
		t.Run("Unauthorized test cases", unauthorizedTestCases)
	})

	t.Run("Valid authorization token cases", func(t *testing.T) {
		claims := &auth_middleware.CustomClaims{
			ID:    1,
			Email: "john@doe.com",
		}
		token := "valid-token"

		jwtServiceSpy.On("ValidateJWT", token).Return(claims, nil)

		t.Run("Get non-existent user", func(t *testing.T) {
			repoSpy.On("UserGet", claims.ID).Return(nil, repository.ErrorRecordNotFound).Once()

			res := httptest.NewRecorder()
			req := NewGetAccountRequest(token)

			r.ServeHTTP(res, req)

			AssertStatusCode(t, res, http.StatusNotFound)
		})

		t.Run("Get existing user", func(t *testing.T) {
			user := &model.User{}

			repoSpy.On("UserGet", claims.ID).Return(user, nil).Once()

			res := httptest.NewRecorder()
			req := NewGetAccountRequest(token)

			r.ServeHTTP(res, req)

			AssertStatusCode(t, res, http.StatusOK)
			AssertResponseBody(t, res, &handlers.Account{})
		})
	})
}

func TestUpdateAccount(t *testing.T) {
	repoSpy := &spies.RepositorySpy{}
	jwtServiceSpy := &spies.JWTServiceSpy{}
	hasherSpy := &spies.PasswordHasherSpy{}

	r := router.Setup(repoSpy, jwtServiceSpy, hasherSpy, router.TestConfig)

	t.Run("Missing/Invalid authorization token cases", func(t *testing.T) {
		account := &handlers.Account{}
		token := "invalid-token"

		missingTokenReq := NewUpdateAccountRequest(account, token)
		invalidTokenReq := NewUpdateAccountRequest(account, token)

		unauthorizedTestCases := UnauthorizedTestCases(missingTokenReq, invalidTokenReq, r, jwtServiceSpy)
		t.Run("Unauthorized test cases", unauthorizedTestCases)
	})

	t.Run("Valid authorization token cases", func(t *testing.T) {
		claims := &auth_middleware.CustomClaims{
			ID:    10,
			Email: "john@doe.com",
		}
		token := "valid-token"

		jwtServiceSpy.On("ValidateJWT", token).Return(claims, nil)

		t.Run("Update non-existent user", func(t *testing.T) {
			account := &handlers.Account{
				FirstName: "Updated First Name",
				LastName:  "Last Name",
				Email:     "john@doe.com",
			}

			repoSpy.On("UserUpdate", claims.ID, account.FirstName, account.LastName, account.Email).Return(nil, repository.ErrorRecordNotFound).Once()

			res := httptest.NewRecorder()
			req := NewUpdateAccountRequest(account, token)

			r.ServeHTTP(res, req)

			AssertStatusCode(t, res, http.StatusNotFound)
		})

		t.Run("Update existing user with empty body", func(t *testing.T) {
			account := &handlers.Account{}

			res := httptest.NewRecorder()
			req := NewUpdateAccountRequest(account, token)

			r.ServeHTTP(res, req)

			wantErrorMessage := handlers.ErrorEmptyBody.Error()

			AssertStatusCode(t, res, http.StatusBadRequest)
			AssertErrorMessage(t, res, wantErrorMessage)
		})

		t.Run("Update existing user with invalid email", func(t *testing.T) {
			account := &handlers.Account{Email: "@"}

			res := httptest.NewRecorder()
			req := NewUpdateAccountRequest(account, token)

			r.ServeHTTP(res, req)

			wantErrorMessage := handlers.ErrorEmail.Error()

			AssertStatusCode(t, res, http.StatusBadRequest)
			AssertErrorMessage(t, res, wantErrorMessage)
		})

		t.Run("Update existing user with valid email", func(t *testing.T) {
			user := &model.User{Email: "john@doe.com"}
			account := &handlers.Account{Email: user.Email}

			repoSpy.On("UserUpdate", claims.ID, user.FirstName, user.LastName, user.Email).Return(user, nil).Once()

			res := httptest.NewRecorder()
			req := NewUpdateAccountRequest(account, token)

			r.ServeHTTP(res, req)

			AssertStatusCode(t, res, http.StatusOK)
			AssertResponseBody(t, res, account)
		})
	})
}

func TestDeleteAccount(t *testing.T) {
	repoSpy := &spies.RepositorySpy{}
	jwtServiceSpy := &spies.JWTServiceSpy{}
	hasherSpy := &spies.PasswordHasherSpy{}

	r := router.Setup(repoSpy, jwtServiceSpy, hasherSpy, router.TestConfig)

	t.Run("Missing/Invalid authorization token cases", func(t *testing.T) {
		token := "invalid-token"

		missingTokenReq := NewDeleteAccountRequest(token)
		invalidTokenReq := NewDeleteAccountRequest(token)

		unauthorizedTestCases := UnauthorizedTestCases(missingTokenReq, invalidTokenReq, r, jwtServiceSpy)
		t.Run("Unauthorized test cases", unauthorizedTestCases)
	})

	t.Run("Valid authorization token cases", func(t *testing.T) {
		claims := &auth_middleware.CustomClaims{
			ID:    10,
			Email: "john@doe.com",
		}
		token := "valid-token"

		jwtServiceSpy.On("ValidateJWT", token).Return(claims, nil)

		t.Run("Delete non-existent user", func(t *testing.T) {
			repoSpy.On("UserDelete", claims.ID).Return(repository.ErrorRecordNotFound).Once()

			res := httptest.NewRecorder()
			req := NewDeleteAccountRequest(token)

			r.ServeHTTP(res, req)

			AssertStatusCode(t, res, http.StatusNotFound)
		})

		t.Run("Delete existing user", func(t *testing.T) {
			repoSpy.On("UserDelete", claims.ID).Return(nil).Once()

			res := httptest.NewRecorder()
			req := NewDeleteAccountRequest(token)

			r.ServeHTTP(res, req)

			AssertStatusCode(t, res, http.StatusNoContent)
		})
	})
}
