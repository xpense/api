package test

import (
	"bytes"
	"encoding/json"
	"errors"
	"expense-api/handlers"
	"expense-api/model"
	"expense-api/repository"
	"expense-api/router"
	"expense-api/router/test/spies"
	"expense-api/utils"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestSignUp(t *testing.T) {
	repoSpy := &spies.RepositorySpy{}
	jwtServiceSpy := &spies.JWTServiceSpy{}
	hasherSpy := &spies.PasswordHasherSpy{}

	r := router.Setup(repoSpy, jwtServiceSpy, hasherSpy)

	newUserRequest := func(user *model.User) *http.Request {
		body := createRequestBody(user)
		req, _ := http.NewRequest(http.MethodPost, "/auth/signup", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		return req
	}

	t.Run("Shouldn't sign up with missing 'first_name'", func(t *testing.T) {
		res := httptest.NewRecorder()
		req := newUserRequest(&model.User{})

		r.ServeHTTP(res, req)

		jsonResponse := parseJSON(t, res)

		haveErrorMessage := jsonResponse["message"].(string)
		wantErrorMessage := model.ErrorName.Error()

		assertStatusCode(t, res, http.StatusBadRequest)
		assertErrorMessage(t, haveErrorMessage, wantErrorMessage)
	})

	t.Run("Shouldn't sign up with missing 'last_name'", func(t *testing.T) {
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

	t.Run("Shouldn't sign up with missing 'email'", func(t *testing.T) {
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

	t.Run("Shouldn't sign up with missing 'password'", func(t *testing.T) {
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

	t.Run("Should sign up user", func(t *testing.T) {
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

func TestLogin(t *testing.T) {
	repoSpy := &spies.RepositorySpy{}
	jwtServiceSpy := &spies.JWTServiceSpy{}
	hasherSpy := &spies.PasswordHasherSpy{}

	r := router.Setup(repoSpy, jwtServiceSpy, hasherSpy)

	newLoginRequest := func(login *handlers.LoginInfo) *http.Request {
		body := createRequestBody(login)
		req, _ := http.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		return req
	}

	t.Run("Invalid request body", func(t *testing.T) {
		testCases := []struct {
			desc    string
			reqBody *handlers.LoginInfo
		}{
			{
				desc:    "Shouldn't allow login with empty body",
				reqBody: &handlers.LoginInfo{},
			},
			{
				desc:    "Shouldn't allow login with missing password",
				reqBody: &handlers.LoginInfo{Email: "john@doe.com"},
			},
			{
				desc:    "Shouldn't allow login with missing email",
				reqBody: &handlers.LoginInfo{Password: "123Password!{}"},
			},
		}

		for _, tC := range testCases {
			t.Run(tC.desc, func(t *testing.T) {
				res := httptest.NewRecorder()
				req := newLoginRequest(tC.reqBody)

				r.ServeHTTP(res, req)

				jsonResponse := parseJSON(t, res)
				haveErrorMessage := jsonResponse["message"].(string)
				wantErrorMessage := handlers.ErrMsgMissingPasswordOrEmail

				assertStatusCode(t, res, http.StatusBadRequest)
				assertErrorMessage(t, haveErrorMessage, wantErrorMessage)
			})
		}
	})

	t.Run("Shouldn't log in non-existent user", func(t *testing.T) {
		reqBody := &handlers.LoginInfo{
			Email:    "john@doe.com",
			Password: "123Password!{}",
		}

		repoSpy.On("UserGetWithEmail", reqBody.Email).Return(nil, repository.ErrorRecordNotFound).Once()

		res := httptest.NewRecorder()
		req := newLoginRequest(reqBody)

		r.ServeHTTP(res, req)

		jsonResponse := parseJSON(t, res)
		haveErrorMessage := jsonResponse["message"].(string)
		wantErrorMessage := handlers.ErrMsgNonExistentUser

		assertStatusCode(t, res, http.StatusNotFound)
		assertErrorMessage(t, haveErrorMessage, wantErrorMessage)
	})

	t.Run("Shouldn't log in if an error occurs while trying to query for user", func(t *testing.T) {
		reqBody := &handlers.LoginInfo{
			Email:    "john@doe.com",
			Password: "123Password!{}",
		}

		repoSpy.On("UserGetWithEmail", reqBody.Email).Return(nil, errors.New("dummy error")).Once()

		res := httptest.NewRecorder()
		req := newLoginRequest(reqBody)

		r.ServeHTTP(res, req)

		assertStatusCode(t, res, http.StatusInternalServerError)
	})

	t.Run("Shouldn't log in if an error occurs while trying to hash password", func(t *testing.T) {
		reqBody := &handlers.LoginInfo{
			Email:    "john@doe.com",
			Password: "123Password!{}",
		}
		user := &model.User{Salt: "salty"}

		repoSpy.On("UserGetWithEmail", reqBody.Email).Return(user, nil).Once()
		hasherSpy.On("HashPassword", reqBody.Password, user.Salt).Return(user, errors.New("dummy error")).Once()

		res := httptest.NewRecorder()
		req := newLoginRequest(reqBody)

		r.ServeHTTP(res, req)

		assertStatusCode(t, res, http.StatusInternalServerError)
	})

	t.Run("Shouldn't log in if there's a password mismatch", func(t *testing.T) {
		reqBody := &handlers.LoginInfo{
			Email:    "john@doe.com",
			Password: "123Password!{}",
		}
		user := &model.User{
			Salt:     "salty",
			Password: "good-password",
		}

		repoSpy.On("UserGetWithEmail", reqBody.Email).Return(user, nil).Once()
		hasherSpy.On("HashPassword", reqBody.Password, user.Salt).Return("bad-password", nil).Once()

		res := httptest.NewRecorder()
		req := newLoginRequest(reqBody)

		r.ServeHTTP(res, req)

		jsonResponse := parseJSON(t, res)
		haveErrorMessage := jsonResponse["message"].(string)
		wantErrorMessage := handlers.ErrMsgWrongPassword

		assertStatusCode(t, res, http.StatusBadRequest)
		assertErrorMessage(t, haveErrorMessage, wantErrorMessage)
	})

	t.Run("Shouldn't log in if there's an error while generating the access token", func(t *testing.T) {
		reqBody := &handlers.LoginInfo{
			Email:    "john@doe.com",
			Password: "123Password!{}",
		}
		user := &model.User{
			Email:    "john@doe.com",
			Salt:     "salty",
			Password: "good-password",
		}

		repoSpy.On("UserGetWithEmail", reqBody.Email).Return(user, nil).Once()
		hasherSpy.On("HashPassword", reqBody.Password, user.Salt).Return(user.Password, nil).Once()
		jwtServiceSpy.On("CreateJWT", user.Email).Return(nil, errors.New("dummy error")).Once()

		res := httptest.NewRecorder()
		req := newLoginRequest(reqBody)

		r.ServeHTTP(res, req)

		assertStatusCode(t, res, http.StatusInternalServerError)
	})

	t.Run("Should log in user and return access token", func(t *testing.T) {
		reqBody := &handlers.LoginInfo{
			Email:    "john@doe.com",
			Password: "123Password!{}",
		}
		user := &model.User{
			Email:    "john@doe.com",
			Salt:     "salty",
			Password: "good-password",
		}
		loginToken := &handlers.LoginToken{
			Token: "token",
		}

		repoSpy.On("UserGetWithEmail", reqBody.Email).Return(user, nil).Once()
		hasherSpy.On("HashPassword", reqBody.Password, user.Salt).Return(user.Password, nil).Once()
		jwtServiceSpy.On("CreateJWT", user.Email).Return(loginToken.Token, nil).Once()

		res := httptest.NewRecorder()
		req := newLoginRequest(reqBody)

		r.ServeHTTP(res, req)

		assertStatusCode(t, res, http.StatusOK)
		assertLoginTokenResponseBody(t, res, loginToken)
	})
}

func assertLoginTokenResponseBody(t *testing.T, res *httptest.ResponseRecorder, expected *handlers.LoginToken) {
	t.Helper()

	var got handlers.LoginToken
	if err := json.NewDecoder(res.Body).Decode(&got); err != nil {
		t.Errorf("couldn't parse json response: %v", err)
	}

	if !reflect.DeepEqual(got, *expected) {
		t.Errorf("expected %+v ;%T, got %+v ;%T", *expected, *expected, got, got)
	}
}
