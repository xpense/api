package integration

import (
	"expense-api/internal/handlers"
	router_test "expense-api/test/router"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	r := Setup()

	var authToken string

	t.Run("Signs up a user and then logs that user in", func(t *testing.T) {
		firstName := "First Name"
		lastName := "Last Name"
		email := "john@doe.com"
		password := "123Password!{}"

		// Sign up

		signUpInfo := &handlers.SignUpInfo{
			FirstName: firstName,
			LastName:  lastName,
			Email:     email,
			Password:  password,
		}

		signUpReq := router_test.NewSignUpRequest(signUpInfo)
		signUpRes := httptest.NewRecorder()

		r.ServeHTTP(signUpRes, signUpReq)

		router_test.AssertStatusCode(t, signUpRes, http.StatusCreated)

		// Login

		loginInfo := &handlers.LoginInfo{
			Email:    email,
			Password: password,
		}

		loginReq := router_test.NewLoginRequest(loginInfo)
		loginRes := httptest.NewRecorder()

		r.ServeHTTP(loginRes, loginReq)

		router_test.AssertStatusCode(t, loginRes, http.StatusOK)

		jsonResponse := router_test.ParseJSON(t, loginRes)

		authToken = jsonResponse["token"].(string)
	})

	t.Run("Creates, updates, deletes, and lists wallets", func(t *testing.T) {
		// Create

		wallet := &handlers.Wallet{
			Name: "cash",
		}

		createWalletReq := router_test.NewCreateWalletRequest(wallet, authToken)
		createWalletRes := httptest.NewRecorder()

		r.ServeHTTP(createWalletRes, createWalletReq)

		router_test.AssertStatusCode(t, createWalletRes, http.StatusCreated)

		jsonResponse := router_test.ParseJSON(t, createWalletRes)

		id := uint(jsonResponse["id"].(float64))

		// Get

		getWalletReq := router_test.NewGetWalletRequest(id, authToken)
		getWalletRes := httptest.NewRecorder()

		r.ServeHTTP(getWalletRes, getWalletReq)

		router_test.AssertStatusCode(t, getWalletRes, http.StatusOK)

		jsonResponse = router_test.ParseJSON(t, getWalletRes)

		if walletName := jsonResponse["name"].(string); wallet.Name != walletName {
			t.Errorf("Expected wallet name: %s, got: %s", wallet.Name, walletName)
		}

		// Update

		updateWallet := &handlers.Wallet{
			Name: "groceries",
		}

		updateWalletReq := router_test.NewUpdateWalletRequest(id, updateWallet, authToken)
		updateWalletRes := httptest.NewRecorder()

		r.ServeHTTP(updateWalletRes, updateWalletReq)

		router_test.AssertStatusCode(t, updateWalletRes, http.StatusOK)

		jsonResponse = router_test.ParseJSON(t, updateWalletRes)

		if updateWalletName := jsonResponse["name"].(string); updateWallet.Name != updateWalletName {
			t.Errorf("Expected wallet name: %s, got: %s", updateWallet.Name, updateWalletName)
		}

		// List

		listWalletsReq := router_test.NewListWalletsRequest(authToken)
		listWalletsRes := httptest.NewRecorder()

		r.ServeHTTP(listWalletsRes, listWalletsReq)

		router_test.AssertStatusCode(t, listWalletsRes, http.StatusOK)

		jsonResponse = router_test.ParseJSON(t, listWalletsRes)

		if count := uint(jsonResponse["count"].(float64)); count != 1 {
			t.Errorf("Expected count: 1, got: %d", count)
		}

		// Delete

		deleteWalletReq := router_test.NewDeleteWalletRequest(id, authToken)
		deleteWalletRes := httptest.NewRecorder()

		r.ServeHTTP(deleteWalletRes, deleteWalletReq)

		router_test.AssertStatusCode(t, deleteWalletRes, http.StatusNoContent)
	})
}
