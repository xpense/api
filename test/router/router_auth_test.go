package router

import (
	"bytes"
	"encoding/json"
	"errors"
	"expense-api/internal/handlers"
	"expense-api/internal/model"
	"expense-api/internal/repository"
	"expense-api/internal/router"
	"expense-api/internal/utils"
	"expense-api/test/spies"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

const baseAuthPath = "/api/v1/auth"

func TestSignUp(t *testing.T) {
	repoSpy := &spies.RepositorySpy{}
	jwtServiceSpy := &spies.JWTServiceSpy{}
	hasherSpy := &spies.PasswordHasherSpy{}

	r := router.Setup(repoSpy, jwtServiceSpy, hasherSpy, router.TestConfig)

	newSignUpRequest := func(signUp *handlers.SignUpInfo) *http.Request {
		body := createRequestBody(signUp)
		req, _ := http.NewRequest(http.MethodPost, baseAuthPath+"/signup", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		return req
	}

	t.Run("Shouldn't sign up with missing 'first_name'", func(t *testing.T) {
		res := httptest.NewRecorder()
		req := newSignUpRequest(&handlers.SignUpInfo{})

		r.ServeHTTP(res, req)

		wantErrorMessage := handlers.ErrorName.Error()

		assertStatusCode(t, res, http.StatusBadRequest)
		assertErrorMessage(t, res, wantErrorMessage)
	})

	t.Run("Shouldn't sign up with missing 'last_name'", func(t *testing.T) {
		res := httptest.NewRecorder()
		req := newSignUpRequest(&handlers.SignUpInfo{
			FirstName: "First Name",
		})

		r.ServeHTTP(res, req)

		wantErrorMessage := handlers.ErrorName.Error()

		assertStatusCode(t, res, http.StatusBadRequest)
		assertErrorMessage(t, res, wantErrorMessage)
	})

	t.Run("Shouldn't sign up with missing 'email'", func(t *testing.T) {
		res := httptest.NewRecorder()
		req := newSignUpRequest(&handlers.SignUpInfo{
			FirstName: "First Name",
			LastName:  "Last Name",
		})

		r.ServeHTTP(res, req)

		wantErrorMessage := handlers.ErrorEmail.Error()

		assertStatusCode(t, res, http.StatusBadRequest)
		assertErrorMessage(t, res, wantErrorMessage)
	})

	t.Run("Shouldn't sign up with missing 'password'", func(t *testing.T) {
		res := httptest.NewRecorder()
		req := newSignUpRequest(&handlers.SignUpInfo{
			FirstName: "First Name",
			LastName:  "Last Name",
			Email:     "john@doe.com",
		})

		r.ServeHTTP(res, req)

		wantErrorMessage := utils.ErrorPasswordLength.Error()

		assertStatusCode(t, res, http.StatusBadRequest)
		assertErrorMessage(t, res, wantErrorMessage)
	})

	t.Run("Shouldn't sign up user with an already registered email", func(t *testing.T) {
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
		repoSpy.On("UserCreate", user.FirstName, user.LastName, user.Email, hashedPassword, salt).Return(nil, repository.ErrorUniqueConstaintViolation).Once()

		res := httptest.NewRecorder()
		req := newSignUpRequest(&handlers.SignUpInfo{
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
			Password:  user.Password,
		})

		r.ServeHTTP(res, req)

		wantErrorMessage := handlers.ErrorEmailConflict.Error()

		assertStatusCode(t, res, http.StatusConflict)
		assertErrorMessage(t, res, wantErrorMessage)
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
		req := newSignUpRequest(&handlers.SignUpInfo{
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
			Password:  user.Password,
		})

		r.ServeHTTP(res, req)

		assertStatusCode(t, res, http.StatusCreated)
	})
}

func TestLogin(t *testing.T) {
	repoSpy := &spies.RepositorySpy{}
	jwtServiceSpy := &spies.JWTServiceSpy{}
	hasherSpy := &spies.PasswordHasherSpy{}

	r := router.Setup(repoSpy, jwtServiceSpy, hasherSpy, router.TestConfig)

	newLoginRequest := func(login *handlers.LoginInfo) *http.Request {
		body := createRequestBody(login)
		req, _ := http.NewRequest(http.MethodPost, baseAuthPath+"/login", bytes.NewReader(body))
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

				wantErrorMessage := handlers.ErrorMissingPasswordOrEmail.Error()

				assertStatusCode(t, res, http.StatusBadRequest)
				assertErrorMessage(t, res, wantErrorMessage)
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

		wantErrorMessage := handlers.ErrorNonExistentUser.Error()

		assertStatusCode(t, res, http.StatusNotFound)
		assertErrorMessage(t, res, wantErrorMessage)
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
		hasherSpy.On("HashPassword", reqBody.Password, user.Salt).Return("", errors.New("dummy error")).Once()

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

		wantErrorMessage := handlers.ErrorWrongPassword.Error()

		assertStatusCode(t, res, http.StatusBadRequest)
		assertErrorMessage(t, res, wantErrorMessage)
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
		user.ID = 1

		repoSpy.On("UserGetWithEmail", reqBody.Email).Return(user, nil).Once()
		hasherSpy.On("HashPassword", reqBody.Password, user.Salt).Return(user.Password, nil).Once()
		jwtServiceSpy.On("CreateJWT", user.ID, user.Email).Return("", errors.New("dummy error")).Once()

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
		user.ID = 1
		loginToken := &handlers.LoginToken{
			Token: "token",
		}

		repoSpy.On("UserGetWithEmail", reqBody.Email).Return(user, nil).Once()
		hasherSpy.On("HashPassword", reqBody.Password, user.Salt).Return(user.Password, nil).Once()
		jwtServiceSpy.On("CreateJWT", user.ID, user.Email).Return(loginToken.Token, nil).Once()

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
