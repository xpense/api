package test

import (
	"bytes"
	"encoding/json"
	"expense-api/handlers"
	"expense-api/model"
	"expense-api/repository"
	"expense-api/router"
	"expense-api/router/test/spies"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestGetUser(t *testing.T) {
	repoSpy := &spies.RepositorySpy{}
	jwtServiceSpy := &spies.JWTServiceSpy{}
	hasherSpy := &spies.PasswordHasherSpy{}

	r := router.Setup(repoSpy, jwtServiceSpy, hasherSpy)

	newUserRequest := func(id uint, token string) *http.Request {
		url := fmt.Sprintf("/user/%d", id)
		req, _ := http.NewRequest(http.MethodGet, url, nil)
		req.Header.Set("Authorization", "Bearer "+token)
		return req
	}

	t.Run("Missing/Invalid authorization token cases", func(t *testing.T) {
		id := uint(1)
		token := "invalid-token"

		missingTokenReq := newUserRequest(id, token)
		invalidTokenReq := newUserRequest(id, token)

		unauthorizedTestCases := UnauthorizedTestCases(missingTokenReq, invalidTokenReq, r, jwtServiceSpy)
		t.Run("Unauthorized test cases", unauthorizedTestCases)
	})

	t.Run("Valid authorization token cases", func(t *testing.T) {
		id := uint(1)
		token := "valid-token"

		jwtServiceSpy.On("ValidateJWT", token).Return(nil, nil)

		t.Run("Get non-existent user", func(t *testing.T) {
			repoSpy.On("UserGet", id).Return(nil, repository.ErrorRecordNotFound).Once()

			res := httptest.NewRecorder()
			req := newUserRequest(id, token)

			r.ServeHTTP(res, req)

			assertStatusCode(t, res, http.StatusNotFound)
		})

		t.Run("Get existing user", func(t *testing.T) {
			user := &model.User{}

			repoSpy.On("UserGet", id).Return(user, nil).Once()

			res := httptest.NewRecorder()
			req := newUserRequest(id, token)

			r.ServeHTTP(res, req)

			assertStatusCode(t, res, http.StatusOK)
			assertUserResponseBody(t, res, user)
		})
	})
}

func TestUpdateUser(t *testing.T) {
	repoSpy := &spies.RepositorySpy{}
	jwtServiceSpy := &spies.JWTServiceSpy{}
	hasherSpy := &spies.PasswordHasherSpy{}

	r := router.Setup(repoSpy, jwtServiceSpy, hasherSpy)

	newUserRequest := func(id uint, user *model.User, token string) *http.Request {
		url := fmt.Sprintf("/user/%d", id)
		body := createRequestBody(user)
		req, _ := http.NewRequest(http.MethodPatch, url, bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		return req
	}

	t.Run("Missing/Invalid authorization token cases", func(t *testing.T) {
		id := uint(1)
		user := &model.User{}
		token := "invalid-token"

		missingTokenReq := newUserRequest(id, user, token)
		invalidTokenReq := newUserRequest(id, user, token)

		unauthorizedTestCases := UnauthorizedTestCases(missingTokenReq, invalidTokenReq, r, jwtServiceSpy)
		t.Run("Unauthorized test cases", unauthorizedTestCases)
	})

	t.Run("Valid authorization token cases", func(t *testing.T) {
		id := uint(10)
		token := "valid-token"

		jwtServiceSpy.On("ValidateJWT", token).Return(nil, nil)

		t.Run("Update non-existent user", func(t *testing.T) {
			user := &model.User{
				FirstName: "Updated First Name",
				LastName:  "Last Name",
				Email:     "john@doe.com",
			}

			repoSpy.On("UserUpdate", id, user.FirstName, user.LastName, user.Email).Return(nil, repository.ErrorRecordNotFound).Once()

			res := httptest.NewRecorder()
			req := newUserRequest(id, user, token)

			r.ServeHTTP(res, req)

			assertStatusCode(t, res, http.StatusNotFound)
		})

		t.Run("Update existing user with empty body", func(t *testing.T) {
			user := &model.User{}

			res := httptest.NewRecorder()
			req := newUserRequest(id, user, token)

			r.ServeHTTP(res, req)

			wantErrorMessage := model.ErrorEmptyBody.Error()

			assertStatusCode(t, res, http.StatusBadRequest)
			assertErrorMessage(t, res, wantErrorMessage)
		})

		t.Run("Update existing user with invalid email", func(t *testing.T) {
			user := &model.User{Email: "@"}

			res := httptest.NewRecorder()
			req := newUserRequest(id, user, token)

			r.ServeHTTP(res, req)

			wantErrorMessage := model.ErrorEmail.Error()

			assertStatusCode(t, res, http.StatusBadRequest)
			assertErrorMessage(t, res, wantErrorMessage)
		})

		t.Run("Update existing user with valid email", func(t *testing.T) {
			user := &model.User{Email: "john@doe.com"}

			repoSpy.On("UserUpdate", id, user.FirstName, user.LastName, user.Email).Return(user, nil).Once()

			res := httptest.NewRecorder()
			req := newUserRequest(id, user, token)

			r.ServeHTTP(res, req)

			assertStatusCode(t, res, http.StatusOK)
			assertUserResponseBody(t, res, user)
		})
	})
}

func TestDeleteUser(t *testing.T) {
	repoSpy := &spies.RepositorySpy{}
	jwtServiceSpy := &spies.JWTServiceSpy{}
	hasherSpy := &spies.PasswordHasherSpy{}

	r := router.Setup(repoSpy, jwtServiceSpy, hasherSpy)

	newUserRequest := func(id uint, token string) *http.Request {
		url := fmt.Sprintf("/user/%d", id)
		req, _ := http.NewRequest(http.MethodDelete, url, nil)
		req.Header.Set("Authorization", "Bearer "+token)
		return req
	}

	t.Run("Missing/Invalid authorization token cases", func(t *testing.T) {
		id := uint(1)
		token := "invalid-token"

		missingTokenReq := newUserRequest(id, token)
		invalidTokenReq := newUserRequest(id, token)

		unauthorizedTestCases := UnauthorizedTestCases(missingTokenReq, invalidTokenReq, r, jwtServiceSpy)
		t.Run("Unauthorized test cases", unauthorizedTestCases)
	})

	t.Run("Valid authorization token cases", func(t *testing.T) {
		id := uint(1)
		token := "valid-token"

		jwtServiceSpy.On("ValidateJWT", token).Return(nil, nil)

		t.Run("Delete non-existent user", func(t *testing.T) {
			repoSpy.On("UserDelete", id).Return(repository.ErrorRecordNotFound).Once()

			res := httptest.NewRecorder()
			req := newUserRequest(id, token)

			r.ServeHTTP(res, req)

			assertStatusCode(t, res, http.StatusNotFound)
		})

		t.Run("Delete existing user", func(t *testing.T) {
			repoSpy.On("UserDelete", id).Return(nil).Once()

			res := httptest.NewRecorder()
			req := newUserRequest(id, token)

			r.ServeHTTP(res, req)

			assertStatusCode(t, res, http.StatusNoContent)
		})
	})
}

func assertUserResponseBody(t *testing.T, res *httptest.ResponseRecorder, user *model.User) {
	t.Helper()

	var got handlers.User
	if err := json.NewDecoder(res.Body).Decode(&got); err != nil {
		t.Errorf("couldn't parse json response: %v", err)
	}

	expected := handlers.UserModelToResponse(user)
	if !reflect.DeepEqual(got, *expected) {
		t.Errorf("expected %+v ;%T, got %+v ;%T", *expected, *expected, got, got)
	}
}
