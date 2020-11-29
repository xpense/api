package router

import (
	"encoding/json"
	"expense-api/internal/handlers"
	"expense-api/internal/middleware/auth"
	"expense-api/internal/model"
	"expense-api/internal/repository"
	"expense-api/internal/router"
	"expense-api/test/spies"
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

	t.Run("Missing/Invalid authorization token cases", func(t *testing.T) {
		transaction := &handlers.Transaction{}
		token := "invalid-token"

		missingTokenReq := NewCreateTransactionRequest(transaction, token)
		invalidTokenReq := NewCreateTransactionRequest(transaction, token)

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
			transaction := &handlers.Transaction{
				Amount: decimal.Zero,
			}

			res := httptest.NewRecorder()
			req := NewCreateTransactionRequest(transaction, token)

			r.ServeHTTP(res, req)

			assertStatusCode(t, res, http.StatusBadRequest)
		})

		t.Run("Create transaction with valid data but missing wallet id", func(t *testing.T) {
			transaction := &handlers.Transaction{
				Timestamp: time.Date(2020, 12, 3, 19, 20, 0, 0, time.UTC),
				Amount:    decimal.NewFromInt32(100),
			}

			res := httptest.NewRecorder()
			req := NewCreateTransactionRequest(transaction, token)

			r.ServeHTTP(res, req)

			wantErrorMessage := handlers.ErrorRequiredWalletID.Error()

			assertStatusCode(t, res, http.StatusBadRequest)
			assertErrorMessage(t, res, wantErrorMessage)
		})

		t.Run("Create transaction with valid data with non-existent wallet id", func(t *testing.T) {
			walletID := uint(1)
			transaction := &handlers.Transaction{
				Timestamp: time.Date(2020, 12, 3, 19, 20, 0, 0, time.UTC),
				Amount:    decimal.NewFromInt32(100),
				WalletID:  walletID,
			}

			repoSpy.On("WalletGet", walletID).Return(nil, repository.ErrorRecordNotFound).Once()

			res := httptest.NewRecorder()
			req := NewCreateTransactionRequest(transaction, token)

			r.ServeHTTP(res, req)

			wantErrorMessage := handlers.ErrorWalletNotFound.Error()

			assertStatusCode(t, res, http.StatusBadRequest)
			assertErrorMessage(t, res, wantErrorMessage)
		})

		t.Run("Create transaction with valid data with existing wallet id that belongs to other user", func(t *testing.T) {
			walletID := uint(1)
			wallet := &model.Wallet{
				UserID: userID + 1,
			}
			transaction := &handlers.Transaction{
				Timestamp: time.Date(2020, 12, 3, 19, 20, 0, 0, time.UTC),
				Amount:    decimal.NewFromInt32(100),
				WalletID:  walletID,
			}

			repoSpy.On("WalletGet", walletID).Return(wallet, nil).Once()

			res := httptest.NewRecorder()
			req := NewCreateTransactionRequest(transaction, token)

			r.ServeHTTP(res, req)

			wantErrorMessage := handlers.ErrorBadWalletID.Error()

			assertStatusCode(t, res, http.StatusForbidden)
			assertErrorMessage(t, res, wantErrorMessage)
		})

		t.Run("Create transaction with valid data but missing party id", func(t *testing.T) {
			walletID := uint(1)
			wallet := &model.Wallet{
				UserID: userID,
			}
			transaction := &handlers.Transaction{
				Timestamp: time.Date(2020, 12, 3, 19, 20, 0, 0, time.UTC),
				Amount:    decimal.NewFromInt32(100),
				WalletID:  walletID,
			}

			repoSpy.On("WalletGet", walletID).Return(wallet, nil).Once()

			res := httptest.NewRecorder()
			req := NewCreateTransactionRequest(transaction, token)

			r.ServeHTTP(res, req)

			wantErrorMessage := handlers.ErrorRequiredWalletID.Error()

			assertStatusCode(t, res, http.StatusBadRequest)
			assertErrorMessage(t, res, wantErrorMessage)
		})

		t.Run("Create transaction with valid data with non-existent party id", func(t *testing.T) {
			walletID := uint(1)
			wallet := &model.Wallet{
				UserID: userID,
			}
			partyID := uint(2)
			transaction := &handlers.Transaction{
				Timestamp: time.Date(2020, 12, 3, 19, 20, 0, 0, time.UTC),
				Amount:    decimal.NewFromInt32(100),
				WalletID:  walletID,
				PartyID:   partyID,
			}

			repoSpy.On("WalletGet", walletID).Return(wallet, nil).Once()
			repoSpy.On("PartyGet", partyID).Return(nil, repository.ErrorRecordNotFound).Once()

			res := httptest.NewRecorder()
			req := NewCreateTransactionRequest(transaction, token)

			r.ServeHTTP(res, req)

			wantErrorMessage := handlers.ErrorPartyNotFound.Error()

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
			transaction := &handlers.Transaction{
				Timestamp: time.Date(2020, 12, 3, 19, 20, 0, 0, time.UTC),
				Amount:    decimal.NewFromInt32(100),
				WalletID:  walletID,
				PartyID:   partyID,
			}

			repoSpy.On("WalletGet", walletID).Return(wallet, nil).Once()
			repoSpy.On("PartyGet", partyID).Return(party, nil).Once()

			res := httptest.NewRecorder()
			req := NewCreateTransactionRequest(transaction, token)

			r.ServeHTTP(res, req)

			wantErrorMessage := handlers.ErrorBadPartyID.Error()

			assertStatusCode(t, res, http.StatusForbidden)
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
				Timestamp: time.Date(2020, 12, 3, 19, 20, 0, 0, time.UTC),
				Amount:    decimal.NewFromInt32(100),
				UserID:    userID,
				WalletID:  walletID,
				PartyID:   partyID,
			}

			repoSpy.On("WalletGet", walletID).Return(wallet, nil).Once()
			repoSpy.On("PartyGet", partyID).Return(party, nil).Once()
			repoSpy.On("TransactionCreate", transaction).Return(nil).Once()

			res := httptest.NewRecorder()
			req := NewCreateTransactionRequest(&handlers.Transaction{
				Timestamp: transaction.Timestamp,
				Amount:    transaction.Amount,
				WalletID:  transaction.WalletID,
				PartyID:   transaction.PartyID,
			}, token)

			r.ServeHTTP(res, req)

			resBody := handlers.TransactionModelToResponse(transaction)

			assertStatusCode(t, res, http.StatusCreated)
			assertSingleTransactionResponseBody(t, res, resBody)
		})
	})
}

func TestGetTransaction(t *testing.T) {
	repoSpy := &spies.RepositorySpy{}
	jwtServiceSpy := &spies.JWTServiceSpy{}
	hasherSpy := &spies.PasswordHasherSpy{}

	r := router.Setup(repoSpy, jwtServiceSpy, hasherSpy, router.TestConfig)

	t.Run("Missing/Invalid authorization token cases", func(t *testing.T) {
		id := uint(1)
		token := "invalid-token"

		missingTokenReq := NewGetTransactionRequest(id, token)
		invalidTokenReq := NewGetTransactionRequest(id, token)

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
			req := NewGetTransactionRequest(id, token)

			r.ServeHTTP(res, req)

			assertStatusCode(t, res, http.StatusBadRequest)
		})

		t.Run("Get transaction with non-existent id", func(t *testing.T) {
			id := uint(10)

			repoSpy.On("TransactionGet", id).Return(nil, repository.ErrorRecordNotFound).Once()

			res := httptest.NewRecorder()
			req := NewGetTransactionRequest(id, token)

			r.ServeHTTP(res, req)

			assertStatusCode(t, res, http.StatusNotFound)
		})

		t.Run("Get transaction with valid id that belongs to another user", func(t *testing.T) {
			id := uint(1)
			transaction := &model.Transaction{
				Timestamp: time.Date(2020, 12, 3, 19, 20, 0, 0, time.UTC),
				Amount:    decimal.NewFromInt32(100),
				UserID:    userID + 1,
			}

			repoSpy.On("TransactionGet", id).Return(transaction, nil).Once()

			res := httptest.NewRecorder()
			req := NewGetTransactionRequest(id, token)

			r.ServeHTTP(res, req)

			assertStatusCode(t, res, http.StatusForbidden)
		})

		t.Run("Get transaction with valid id", func(t *testing.T) {
			id := uint(1)
			transaction := &model.Transaction{
				Timestamp: time.Date(2020, 12, 3, 19, 20, 0, 0, time.UTC),
				Amount:    decimal.NewFromInt32(100),
				UserID:    userID,
			}

			repoSpy.On("TransactionGet", id).Return(transaction, nil).Twice()

			res := httptest.NewRecorder()
			req := NewGetTransactionRequest(id, token)

			r.ServeHTTP(res, req)

			resBody := handlers.TransactionModelToResponse(transaction)

			assertStatusCode(t, res, http.StatusOK)
			assertSingleTransactionResponseBody(t, res, resBody)
		})
	})
}

func TestUpdateTransaction(t *testing.T) {
	repoSpy := &spies.RepositorySpy{}
	jwtServiceSpy := &spies.JWTServiceSpy{}
	hasherSpy := &spies.PasswordHasherSpy{}

	r := router.Setup(repoSpy, jwtServiceSpy, hasherSpy, router.TestConfig)

	t.Run("Missing/Invalid authorization token cases", func(t *testing.T) {
		id := uint(1)
		transaction := &handlers.Transaction{}
		token := "invalid-token"

		missingTokenReq := NewUpdateTransactionRequest(id, transaction, token)
		invalidTokenReq := NewUpdateTransactionRequest(id, transaction, token)

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
			transaction := &handlers.Transaction{
				Amount: decimal.NewFromInt32(100),
			}

			repoSpy.On("TransactionGet", id).Return(nil, repository.ErrorRecordNotFound).Once()

			res := httptest.NewRecorder()
			req := NewUpdateTransactionRequest(id, transaction, token)

			r.ServeHTTP(res, req)

			assertStatusCode(t, res, http.StatusNotFound)
		})

		t.Run("Try to update transaction with valid id that belongs to another user", func(t *testing.T) {
			id := uint(1)
			transaction := &model.Transaction{
				Timestamp: time.Date(2020, 12, 3, 19, 20, 0, 0, time.UTC),
				Amount:    decimal.NewFromInt32(100),
				UserID:    userID + 1,
			}

			repoSpy.On("TransactionGet", id).Return(transaction, nil).Once()

			res := httptest.NewRecorder()
			req := NewUpdateTransactionRequest(id, &handlers.Transaction{
				Amount: transaction.Amount.Add(decimal.NewFromInt(100)),
			}, token)

			r.ServeHTTP(res, req)

			assertStatusCode(t, res, http.StatusForbidden)
		})

		t.Run("Change a transaction's wallet ID to one that doesn't exist", func(t *testing.T) {
			id := uint(1)
			walletID := uint(3)

			transaction := &model.Transaction{
				UserID: userID,
			}

			updateTransaction := &handlers.Transaction{
				WalletID: walletID,
			}

			repoSpy.On("TransactionGet", id).Return(transaction, nil).Once()
			repoSpy.On("WalletGet", walletID).Return(nil, repository.ErrorRecordNotFound).Once()

			res := httptest.NewRecorder()
			req := NewUpdateTransactionRequest(id, updateTransaction, token)

			r.ServeHTTP(res, req)

			wantErrorMessage := handlers.ErrorWalletNotFound.Error()

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

			updateTransaction := &handlers.Transaction{
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
			req := NewUpdateTransactionRequest(id, updateTransaction, token)

			r.ServeHTTP(res, req)

			wantErrorMessage := handlers.ErrorBadWalletID.Error()

			assertStatusCode(t, res, http.StatusForbidden)
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

			updateTransaction := &handlers.Transaction{
				PartyID: nonExistentPartyID,
			}

			repoSpy.On("TransactionGet", id).Return(transaction, nil).Once()
			repoSpy.On("WalletGet", walletID).Return(wallet, nil).Once()
			repoSpy.On("PartyGet", nonExistentPartyID).Return(nil, repository.ErrorRecordNotFound).Once()

			res := httptest.NewRecorder()
			req := NewUpdateTransactionRequest(id, updateTransaction, token)

			r.ServeHTTP(res, req)

			wantErrorMessage := handlers.ErrorPartyNotFound.Error()

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

			updateTransaction := &handlers.Transaction{
				PartyID: anotherUsersPartyID,
			}

			repoSpy.On("TransactionGet", id).Return(transaction, nil).Once()
			repoSpy.On("WalletGet", walletID).Return(wallet, nil).Once()
			repoSpy.On("PartyGet", anotherUsersPartyID).Return(anotherUsersParty, nil).Once()

			res := httptest.NewRecorder()
			req := NewUpdateTransactionRequest(id, updateTransaction, token)

			r.ServeHTTP(res, req)

			wantErrorMessage := handlers.ErrorBadPartyID.Error()

			assertStatusCode(t, res, http.StatusForbidden)
			assertErrorMessage(t, res, wantErrorMessage)
		})

		t.Run("Update existing transaction with valid arguments", func(t *testing.T) {
			id := uint(3)
			transaction := &model.Transaction{
				Amount: decimal.NewFromInt32(100),
				UserID: userID,
			}

			updateTransaction := &model.Transaction{
				Amount: transaction.Amount.Add(decimal.NewFromInt(100)),
				UserID: userID,
			}

			repoSpy.On("TransactionGet", id).Return(transaction, nil).Once()
			repoSpy.On("TransactionUpdate", id, updateTransaction).Return(updateTransaction, nil).Once()

			res := httptest.NewRecorder()
			req := NewUpdateTransactionRequest(id, &handlers.Transaction{
				Amount: updateTransaction.Amount,
			}, token)

			r.ServeHTTP(res, req)

			resBody := handlers.TransactionModelToResponse(updateTransaction)

			assertStatusCode(t, res, http.StatusOK)
			assertSingleTransactionResponseBody(t, res, resBody)
		})
	})
}

func TestDeleteTransaction(t *testing.T) {
	repoSpy := &spies.RepositorySpy{}
	jwtServiceSpy := &spies.JWTServiceSpy{}
	hasherSpy := &spies.PasswordHasherSpy{}

	r := router.Setup(repoSpy, jwtServiceSpy, hasherSpy, router.TestConfig)

	t.Run("Missing/Invalid authorization token cases", func(t *testing.T) {
		id := uint(1)
		token := "invalid-token"

		missingTokenReq := NewDeleteTransactionRequest(id, token)
		invalidTokenReq := NewDeleteTransactionRequest(id, token)

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
			req := NewDeleteTransactionRequest(id, token)

			r.ServeHTTP(res, req)

			assertStatusCode(t, res, http.StatusNotFound)
		})

		t.Run("Try to delete transaction with valid id that belongs to another user", func(t *testing.T) {
			id := uint(1)
			transaction := &model.Transaction{
				Timestamp: time.Date(2020, 12, 3, 19, 20, 0, 0, time.UTC),
				Amount:    decimal.NewFromInt32(100),
				UserID:    userID + 1,
			}

			repoSpy.On("TransactionGet", id).Return(transaction, nil).Once()

			res := httptest.NewRecorder()
			req := NewDeleteTransactionRequest(id, token)

			r.ServeHTTP(res, req)

			assertStatusCode(t, res, http.StatusForbidden)
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
			req := NewDeleteTransactionRequest(id, token)

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

	newTransactionListResponse := func(slice []*handlers.Transaction) *transactionListResponse {
		return &transactionListResponse{
			Count:   len(slice),
			Entries: slice,
		}
	}

	t.Run("Missing/Invalid authorization token cases", func(t *testing.T) {
		token := "invalid-token"

		missingTokenReq := NewListTransactionsRequest(token)
		invalidTokenReq := NewListTransactionsRequest(token)

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
			req := NewListTransactionsRequest(token)

			r.ServeHTTP(res, req)

			expected := newTransactionListResponse([]*handlers.Transaction{})

			assertStatusCode(t, res, http.StatusOK)
			assertListTransactionResponseBody(t, res, expected)
		})

		t.Run("List transactions when there are non-zero transactions", func(t *testing.T) {
			transactions := []*model.Transaction{{}}

			repoSpy.On("TransactionList", userID).Return(transactions, nil).Once()

			res := httptest.NewRecorder()
			req := NewListTransactionsRequest(token)

			r.ServeHTTP(res, req)

			expected := newTransactionListResponse([]*handlers.Transaction{{}})

			assertStatusCode(t, res, http.StatusOK)
			assertListTransactionResponseBody(t, res, expected)
		})
	})
}

type transactionListResponse struct {
	Count   int                     `json:"count"`
	Entries []*handlers.Transaction `json:"entries"`
}

func assertSingleTransactionResponseBody(t *testing.T, res *httptest.ResponseRecorder, transaction *handlers.Transaction) {
	t.Helper()

	var got handlers.Transaction
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
