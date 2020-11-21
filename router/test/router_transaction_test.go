package test

import (
	"bytes"
	"encoding/json"
	"expense-api/handlers"
	"expense-api/middleware/auth"
	"expense-api/model"
	"expense-api/repository"
	"expense-api/router"
	"expense-api/router/test/spies"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/shopspring/decimal"
)

func TestCreateTransaction(t *testing.T) {
	repoSpy := &spies.RepositorySpy{}
	jwtServiceSpy := &spies.JWTServiceSpy{}
	hasherSpy := &spies.PasswordHasherSpy{}

	r := router.Setup(repoSpy, jwtServiceSpy, hasherSpy, router.TestConfig)

	newTransactionRequest := func(transaction *model.Transaction, token string) *http.Request {
		body := createRequestBody(transaction)
		req, _ := http.NewRequest(http.MethodPost, "/transaction/", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		return req
	}

	t.Run("Missing/Invalid authorization token cases", func(t *testing.T) {
		transaction := &model.Transaction{}
		token := "invalid-token"

		missingTokenReq := newTransactionRequest(transaction, token)
		invalidTokenReq := newTransactionRequest(transaction, token)

		unauthorizedTestCases := UnauthorizedTestCases(missingTokenReq, invalidTokenReq, r, jwtServiceSpy)
		t.Run("Unauthorized test cases", unauthorizedTestCases)
	})

	t.Run("Valid authorization token cases", func(t *testing.T) {
		token := "valid-token"
		userID := uint(1)
		claims := auth.CustomClaims{
			ID: userID,
		}
		jwtServiceSpy.On("ValidateJWT", token).Return(&claims, nil)

		t.Run("Create transaction with amount = 0", func(t *testing.T) {
			transaction := &model.Transaction{
				Amount: decimal.Zero,
			}

			res := httptest.NewRecorder()
			req := newTransactionRequest(transaction, token)

			r.ServeHTTP(res, req)

			assertStatusCode(t, res, http.StatusBadRequest)
		})

		t.Run("Create transaction with valid data but missing wallet id", func(t *testing.T) {
			transaction := &model.Transaction{
				Timestamp: time.Now().Round(0),
				Amount:    decimal.NewFromInt32(100),
				UserID:    userID,
			}

			res := httptest.NewRecorder()
			req := newTransactionRequest(transaction, token)

			r.ServeHTTP(res, req)

			wantErrorMessage := handlers.ErrMsgRequiredWalletID

			assertStatusCode(t, res, http.StatusBadRequest)
			assertErrorMessage(t, res, wantErrorMessage)
		})

		t.Run("Create transaction with valid data with non-existent wallet id", func(t *testing.T) {
			walletID := uint(1)
			transaction := &model.Transaction{
				Timestamp: time.Now().Round(0),
				Amount:    decimal.NewFromInt32(100),
				UserID:    userID,
				WalletID:  walletID,
			}

			repoSpy.On("WalletGet", walletID).Return(nil, repository.ErrorRecordNotFound).Once()

			res := httptest.NewRecorder()
			req := newTransactionRequest(transaction, token)

			r.ServeHTTP(res, req)

			wantErrorMessage := handlers.ErrMsgWalletNotFound

			assertStatusCode(t, res, http.StatusBadRequest)
			assertErrorMessage(t, res, wantErrorMessage)
		})

		t.Run("Create transaction with valid data with existing wallet id that belongs to other user", func(t *testing.T) {
			walletID := uint(1)
			wallet := &model.Wallet{
				UserID: userID + 1,
			}
			transaction := &model.Transaction{
				Timestamp: time.Now().Round(0),
				Amount:    decimal.NewFromInt32(100),
				UserID:    userID,
				WalletID:  walletID,
			}

			repoSpy.On("WalletGet", walletID).Return(wallet, nil).Once()

			res := httptest.NewRecorder()
			req := newTransactionRequest(transaction, token)

			r.ServeHTTP(res, req)

			wantErrorMessage := handlers.ErrMsgBadWalletID

			assertStatusCode(t, res, http.StatusUnauthorized)
			assertErrorMessage(t, res, wantErrorMessage)
		})

		t.Run("Create transaction with valid data but missing party id", func(t *testing.T) {
			walletID := uint(1)
			wallet := &model.Wallet{
				UserID: userID,
			}
			transaction := &model.Transaction{
				Timestamp: time.Now().Round(0),
				Amount:    decimal.NewFromInt32(100),
				UserID:    userID,
				WalletID:  walletID,
			}

			repoSpy.On("WalletGet", walletID).Return(wallet, nil).Once()

			res := httptest.NewRecorder()
			req := newTransactionRequest(transaction, token)

			r.ServeHTTP(res, req)

			wantErrorMessage := handlers.ErrMsgRequiredWalletID

			assertStatusCode(t, res, http.StatusBadRequest)
			assertErrorMessage(t, res, wantErrorMessage)
		})

		t.Run("Create transaction with valid data with non-existent party id", func(t *testing.T) {
			walletID := uint(1)
			wallet := &model.Wallet{
				UserID: userID,
			}
			partyID := uint(2)
			transaction := &model.Transaction{
				Timestamp: time.Now().Round(0),
				Amount:    decimal.NewFromInt32(100),
				UserID:    userID,
				WalletID:  walletID,
				PartyID:   partyID,
			}

			repoSpy.On("WalletGet", walletID).Return(wallet, nil).Once()
			repoSpy.On("PartyGet", partyID).Return(nil, repository.ErrorRecordNotFound).Once()

			res := httptest.NewRecorder()
			req := newTransactionRequest(transaction, token)

			r.ServeHTTP(res, req)

			wantErrorMessage := handlers.ErrMsgPartyNotFound

			assertStatusCode(t, res, http.StatusBadRequest)
			assertErrorMessage(t, res, wantErrorMessage)
		})

		t.Run("Create transaction with valid data with existing party id that belongs to another user", func(t *testing.T) {
			walletID := uint(1)
			wallet := &model.Wallet{
				UserID: userID,
			}
			partyID := uint(1)
			party := &model.Party{
				UserID: userID + 1,
			}
			transaction := &model.Transaction{
				Timestamp: time.Now().Round(0),
				Amount:    decimal.NewFromInt32(100),
				UserID:    userID,
				WalletID:  walletID,
				PartyID:   partyID,
			}

			repoSpy.On("WalletGet", walletID).Return(wallet, nil).Once()
			repoSpy.On("PartyGet", partyID).Return(party, nil).Once()

			res := httptest.NewRecorder()
			req := newTransactionRequest(transaction, token)

			r.ServeHTTP(res, req)

			wantErrorMessage := handlers.ErrMsgBadPartyID

			assertStatusCode(t, res, http.StatusUnauthorized)
			assertErrorMessage(t, res, wantErrorMessage)
		})

		t.Run("Create transaction with valid data", func(t *testing.T) {
			walletID := uint(1)
			wallet := &model.Wallet{
				UserID: userID,
			}
			partyID := uint(1)
			party := &model.Party{
				UserID: userID,
			}
			transaction := &model.Transaction{
				Timestamp: time.Now().Round(0),
				Amount:    decimal.NewFromInt32(100),
				UserID:    userID,
				WalletID:  walletID,
				PartyID:   partyID,
			}

			repoSpy.On("WalletGet", walletID).Return(wallet, nil).Once()
			repoSpy.On("PartyGet", partyID).Return(party, nil).Once()
			repoSpy.On("TransactionCreate", transaction).Return(nil).Once()

			res := httptest.NewRecorder()
			req := newTransactionRequest(transaction, token)

			r.ServeHTTP(res, req)

			transaction.UserID = 0

			assertStatusCode(t, res, http.StatusCreated)
			assertSingleTransactionResponseBody(t, res, transaction)
		})
	})
}

func TestGetTransaction(t *testing.T) {
	repoSpy := &spies.RepositorySpy{}
	jwtServiceSpy := &spies.JWTServiceSpy{}
	hasherSpy := &spies.PasswordHasherSpy{}

	r := router.Setup(repoSpy, jwtServiceSpy, hasherSpy, router.TestConfig)

	newTransactionRequest := func(id uint, token string) *http.Request {
		url := fmt.Sprintf("/transaction/%d", id)
		req, _ := http.NewRequest(http.MethodGet, url, nil)
		req.Header.Set("Authorization", "Bearer "+token)
		return req
	}

	t.Run("Missing/Invalid authorization token cases", func(t *testing.T) {
		id := uint(1)
		token := "invalid-token"

		missingTokenReq := newTransactionRequest(id, token)
		invalidTokenReq := newTransactionRequest(id, token)

		unauthorizedTestCases := UnauthorizedTestCases(missingTokenReq, invalidTokenReq, r, jwtServiceSpy)
		t.Run("Unauthorized test cases", unauthorizedTestCases)
	})

	t.Run("Valid authorization token cases", func(t *testing.T) {
		token := "valid-token"
		userID := uint(1)
		claims := auth.CustomClaims{
			ID: userID,
		}
		jwtServiceSpy.On("ValidateJWT", token).Return(&claims, nil)

		t.Run("Get transaction with id = 0", func(t *testing.T) {
			id := uint(0)

			res := httptest.NewRecorder()
			req := newTransactionRequest(id, token)

			r.ServeHTTP(res, req)

			assertStatusCode(t, res, http.StatusBadRequest)
		})

		t.Run("Get transaction with non-existent id", func(t *testing.T) {
			id := uint(10)

			repoSpy.On("TransactionGet", id).Return(nil, repository.ErrorRecordNotFound).Once()

			res := httptest.NewRecorder()
			req := newTransactionRequest(id, token)

			r.ServeHTTP(res, req)

			assertStatusCode(t, res, http.StatusNotFound)
		})

		t.Run("Get transaction with valid id that belongs to another user", func(t *testing.T) {
			id := uint(1)
			transaction := &model.Transaction{
				Timestamp: time.Now().Round(0),
				Amount:    decimal.NewFromInt32(100),
				UserID:    userID + 1,
			}

			repoSpy.On("TransactionGet", id).Return(transaction, nil).Once()

			res := httptest.NewRecorder()
			req := newTransactionRequest(id, token)

			r.ServeHTTP(res, req)

			assertStatusCode(t, res, http.StatusUnauthorized)
		})

		t.Run("Get transaction with valid id", func(t *testing.T) {
			id := uint(1)
			transaction := &model.Transaction{
				Timestamp: time.Now().Round(0),
				Amount:    decimal.NewFromInt32(100),
				UserID:    userID,
			}

			repoSpy.On("TransactionGet", id).Return(transaction, nil).Twice()

			res := httptest.NewRecorder()
			req := newTransactionRequest(id, token)

			r.ServeHTTP(res, req)

			transaction.UserID = 0

			assertStatusCode(t, res, http.StatusOK)
			assertSingleTransactionResponseBody(t, res, transaction)
		})
	})
}

func TestUpdateTransaction(t *testing.T) {
	repoSpy := &spies.RepositorySpy{}
	jwtServiceSpy := &spies.JWTServiceSpy{}
	hasherSpy := &spies.PasswordHasherSpy{}

	r := router.Setup(repoSpy, jwtServiceSpy, hasherSpy, router.TestConfig)

	newTransactionRequest := func(id uint, transaction *model.Transaction, token string) *http.Request {
		url := fmt.Sprintf("/transaction/%d", id)
		body := createRequestBody(transaction)
		req, _ := http.NewRequest(http.MethodPatch, url, bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		return req
	}

	t.Run("Missing/Invalid authorization token cases", func(t *testing.T) {
		id := uint(1)
		transaction := &model.Transaction{}
		token := "invalid-token"

		missingTokenReq := newTransactionRequest(id, transaction, token)
		invalidTokenReq := newTransactionRequest(id, transaction, token)

		unauthorizedTestCases := UnauthorizedTestCases(missingTokenReq, invalidTokenReq, r, jwtServiceSpy)
		t.Run("Unauthorized test cases", unauthorizedTestCases)
	})

	t.Run("Valid authorization token cases", func(t *testing.T) {
		token := "valid-token"
		userID := uint(1)
		claims := auth.CustomClaims{
			ID: userID,
		}
		jwtServiceSpy.On("ValidateJWT", token).Return(&claims, nil)

		t.Run("Update non-existent transaction", func(t *testing.T) {
			id := uint(1)
			transaction := &model.Transaction{
				Amount: decimal.NewFromInt32(100),
			}

			repoSpy.On("TransactionGet", id).Return(nil, repository.ErrorRecordNotFound).Once()

			res := httptest.NewRecorder()
			req := newTransactionRequest(id, transaction, token)

			r.ServeHTTP(res, req)

			assertStatusCode(t, res, http.StatusNotFound)
		})

		t.Run("Try to update transaction with valid id that belongs to another user", func(t *testing.T) {
			id := uint(1)
			transaction := &model.Transaction{
				Timestamp: time.Now().Round(0),
				Amount:    decimal.NewFromInt32(100),
				UserID:    userID + 1,
			}

			repoSpy.On("TransactionGet", id).Return(transaction, nil).Once()

			res := httptest.NewRecorder()
			req := newTransactionRequest(id, transaction, token)

			r.ServeHTTP(res, req)

			assertStatusCode(t, res, http.StatusUnauthorized)
		})

		t.Run("Change a transaction's wallet ID to one that doesn't exist", func(t *testing.T) {
			id := uint(1)
			walletID := uint(3)

			transaction := &model.Transaction{
				UserID: userID,
			}

			updateTransaction := &model.Transaction{
				WalletID: walletID,
			}

			repoSpy.On("TransactionGet", id).Return(transaction, nil).Once()
			repoSpy.On("WalletGet", walletID).Return(nil, repository.ErrorRecordNotFound).Once()

			res := httptest.NewRecorder()
			req := newTransactionRequest(id, updateTransaction, token)

			r.ServeHTTP(res, req)

			wantErrorMessage := handlers.ErrMsgWalletNotFound

			assertStatusCode(t, res, http.StatusBadRequest)
			assertErrorMessage(t, res, wantErrorMessage)
		})

		t.Run("Change a transaction's wallet ID to one that belongs to another user", func(t *testing.T) {
			id := uint(1)
			anotherUsersID := uint(199999)
			anotherUsersWalletID := uint(3)

			transaction := &model.Transaction{
				UserID: userID,
			}

			updateTransaction := &model.Transaction{
				WalletID: anotherUsersWalletID,
			}

			anotherUsersWallet := &model.Wallet{
				UserID: anotherUsersID,
			}
			anotherUsersWallet.ID = anotherUsersWalletID

			repoSpy.On("TransactionGet", id).Return(transaction, nil).Once()
			repoSpy.On("WalletGet", anotherUsersWalletID).Return(anotherUsersWallet, nil).Once()
			repoSpy.On("TransactionUpdate", id, updateTransaction).Return(nil, repository.ErrorUniqueConstaintViolation).Once()

			res := httptest.NewRecorder()
			req := newTransactionRequest(id, updateTransaction, token)

			r.ServeHTTP(res, req)

			wantErrorMessage := handlers.ErrMsgBadWalletID

			assertStatusCode(t, res, http.StatusUnauthorized)
			assertErrorMessage(t, res, wantErrorMessage)
		})

		t.Run("Change a transaction's party ID to one that doesn't exist", func(t *testing.T) {
			id := uint(1)
			nonExistentPartyID := uint(3)

			walletID := uint(1)
			wallet := &model.Wallet{
				UserID: userID,
			}

			transaction := &model.Transaction{
				UserID: userID,
			}

			updateTransaction := &model.Transaction{
				PartyID: nonExistentPartyID,
				UserID:  userID,
			}

			repoSpy.On("TransactionGet", id).Return(transaction, nil).Once()
			repoSpy.On("WalletGet", walletID).Return(wallet, nil).Once()
			repoSpy.On("PartyGet", nonExistentPartyID).Return(nil, repository.ErrorRecordNotFound).Once()

			res := httptest.NewRecorder()
			req := newTransactionRequest(id, updateTransaction, token)

			r.ServeHTTP(res, req)

			wantErrorMessage := handlers.ErrMsgPartyNotFound

			assertStatusCode(t, res, http.StatusBadRequest)
			assertErrorMessage(t, res, wantErrorMessage)
		})

		t.Run("Change a transaction's party ID to one that belongs to another user", func(t *testing.T) {
			id := uint(1)

			walletID := uint(1)
			wallet := &model.Wallet{
				UserID: userID,
			}

			anotherUsersID := uint(6)
			anotherUsersPartyID := uint(5)
			anotherUsersParty := &model.Party{
				UserID: anotherUsersID,
			}

			transaction := &model.Transaction{
				UserID: userID,
			}

			updateTransaction := &model.Transaction{
				PartyID: anotherUsersPartyID,
			}

			repoSpy.On("TransactionGet", id).Return(transaction, nil).Once()
			repoSpy.On("WalletGet", walletID).Return(wallet, nil).Once()
			repoSpy.On("PartyGet", anotherUsersPartyID).Return(anotherUsersParty, nil).Once()

			res := httptest.NewRecorder()
			req := newTransactionRequest(id, updateTransaction, token)

			r.ServeHTTP(res, req)

			wantErrorMessage := handlers.ErrMsgBadPartyID

			assertStatusCode(t, res, http.StatusUnauthorized)
			assertErrorMessage(t, res, wantErrorMessage)
		})

		t.Run("Update existing transaction with valid arguments", func(t *testing.T) {
			id := uint(3)
			transaction := &model.Transaction{
				Amount: decimal.NewFromInt32(100),
				UserID: userID,
			}

			repoSpy.On("TransactionGet", id).Return(transaction, nil).Once()
			repoSpy.On("TransactionUpdate", id, transaction).Return(transaction, nil).Once()

			res := httptest.NewRecorder()
			req := newTransactionRequest(id, transaction, token)

			r.ServeHTTP(res, req)

			transaction.UserID = 0

			assertStatusCode(t, res, http.StatusOK)
			assertSingleTransactionResponseBody(t, res, transaction)
		})
	})
}

func TestDeleteTransaction(t *testing.T) {
	repoSpy := &spies.RepositorySpy{}
	jwtServiceSpy := &spies.JWTServiceSpy{}
	hasherSpy := &spies.PasswordHasherSpy{}

	r := router.Setup(repoSpy, jwtServiceSpy, hasherSpy, router.TestConfig)

	newTransactionRequest := func(id uint, token string) *http.Request {
		url := fmt.Sprintf("/transaction/%d", id)
		req, _ := http.NewRequest(http.MethodDelete, url, nil)
		req.Header.Set("Authorization", "Bearer "+token)
		return req
	}

	t.Run("Missing/Invalid authorization token cases", func(t *testing.T) {
		id := uint(1)
		token := "invalid-token"

		missingTokenReq := newTransactionRequest(id, token)
		invalidTokenReq := newTransactionRequest(id, token)

		unauthorizedTestCases := UnauthorizedTestCases(missingTokenReq, invalidTokenReq, r, jwtServiceSpy)
		t.Run("Unauthorized test cases", unauthorizedTestCases)
	})

	t.Run("Valid authorization token cases", func(t *testing.T) {
		token := "valid-token"
		userID := uint(1)
		claims := auth.CustomClaims{
			ID: userID,
		}
		jwtServiceSpy.On("ValidateJWT", token).Return(&claims, nil)

		t.Run("Delete non-existent transaction", func(t *testing.T) {
			id := uint(1)

			repoSpy.On("TransactionGet", id).Return(nil, repository.ErrorRecordNotFound).Once()

			res := httptest.NewRecorder()
			req := newTransactionRequest(id, token)

			r.ServeHTTP(res, req)

			assertStatusCode(t, res, http.StatusNotFound)
		})

		t.Run("Try to delete transaction with valid id that belongs to another user", func(t *testing.T) {
			id := uint(1)
			transaction := &model.Transaction{
				Timestamp: time.Now().Round(0),
				Amount:    decimal.NewFromInt32(100),
				UserID:    userID + 1,
			}

			repoSpy.On("TransactionGet", id).Return(transaction, nil).Once()

			res := httptest.NewRecorder()
			req := newTransactionRequest(id, token)

			r.ServeHTTP(res, req)

			assertStatusCode(t, res, http.StatusUnauthorized)
		})

		t.Run("Delete existing transaction", func(t *testing.T) {
			id := uint(2)
			transaction := &model.Transaction{
				Amount: decimal.NewFromInt32(100),
				UserID: userID,
			}

			repoSpy.On("TransactionGet", id).Return(transaction, nil).Once()
			repoSpy.On("TransactionDelete", id).Return(nil).Once()

			res := httptest.NewRecorder()
			req := newTransactionRequest(id, token)

			r.ServeHTTP(res, req)

			assertStatusCode(t, res, http.StatusNoContent)
		})
	})
}

func TestListTransactions(t *testing.T) {
	repoSpy := &spies.RepositorySpy{}
	jwtServiceSpy := &spies.JWTServiceSpy{}
	hasherSpy := &spies.PasswordHasherSpy{}

	r := router.Setup(repoSpy, jwtServiceSpy, hasherSpy, router.TestConfig)

	newTransactionListResponse := func(slice []*model.Transaction) *transactionListResponse {
		return &transactionListResponse{
			Count:   len(slice),
			Entries: slice,
		}
	}

	newTransactionRequest := func(token string) *http.Request {
		req, _ := http.NewRequest(http.MethodGet, "/transaction/", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		return req
	}

	t.Run("Missing/Invalid authorization token cases", func(t *testing.T) {
		token := "invalid-token"

		missingTokenReq := newTransactionRequest(token)
		invalidTokenReq := newTransactionRequest(token)

		unauthorizedTestCases := UnauthorizedTestCases(missingTokenReq, invalidTokenReq, r, jwtServiceSpy)
		t.Run("Unauthorized test cases", unauthorizedTestCases)
	})

	t.Run("Valid authorization token cases", func(t *testing.T) {
		token := "valid-token"
		userID := uint(1)
		claims := auth.CustomClaims{
			ID: userID,
		}
		jwtServiceSpy.On("ValidateJWT", token).Return(&claims, nil)

		t.Run("List transactions when there are no transactions", func(t *testing.T) {
			transactions := []*model.Transaction{}
			repoSpy.On("TransactionList", userID).Return(transactions, nil).Once()

			res := httptest.NewRecorder()
			req := newTransactionRequest(token)

			r.ServeHTTP(res, req)

			expected := newTransactionListResponse(transactions)

			assertStatusCode(t, res, http.StatusOK)
			assertListTransactionResponseBody(t, res, expected)
		})

		t.Run("List transactions when there are non-zero transactions", func(t *testing.T) {
			transactions := []*model.Transaction{{}}

			repoSpy.On("TransactionList", userID).Return(transactions, nil).Once()

			res := httptest.NewRecorder()
			req := newTransactionRequest(token)

			r.ServeHTTP(res, req)

			expected := newTransactionListResponse(transactions)

			assertStatusCode(t, res, http.StatusOK)
			assertListTransactionResponseBody(t, res, expected)
		})
	})
}

type transactionListResponse struct {
	Count   int                  `json:"count"`
	Entries []*model.Transaction `json:"entries"`
}

func assertSingleTransactionResponseBody(t *testing.T, res *httptest.ResponseRecorder, transaction *model.Transaction) {
	t.Helper()

	var got model.Transaction
	if err := json.NewDecoder(res.Body).Decode(&got); err != nil {
		t.Errorf("couldn't parse json response: %v", err)
	}

	if !cmp.Equal(got, *transaction) {
		t.Errorf("expected %+v, got %+v", *transaction, got)
	}
}

func assertListTransactionResponseBody(t *testing.T, res *httptest.ResponseRecorder, expected *transactionListResponse) {
	t.Helper()

	var got transactionListResponse
	if err := json.NewDecoder(res.Body).Decode(&got); err != nil {
		t.Errorf("couldn't parse json response: %v", err)
	}

	if !cmp.Equal(got, *expected) {
		t.Errorf("expected %+v ;%T, got %+v ;%T", *expected, *expected, got, got)
	}
}
