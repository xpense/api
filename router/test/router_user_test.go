package test

import (
	"bytes"
	"encoding/json"
	"expense-api/handlers"
	"expense-api/model"
	"expense-api/repository"
	"expense-api/router"
	"expense-api/router/test/spies"
	"expense-api/utils"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestCreateUser(t *testing.T) {
	repoSpy := &spies.RepositorySpy{}
	hasherSpy := &spies.PasswordHasherSpy{}

	r := router.Setup(repoSpy, hasherSpy)

	newUserRequest := func(user *model.User) *http.Request {
		body := createRequestBody(user)
		req, _ := http.NewRequest(http.MethodPost, "/user/", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		return req
	}

	t.Run("Shouldn't create user with missing 'first_name'", func(t *testing.T) {
		res := httptest.NewRecorder()
		req := newUserRequest(&model.User{})

		r.ServeHTTP(res, req)

		jsonResponse := parseJSON(t, res)

		haveErrorMessage := jsonResponse["message"].(string)
		wantErrorMessage := model.ErrorName.Error()

		assertStatusCode(t, res, http.StatusBadRequest)
		assertErrorMessage(t, haveErrorMessage, wantErrorMessage)
	})

	t.Run("Shouldn't create user with missing 'last_name'", func(t *testing.T) {
		res := httptest.NewRecorder()
		req := newUserRequest(&model.User{
			FirstName: "First Name",
		})

		r.ServeHTTP(res, req)

		jsonResponse := parseJSON(t, res)

		haveErrorMessage := jsonResponse["message"].(string)
		wantErrorMessage := model.ErrorName.Error()

		assertStatusCode(t, res, http.StatusBadRequest)
		assertErrorMessage(t, haveErrorMessage, wantErrorMessage)
	})

	t.Run("Shouldn't create user with missing 'email'", func(t *testing.T) {
		res := httptest.NewRecorder()
		req := newUserRequest(&model.User{
			FirstName: "First Name",
			LastName:  "Last Name",
		})

		r.ServeHTTP(res, req)

		jsonResponse := parseJSON(t, res)

		haveErrorMessage := jsonResponse["message"].(string)
		wantErrorMessage := model.ErrorEmail.Error()

		assertStatusCode(t, res, http.StatusBadRequest)
		assertErrorMessage(t, haveErrorMessage, wantErrorMessage)
	})

	t.Run("Shouldn't create user with missing 'password'", func(t *testing.T) {
		res := httptest.NewRecorder()
		req := newUserRequest(&model.User{
			FirstName: "First Name",
			LastName:  "Last Name",
			Email:     "john@doe.com",
		})

		r.ServeHTTP(res, req)

		jsonResponse := parseJSON(t, res)

		haveErrorMessage := jsonResponse["message"].(string)
		wantErrorMessage := utils.ErrorPasswordLength.Error()

		assertStatusCode(t, res, http.StatusBadRequest)
		assertErrorMessage(t, haveErrorMessage, wantErrorMessage)
	})

	t.Run("Should create user", func(t *testing.T) {
		user := &model.User{
			FirstName: "First Name",
			LastName:  "Last Name",
			Email:     "john@doe.com",
			Password:  "123Password!{}",
		}
		user.ID = 1

		salt := "saltystring"
		hashedPassword := "hashedPassword"

		hasherSpy.On("GenerateSalt").Return(salt, nil).Once()
		hasherSpy.On("HashPassword", user.Password, salt).Return(hashedPassword, nil).Once()
		repoSpy.On("UserCreate", user.FirstName, user.LastName, user.Email, hashedPassword, salt).Return(user, nil).Once()

		res := httptest.NewRecorder()
		req := newUserRequest(user)

		r.ServeHTTP(res, req)

		assertStatusCode(t, res, http.StatusCreated)
	})
}

func TestGetUser(t *testing.T) {
	repoSpy := &spies.RepositorySpy{}
	hasherSpy := &spies.PasswordHasherSpy{}

	r := router.Setup(repoSpy, hasherSpy)

	newUserRequest := func(id uint) *http.Request {
		url := fmt.Sprintf("/user/%d", id)
		req, _ := http.NewRequest(http.MethodGet, url, nil)
		return req
	}

	t.Run("Get non-existent user", func(t *testing.T) {
		id := uint(1)
		repoSpy.On("UserGet", id).Return(nil, repository.ErrorRecordNotFound).Once()

		res := httptest.NewRecorder()
		req := newUserRequest(id)

		r.ServeHTTP(res, req)

		assertStatusCode(t, res, http.StatusNotFound)
	})

	t.Run("Get existing user", func(t *testing.T) {
		id := uint(1)
		user := &model.User{}

		repoSpy.On("UserGet", id).Return(user, nil).Once()

		res := httptest.NewRecorder()
		req := newUserRequest(id)

		r.ServeHTTP(res, req)

		assertStatusCode(t, res, http.StatusOK)
		assertUserResponseBody(t, res, user)
	})
}

func TestUpdateUser(t *testing.T) {
	repoSpy := &spies.RepositorySpy{}
	hasherSpy := &spies.PasswordHasherSpy{}

	r := router.Setup(repoSpy, hasherSpy)

	newUserRequest := func(id uint, user *model.User) *http.Request {
		url := fmt.Sprintf("/user/%d", id)
		body := createRequestBody(user)
		req, _ := http.NewRequest(http.MethodPatch, url, bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		return req
	}

	t.Run("Update non-existent user", func(t *testing.T) {
		id := uint(10)
		user := &model.User{
			FirstName: "Updated First Name",
			LastName:  "Last Name",
			Email:     "john@doe.com",
		}

		repoSpy.On("UserUpdate", id, user.FirstName, user.LastName, user.Email).Return(nil, repository.ErrorRecordNotFound).Once()

		res := httptest.NewRecorder()
		req := newUserRequest(id, user)

		r.ServeHTTP(res, req)

		assertStatusCode(t, res, http.StatusNotFound)
	})

	t.Run("Update existing user with empty body", func(t *testing.T) {
		id := uint(10)
		user := &model.User{}

		res := httptest.NewRecorder()
		req := newUserRequest(id, user)

		r.ServeHTTP(res, req)

		jsonResponse := parseJSON(t, res)

		haveErrorMessage := jsonResponse["message"].(string)
		wantErrorMessage := model.ErrorEmptyBody.Error()

		assertStatusCode(t, res, http.StatusBadRequest)
		assertErrorMessage(t, haveErrorMessage, wantErrorMessage)
	})

	t.Run("Update existing user with invalid email", func(t *testing.T) {
		id := uint(10)
		user := &model.User{Email: "@"}

		res := httptest.NewRecorder()
		req := newUserRequest(id, user)

		r.ServeHTTP(res, req)

		jsonResponse := parseJSON(t, res)

		haveErrorMessage := jsonResponse["message"].(string)
		wantErrorMessage := model.ErrorEmail.Error()

		assertStatusCode(t, res, http.StatusBadRequest)
		assertErrorMessage(t, haveErrorMessage, wantErrorMessage)
	})

	t.Run("Update existing user with valid email", func(t *testing.T) {
		id := uint(10)
		user := &model.User{Email: "john@doe.com"}

		repoSpy.On("UserUpdate", id, user.FirstName, user.LastName, user.Email).Return(user, nil).Once()

		res := httptest.NewRecorder()
		req := newUserRequest(id, user)

		r.ServeHTTP(res, req)

		assertStatusCode(t, res, http.StatusOK)
		assertUserResponseBody(t, res, user)
	})
}

func TestDeleteUser(t *testing.T) {
	repoSpy := &spies.RepositorySpy{}
	hasherSpy := &spies.PasswordHasherSpy{}

	r := router.Setup(repoSpy, hasherSpy)

	newUserRequest := func(id uint) *http.Request {
		url := fmt.Sprintf("/user/%d", id)
		req, _ := http.NewRequest(http.MethodDelete, url, nil)
		return req
	}

	t.Run("Delete non-existent user", func(t *testing.T) {
		id := uint(1)
		repoSpy.On("UserDelete", id).Return(repository.ErrorRecordNotFound).Once()

		res := httptest.NewRecorder()
		req := newUserRequest(id)

		r.ServeHTTP(res, req)

		assertStatusCode(t, res, http.StatusNotFound)
	})

	t.Run("Delete existing user", func(t *testing.T) {
		id := uint(1)
		repoSpy.On("UserDelete", id).Return(nil).Once()

		res := httptest.NewRecorder()
		req := newUserRequest(id)

		r.ServeHTTP(res, req)

		assertStatusCode(t, res, http.StatusNoContent)
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
