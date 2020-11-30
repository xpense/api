package integration

import (
	"expense-api/internal/handlers"
	router_test "expense-api/test/router"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/shopspring/decimal"
)

func TestIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	r := Setup()

	var authToken string

	var (
		walletID      uint
		partyID       uint
		transactionID uint
	)

	var (
		firstName = "First Name"
		lastName  = "Last Name"
		email     = "john@doe.com"
		password  = "123Password!{}"
	)

	t.Run("Signs up a user and then logs that user in", func(t *testing.T) {
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

		var loginTokenResponseBody handlers.LoginToken
		router_test.ParseLoginTokenJSONResponse(t, loginRes, &loginTokenResponseBody)
		// Set auth token
		authToken = loginTokenResponseBody.Token
	})

	t.Run("Creates, updates, and lists wallets", func(t *testing.T) {
		{
			// Create wallet
			wallet := &handlers.Wallet{
				Name: "cash",
			}

			createWalletReq := router_test.NewCreateWalletRequest(wallet, authToken)
			createWalletRes := httptest.NewRecorder()

			r.ServeHTTP(createWalletRes, createWalletReq)
			router_test.AssertStatusCode(t, createWalletRes, http.StatusCreated)

			var createWalletResponseBody handlers.Wallet
			router_test.ParseWalletJSONResponse(t, createWalletRes, &createWalletResponseBody)

			// Set wallet id
			walletID = createWalletResponseBody.ID

			// Get wallet
			getWalletReq := router_test.NewGetWalletRequest(walletID, authToken)
			getWalletRes := httptest.NewRecorder()

			r.ServeHTTP(getWalletRes, getWalletReq)
			router_test.AssertStatusCode(t, getWalletRes, http.StatusOK)

			var getWalletResponseBody handlers.Wallet
			router_test.ParseWalletJSONResponse(t, getWalletRes, &getWalletResponseBody)

			if wallet.Name != getWalletResponseBody.Name {
				t.Errorf("Expected wallet name: %s, got: %s", wallet.Name, getWalletResponseBody.Name)
			}
		}

		{
			// Update wallet
			updateWallet := &handlers.Wallet{
				Name: "groceries",
			}

			updateWalletReq := router_test.NewUpdateWalletRequest(walletID, updateWallet, authToken)
			updateWalletRes := httptest.NewRecorder()

			r.ServeHTTP(updateWalletRes, updateWalletReq)
			router_test.AssertStatusCode(t, updateWalletRes, http.StatusOK)

			var updateWalletResponseBody handlers.Wallet
			router_test.ParseWalletJSONResponse(t, updateWalletRes, &updateWalletResponseBody)

			if updateWallet.Name != updateWalletResponseBody.Name {
				t.Errorf("Expected wallet name: %s, got: %s", updateWallet.Name, updateWalletResponseBody.Name)
			}
		}

		{
			// List wallets
			listWalletsReq := router_test.NewListWalletsRequest(authToken)
			listWalletsRes := httptest.NewRecorder()

			r.ServeHTTP(listWalletsRes, listWalletsReq)
			router_test.AssertStatusCode(t, listWalletsRes, http.StatusOK)

			var listWalletsResponseBody router_test.WalletListResponse
			router_test.ParseWalletListJSONResponse(t, listWalletsRes, &listWalletsResponseBody)

			if count := listWalletsResponseBody.Count; count != 1 {
				t.Errorf("Expected count: 1, got: %d", count)
			}
		}
	})

	t.Run("Creates, updates, and lists parties", func(t *testing.T) {
		{
			// Create party
			party := &handlers.Party{
				Name: "soviets",
			}

			createPartyReq := router_test.NewCreatePartyRequest(party, authToken)
			createPartyRes := httptest.NewRecorder()

			r.ServeHTTP(createPartyRes, createPartyReq)
			router_test.AssertStatusCode(t, createPartyRes, http.StatusCreated)

			var createPartyResponseBody handlers.Party
			router_test.ParsePartyJSONResponse(t, createPartyRes, &createPartyResponseBody)

			// Set party id
			partyID = createPartyResponseBody.ID

			// Get party
			getWalletReq := router_test.NewGetPartyRequest(partyID, authToken)
			getWalletRes := httptest.NewRecorder()

			r.ServeHTTP(getWalletRes, getWalletReq)
			router_test.AssertStatusCode(t, getWalletRes, http.StatusOK)

			var getPartyResponseBody handlers.Party
			router_test.ParsePartyJSONResponse(t, getWalletRes, &getPartyResponseBody)

			if party.Name != getPartyResponseBody.Name {
				t.Errorf("Expected wallet name: %s, got: %s", party.Name, getPartyResponseBody.Name)
			}
		}

		{
			// Update party
			updateParty := &handlers.Party{
				Name: "groceries",
			}

			updatePartyReq := router_test.NewUpdatePartyRequest(partyID, updateParty, authToken)
			updatePartyRes := httptest.NewRecorder()

			r.ServeHTTP(updatePartyRes, updatePartyReq)
			router_test.AssertStatusCode(t, updatePartyRes, http.StatusOK)

			var updatePartyResponseBody handlers.Party
			router_test.ParsePartyJSONResponse(t, updatePartyRes, &updatePartyResponseBody)

			if updateParty.Name != updatePartyResponseBody.Name {
				t.Errorf("Expected wallet name: %s, got: %s", updateParty.Name, updatePartyResponseBody.Name)
			}
		}

		{
			// List parties
			listPartiesReq := router_test.NewListPartiesRequest(authToken)
			listPartiesRes := httptest.NewRecorder()

			r.ServeHTTP(listPartiesRes, listPartiesReq)
			router_test.AssertStatusCode(t, listPartiesRes, http.StatusOK)

			var listPartiesResponseBody router_test.PartyListResponse
			router_test.ParsePartyListJSONResponse(t, listPartiesRes, &listPartiesResponseBody)

			if count := listPartiesResponseBody.Count; count != 1 {
				t.Errorf("Expected count: 1, got: %d", count)
			}
		}
	})

	t.Run("Creates, updates, and lists transactions", func(t *testing.T) {
		{
			// Create transaction
			transaction := &handlers.Transaction{
				Amount:   decimal.NewFromFloat(123.45),
				WalletID: walletID,
				PartyID:  partyID,
			}

			createTransactionReq := router_test.NewCreateTransactionRequest(transaction, authToken)
			createTransactionRes := httptest.NewRecorder()

			r.ServeHTTP(createTransactionRes, createTransactionReq)
			router_test.AssertStatusCode(t, createTransactionRes, http.StatusCreated)

			var createTransactionResponseBody handlers.Transaction
			router_test.ParseTransactionJSONResponse(t, createTransactionRes, &createTransactionResponseBody)

			// Set transaction id
			transactionID = createTransactionResponseBody.ID

			// Get transaction
			getTransactionReq := router_test.NewGetTransactionRequest(transactionID, authToken)
			getTransactionRes := httptest.NewRecorder()

			r.ServeHTTP(getTransactionRes, getTransactionReq)
			router_test.AssertStatusCode(t, getTransactionRes, http.StatusOK)

			var getTransactionResponseBody handlers.Transaction
			router_test.ParseTransactionJSONResponse(t, getTransactionRes, &getTransactionResponseBody)

			if transaction.Amount.Cmp(getTransactionResponseBody.Amount) != 0 {
				t.Errorf("Expected transaction amount: %v, got: %v", transaction.Amount, getTransactionResponseBody.Amount)
			}
		}

		{
			// Update transaction
			updateTransaction := &handlers.Transaction{
				Amount: decimal.NewFromFloat(99.99),
			}

			updateTransactionReq := router_test.NewUpdateTransactionRequest(transactionID, updateTransaction, authToken)
			updateTransactionRes := httptest.NewRecorder()

			r.ServeHTTP(updateTransactionRes, updateTransactionReq)
			router_test.AssertStatusCode(t, updateTransactionRes, http.StatusOK)

			var updateTransactionResponseBody handlers.Transaction
			router_test.ParseTransactionJSONResponse(t, updateTransactionRes, &updateTransactionResponseBody)

			if updateTransaction.Amount.Cmp(updateTransactionResponseBody.Amount) != 0 {
				t.Errorf("Expected transaction amount: %v, got: %v", updateTransaction.Amount, updateTransactionResponseBody.Amount)
			}
		}

		{
			// List transactions
			listTransactionsReq := router_test.NewListTransactionsRequest(authToken)
			listTransactionsRes := httptest.NewRecorder()

			r.ServeHTTP(listTransactionsRes, listTransactionsReq)
			router_test.AssertStatusCode(t, listTransactionsRes, http.StatusOK)

			var listTransactionsResponseBody router_test.TransactionListResponse
			router_test.ParseTransactionListJSONResponse(t, listTransactionsRes, &listTransactionsResponseBody)

			if count := listTransactionsResponseBody.Count; count != 1 {
				t.Errorf("Expected count: 1, got: %d", count)
			}
		}
	})

	t.Run("List transactions by wallet and party", func(t *testing.T) {
		{
			// List transactions by wallet
			listTransactionsByWalletReq := router_test.NewListTransactionsByWalletRequest(walletID, authToken)
			listTransactionsByWalletRes := httptest.NewRecorder()

			r.ServeHTTP(listTransactionsByWalletRes, listTransactionsByWalletReq)
			router_test.AssertStatusCode(t, listTransactionsByWalletRes, http.StatusOK)

			var transactions router_test.TransactionListResponse
			router_test.ParseTransactionListJSONResponse(t, listTransactionsByWalletRes, &transactions)

			if count := transactions.Count; count != 1 {
				t.Errorf("Expected count: 1, got: %d", count)
			}

			for _, transaction := range transactions.Entries {
				if transaction.WalletID != walletID {
					t.Errorf("Expected wallet id: %d, got: %d", walletID, transaction.WalletID)
				}
			}
		}

		{
			// List transactions by party
			listTransactionsByPartyReq := router_test.NewListTransactionsByPartyRequest(partyID, authToken)
			listTransactionsByPartyRes := httptest.NewRecorder()

			r.ServeHTTP(listTransactionsByPartyRes, listTransactionsByPartyReq)
			router_test.AssertStatusCode(t, listTransactionsByPartyRes, http.StatusOK)

			var transactions router_test.TransactionListResponse
			router_test.ParseTransactionListJSONResponse(t, listTransactionsByPartyRes, &transactions)

			if count := transactions.Count; count != 1 {
				t.Errorf("Expected count: 1, got: %d", count)
			}

			for _, transaction := range transactions.Entries {
				if transaction.PartyID != partyID {
					t.Errorf("Expected party id: %d, got: %d", partyID, transaction.PartyID)
				}
			}
		}
	})

	t.Run("Delete transaction, wallet, party and account", func(t *testing.T) {
		{
			// Delete transaction
			deleteTransactionReq := router_test.NewDeleteTransactionRequest(transactionID, authToken)
			deleteTransactionRes := httptest.NewRecorder()

			r.ServeHTTP(deleteTransactionRes, deleteTransactionReq)
			router_test.AssertStatusCode(t, deleteTransactionRes, http.StatusNoContent)

			// Get transaction
			getTransactionReq := router_test.NewGetTransactionRequest(transactionID, authToken)
			getTransactionRes := httptest.NewRecorder()

			r.ServeHTTP(getTransactionRes, getTransactionReq)
			router_test.AssertStatusCode(t, getTransactionRes, http.StatusNotFound)
		}

		{
			// Delete wallet
			deleteWalletReq := router_test.NewDeleteWalletRequest(walletID, authToken)
			deleteWalletRes := httptest.NewRecorder()

			r.ServeHTTP(deleteWalletRes, deleteWalletReq)
			router_test.AssertStatusCode(t, deleteWalletRes, http.StatusNoContent)

			// Get wallet
			getWalletReq := router_test.NewGetWalletRequest(walletID, authToken)
			getWalletRes := httptest.NewRecorder()

			r.ServeHTTP(getWalletRes, getWalletReq)
			router_test.AssertStatusCode(t, getWalletRes, http.StatusNotFound)
		}

		{
			// Delete party
			deletePartyReq := router_test.NewDeletePartyRequest(partyID, authToken)
			deletePartyRes := httptest.NewRecorder()

			r.ServeHTTP(deletePartyRes, deletePartyReq)
			router_test.AssertStatusCode(t, deletePartyRes, http.StatusNoContent)

			// Get party
			getWalletReq := router_test.NewGetPartyRequest(partyID, authToken)
			getWalletRes := httptest.NewRecorder()

			r.ServeHTTP(getWalletRes, getWalletReq)
			router_test.AssertStatusCode(t, getWalletRes, http.StatusNotFound)
		}

		{
			// Delete account
			deleteAccountReq := router_test.NewDeleteAccountRequest(authToken)
			deleteAccountRes := httptest.NewRecorder()

			r.ServeHTTP(deleteAccountRes, deleteAccountReq)
			router_test.AssertStatusCode(t, deleteAccountRes, http.StatusNoContent)

			// Login
			loginInfo := &handlers.LoginInfo{
				Email:    email,
				Password: password,
			}

			loginReq := router_test.NewLoginRequest(loginInfo)
			loginRes := httptest.NewRecorder()

			r.ServeHTTP(loginRes, loginReq)

			router_test.AssertStatusCode(t, loginRes, http.StatusNotFound)
		}
	})
}
