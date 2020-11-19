package test

import (
	"bytes"
	"encoding/json"
	"expense-api/handlers"
	auth_middleware "expense-api/middleware/auth"
	"expense-api/model"
	"expense-api/repository"
	"expense-api/router"
	"expense-api/router/test/spies"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestGetAccount(t *testing.T) {
	repoSpy := &spies.RepositorySpy{}
	jwtServiceSpy := &spies.JWTServiceSpy{}
	hasherSpy := &spies.PasswordHasherSpy{}

	r := router.Setup(repoSpy, jwtServiceSpy, hasherSpy, router.TestConfig)

	newAccountRequest := func(token string) *http.Request {
		req, _ := http.NewRequest(http.MethodGet, "/account/", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		return req
	}

	t.Run("Missing/Invalid authorization token cases", func(t *testing.T) {
		token := "invalid-token"

		missingTokenReq := newAccountRequest(token)
		invalidTokenReq := newAccountRequest(token)

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
			req := newAccountRequest(token)

			r.ServeHTTP(res, req)

			assertStatusCode(t, res, http.StatusNotFound)
		})

		t.Run("Get existing user", func(t *testing.T) {
			user := &model.User{}

			repoSpy.On("UserGet", claims.ID).Return(user, nil).Once()

			res := httptest.NewRecorder()
			req := newAccountRequest(token)

			r.ServeHTTP(res, req)

			assertStatusCode(t, res, http.StatusOK)
			assertUserResponseBody(t, res, user)
		})
	})
}

func TestUpdateAccount(t *testing.T) {
	repoSpy := &spies.RepositorySpy{}
	jwtServiceSpy := &spies.JWTServiceSpy{}
	hasherSpy := &spies.PasswordHasherSpy{}

	r := router.Setup(repoSpy, jwtServiceSpy, hasherSpy, router.TestConfig)

	newAccountRequest := func(user *model.User, token string) *http.Request {
		body := createRequestBody(user)
		req, _ := http.NewRequest(http.MethodPatch, "/account/", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		return req
	}

	t.Run("Missing/Invalid authorization token cases", func(t *testing.T) {
		user := &model.User{}
		token := "invalid-token"

		missingTokenReq := newAccountRequest(user, token)
		invalidTokenReq := newAccountRequest(user, token)

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
			user := &model.User{
				FirstName: "Updated First Name",
				LastName:  "Last Name",
				Email:     "john@doe.com",
			}

			repoSpy.On("UserUpdate", claims.ID, user.FirstName, user.LastName, user.Email).Return(nil, repository.ErrorRecordNotFound).Once()

			res := httptest.NewRecorder()
			req := newAccountRequest(user, token)

			r.ServeHTTP(res, req)

			assertStatusCode(t, res, http.StatusNotFound)
		})

		t.Run("Update existing user with empty body", func(t *testing.T) {
			user := &model.User{}

			res := httptest.NewRecorder()
			req := newAccountRequest(user, token)

			r.ServeHTTP(res, req)

			wantErrorMessage := model.ErrorEmptyBody.Error()

			assertStatusCode(t, res, http.StatusBadRequest)
			assertErrorMessage(t, res, wantErrorMessage)
		})

		t.Run("Update existing user with invalid email", func(t *testing.T) {
			user := &model.User{Email: "@"}

			res := httptest.NewRecorder()
			req := newAccountRequest(user, token)

			r.ServeHTTP(res, req)

			wantErrorMessage := model.ErrorEmail.Error()

			assertStatusCode(t, res, http.StatusBadRequest)
			assertErrorMessage(t, res, wantErrorMessage)
		})

		t.Run("Update existing user with valid email", func(t *testing.T) {
			user := &model.User{Email: "john@doe.com"}

			repoSpy.On("UserUpdate", claims.ID, user.FirstName, user.LastName, user.Email).Return(user, nil).Once()

			res := httptest.NewRecorder()
			req := newAccountRequest(user, token)

			r.ServeHTTP(res, req)

			assertStatusCode(t, res, http.StatusOK)
			assertUserResponseBody(t, res, user)
		})
	})
}

func TestDeleteAccount(t *testing.T) {
	repoSpy := &spies.RepositorySpy{}
	jwtServiceSpy := &spies.JWTServiceSpy{}
	hasherSpy := &spies.PasswordHasherSpy{}

	r := router.Setup(repoSpy, jwtServiceSpy, hasherSpy, router.TestConfig)

	newAccountRequest := func(token string) *http.Request {
		req, _ := http.NewRequest(http.MethodDelete, "/account/", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		return req
	}

	t.Run("Missing/Invalid authorization token cases", func(t *testing.T) {
		token := "invalid-token"

		missingTokenReq := newAccountRequest(token)
		invalidTokenReq := newAccountRequest(token)

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
			req := newAccountRequest(token)

			r.ServeHTTP(res, req)

			assertStatusCode(t, res, http.StatusNotFound)
		})

		t.Run("Delete existing user", func(t *testing.T) {
			repoSpy.On("UserDelete", claims.ID).Return(nil).Once()

			res := httptest.NewRecorder()
			req := newAccountRequest(token)

			r.ServeHTTP(res, req)

			assertStatusCode(t, res, http.StatusNoContent)
		})
	})
}

func assertUserResponseBody(t *testing.T, res *httptest.ResponseRecorder, user *model.User) {
	t.Helper()

	var got handlers.Account
	if err := json.NewDecoder(res.Body).Decode(&got); err != nil {
		t.Errorf("couldn't parse json response: %v", err)
	}

	expected := handlers.UserModelToAccountResponse(user)
	if !reflect.DeepEqual(got, *expected) {
		t.Errorf("expected %+v ;%T, got %+v ;%T", *expected, *expected, got, got)
	}
}
