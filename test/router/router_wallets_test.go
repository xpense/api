package router

import (
	"expense-api/internal/handlers"
	"expense-api/internal/middleware/auth"
	"expense-api/internal/model"
	"expense-api/internal/repository"
	"expense-api/internal/router"
	"expense-api/test/spies"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateWallet(t *testing.T) {
	repoSpy := &spies.RepositorySpy{}
	jwtServiceSpy := &spies.JWTServiceSpy{}
	hasherSpy := &spies.PasswordHasherSpy{}

	r := router.Setup(repoSpy, jwtServiceSpy, hasherSpy, router.TestConfig)

	t.Run("Missing/Invalid authorization token cases", func(t *testing.T) {
		wallet := &handlers.Wallet{}
		token := "invalid-token"

		missingTokenReq := NewCreateWalletRequest(wallet, token)
		invalidTokenReq := NewCreateWalletRequest(wallet, token)

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

		t.Run("Try to create a wallet with already existing name, belonging to the same user", func(t *testing.T) {
			wallet := &model.Wallet{
				Name:   "cash",
				UserID: userID,
			}

			repoSpy.On("WalletCreate", wallet).Return(repository.ErrorUniqueConstaintViolation).Once()

			res := httptest.NewRecorder()
			req := NewCreateWalletRequest(&handlers.Wallet{
				Name: wallet.Name,
			}, token)

			r.ServeHTTP(res, req)

			wantErrorMessage := handlers.ErrorWalletNameTaken.Error()

			AssertStatusCode(t, res, http.StatusConflict)
			AssertErrorMessage(t, res, wantErrorMessage)
		})

		t.Run("Create wallet with valid data", func(t *testing.T) {
			wallet := &model.Wallet{
				Name:   "cash",
				UserID: userID,
			}

			repoSpy.On("WalletCreate", wallet).Return(nil).Once()

			res := httptest.NewRecorder()
			req := NewCreateWalletRequest(&handlers.Wallet{
				Name: wallet.Name,
			}, token)

			r.ServeHTTP(res, req)

			resBody := handlers.WalletModelToResponse(wallet)

			AssertStatusCode(t, res, http.StatusCreated)
			AssertResponseBody(t, res, resBody)
		})
	})
}

func TestGetWallet(t *testing.T) {
	repoSpy := &spies.RepositorySpy{}
	jwtServiceSpy := &spies.JWTServiceSpy{}
	hasherSpy := &spies.PasswordHasherSpy{}

	r := router.Setup(repoSpy, jwtServiceSpy, hasherSpy, router.TestConfig)

	t.Run("Missing/Invalid authorization token cases", func(t *testing.T) {
		id := uint(1)
		token := "invalid-token"

		missingTokenReq := NewGetWalletRequest(id, token)
		invalidTokenReq := NewGetWalletRequest(id, token)

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

		t.Run("Get wallet with id = 0", func(t *testing.T) {
			id := uint(0)

			res := httptest.NewRecorder()
			req := NewGetWalletRequest(id, token)

			r.ServeHTTP(res, req)

			AssertStatusCode(t, res, http.StatusBadRequest)
		})

		t.Run("Get wallet with non-existent id", func(t *testing.T) {
			id := uint(10)

			repoSpy.On("WalletGet", id).Return(nil, repository.ErrorRecordNotFound).Once()

			res := httptest.NewRecorder()
			req := NewGetWalletRequest(id, token)

			r.ServeHTTP(res, req)

			AssertStatusCode(t, res, http.StatusNotFound)
		})

		t.Run("Get wallet with valid id that belongs to another user", func(t *testing.T) {
			id := uint(1)
			wallet := &model.Wallet{
				Name:   "new wallet",
				UserID: userID + 1,
			}

			repoSpy.On("WalletGet", id).Return(wallet, nil).Once()

			res := httptest.NewRecorder()
			req := NewGetWalletRequest(id, token)

			r.ServeHTTP(res, req)

			AssertStatusCode(t, res, http.StatusForbidden)
		})

		t.Run("Get wallet with valid id", func(t *testing.T) {
			id := uint(1)
			wallet := &model.Wallet{
				Name:   "new wallet",
				UserID: userID,
			}

			repoSpy.On("WalletGet", id).Return(wallet, nil).Twice()

			res := httptest.NewRecorder()
			req := NewGetWalletRequest(id, token)

			r.ServeHTTP(res, req)

			resBody := handlers.WalletModelToResponse(wallet)

			AssertStatusCode(t, res, http.StatusOK)
			AssertResponseBody(t, res, resBody)
		})
	})
}

func TestUpdateWallet(t *testing.T) {
	repoSpy := &spies.RepositorySpy{}
	jwtServiceSpy := &spies.JWTServiceSpy{}
	hasherSpy := &spies.PasswordHasherSpy{}

	r := router.Setup(repoSpy, jwtServiceSpy, hasherSpy, router.TestConfig)

	t.Run("Missing/Invalid authorization token cases", func(t *testing.T) {
		id := uint(1)
		wallet := &handlers.Wallet{}
		token := "invalid-token"

		missingTokenReq := NewUpdateWalletRequest(id, wallet, token)
		invalidTokenReq := NewUpdateWalletRequest(id, wallet, token)

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

		t.Run("Update non-existent wallet", func(t *testing.T) {
			id := uint(1)
			wallet := &handlers.Wallet{
				Name: "new wallet",
			}

			repoSpy.On("WalletGet", id).Return(nil, repository.ErrorRecordNotFound).Once()

			res := httptest.NewRecorder()
			req := NewUpdateWalletRequest(id, wallet, token)

			r.ServeHTTP(res, req)

			AssertStatusCode(t, res, http.StatusNotFound)
		})

		t.Run("Try to update wallet with valid id that belongs to another user", func(t *testing.T) {
			id := uint(1)
			wallet := &model.Wallet{
				Name:   "new wallet",
				UserID: userID + 1,
			}

			repoSpy.On("WalletGet", id).Return(wallet, nil).Once()

			res := httptest.NewRecorder()
			req := NewUpdateWalletRequest(id, &handlers.Wallet{Name: wallet.Name}, token)

			r.ServeHTTP(res, req)

			AssertStatusCode(t, res, http.StatusForbidden)
		})

		t.Run("Try to update a wallet with already existing name, belonging to the same user", func(t *testing.T) {
			id := uint(1)
			wallet := &model.Wallet{
				Name:   "cash",
				UserID: userID,
			}

			repoSpy.On("WalletGet", id).Return(wallet, nil).Once()
			repoSpy.On("WalletUpdate", id, wallet).Return(nil, repository.ErrorUniqueConstaintViolation).Once()

			res := httptest.NewRecorder()
			req := NewUpdateWalletRequest(id, &handlers.Wallet{Name: wallet.Name}, token)

			r.ServeHTTP(res, req)

			wantErrorMessage := handlers.ErrorWalletNameTaken.Error()

			AssertStatusCode(t, res, http.StatusConflict)
			AssertErrorMessage(t, res, wantErrorMessage)
		})

		t.Run("Update existing wallet with valid arguments", func(t *testing.T) {
			id := uint(3)
			wallet := &model.Wallet{
				Name:   "new wallet",
				UserID: userID,
			}

			repoSpy.On("WalletGet", id).Return(wallet, nil).Once()
			repoSpy.On("WalletUpdate", id, wallet).Return(wallet, nil).Once()

			res := httptest.NewRecorder()
			req := NewUpdateWalletRequest(id, &handlers.Wallet{Name: wallet.Name}, token)

			r.ServeHTTP(res, req)

			resBody := handlers.WalletModelToResponse(wallet)

			AssertStatusCode(t, res, http.StatusOK)
			AssertResponseBody(t, res, resBody)
		})
	})
}

func TestDeleteWallet(t *testing.T) {
	repoSpy := &spies.RepositorySpy{}
	jwtServiceSpy := &spies.JWTServiceSpy{}
	hasherSpy := &spies.PasswordHasherSpy{}

	r := router.Setup(repoSpy, jwtServiceSpy, hasherSpy, router.TestConfig)

	t.Run("Missing/Invalid authorization token cases", func(t *testing.T) {
		id := uint(1)
		token := "invalid-token"

		missingTokenReq := NewDeleteWalletRequest(id, token)
		invalidTokenReq := NewDeleteWalletRequest(id, token)

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

		t.Run("Delete non-existent wallet", func(t *testing.T) {
			id := uint(1)

			repoSpy.On("WalletGet", id).Return(nil, repository.ErrorRecordNotFound).Once()

			res := httptest.NewRecorder()
			req := NewDeleteWalletRequest(id, token)

			r.ServeHTTP(res, req)

			AssertStatusCode(t, res, http.StatusNotFound)
		})

		t.Run("Try to delete wallet with valid id that belongs to another user", func(t *testing.T) {
			id := uint(1)
			wallet := &model.Wallet{
				Name:   "new wallet",
				UserID: userID + 1,
			}

			repoSpy.On("WalletGet", id).Return(wallet, nil).Once()

			res := httptest.NewRecorder()
			req := NewDeleteWalletRequest(id, token)

			r.ServeHTTP(res, req)

			AssertStatusCode(t, res, http.StatusForbidden)
		})

		t.Run("Delete existing wallet", func(t *testing.T) {
			id := uint(2)
			wallet := &model.Wallet{
				Name:   "new wallet",
				UserID: userID,
			}

			repoSpy.On("WalletGet", id).Return(wallet, nil).Once()
			repoSpy.On("WalletDelete", id).Return(nil).Once()

			res := httptest.NewRecorder()
			req := NewDeleteWalletRequest(id, token)

			r.ServeHTTP(res, req)

			AssertStatusCode(t, res, http.StatusNoContent)
		})
	})
}

func TestListWallets(t *testing.T) {
	repoSpy := &spies.RepositorySpy{}
	jwtServiceSpy := &spies.JWTServiceSpy{}
	hasherSpy := &spies.PasswordHasherSpy{}

	r := router.Setup(repoSpy, jwtServiceSpy, hasherSpy, router.TestConfig)

	newWalletListResponse := func(slice []*handlers.Wallet) *walletListResponse {
		return &walletListResponse{
			Count:   len(slice),
			Entries: slice,
		}
	}

	t.Run("Missing/Invalid authorization token cases", func(t *testing.T) {
		token := "invalid-token"

		missingTokenReq := NewListWalletsRequest(token)
		invalidTokenReq := NewListWalletsRequest(token)

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

		t.Run("List wallets when there are no wallets", func(t *testing.T) {
			wallets := []*model.Wallet{}
			repoSpy.On("WalletList", userID).Return(wallets, nil).Once()

			res := httptest.NewRecorder()
			req := NewListWalletsRequest(token)

			r.ServeHTTP(res, req)

			expected := newWalletListResponse([]*handlers.Wallet{})

			AssertStatusCode(t, res, http.StatusOK)
			AssertResponseBody(t, res, expected)
		})

		t.Run("List wallets when there are non-zero wallets", func(t *testing.T) {
			wallets := []*model.Wallet{{}}

			repoSpy.On("WalletList", userID).Return(wallets, nil).Once()

			res := httptest.NewRecorder()
			req := NewListWalletsRequest(token)

			r.ServeHTTP(res, req)

			expected := newWalletListResponse([]*handlers.Wallet{{}})

			AssertStatusCode(t, res, http.StatusOK)
			AssertResponseBody(t, res, expected)
		})
	})
}

func TestListTransactionsByWallet(t *testing.T) {
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
		id := uint(1)
		token := "invalid-token"

		missingTokenReq := NewListTransactionsByWalletRequest(id, token)
		invalidTokenReq := NewListTransactionsByWalletRequest(id, token)

		unauthorizedTestCases := UnauthorizedTestCases(missingTokenReq, invalidTokenReq, r, jwtServiceSpy)
		t.Run("Unauthorized test cases", unauthorizedTestCases)
	})

	t.Run("Valid authorization token cases", func(t *testing.T) {
		id := uint(1)
		token := "valid-token"
		userID := uint(1)
		claims := auth.CustomClaims{
			ID: userID,
		}
		jwtServiceSpy.On("ValidateJWT", token).Return(&claims, nil)

		t.Run("List transactions of a non-existent wallet", func(t *testing.T) {
			repoSpy.On("WalletGet", id).Return(nil, repository.ErrorRecordNotFound).Once()

			res := httptest.NewRecorder()
			req := NewListTransactionsByWalletRequest(id, token)

			r.ServeHTTP(res, req)

			AssertStatusCode(t, res, http.StatusNotFound)
		})

		t.Run("List transactions of a wallet that belongs to another user", func(t *testing.T) {
			wallet := &model.Wallet{
				UserID: userID + 1,
			}
			repoSpy.On("WalletGet", id).Return(wallet, nil).Once()

			res := httptest.NewRecorder()
			req := NewListTransactionsByWalletRequest(id, token)

			r.ServeHTTP(res, req)

			AssertStatusCode(t, res, http.StatusForbidden)
		})

		t.Run("List transactions when there are no transactions", func(t *testing.T) {
			wallet := &model.Wallet{
				UserID: userID,
			}
			transactions := []*model.Transaction{}

			repoSpy.On("WalletGet", id).Return(wallet, nil).Once()
			repoSpy.On("TransactionListByWallet", userID, id).Return(transactions, nil).Once()

			res := httptest.NewRecorder()
			req := NewListTransactionsByWalletRequest(id, token)

			r.ServeHTTP(res, req)

			expected := newTransactionListResponse([]*handlers.Transaction{})

			AssertStatusCode(t, res, http.StatusOK)
			AssertResponseBody(t, res, expected)
		})

		t.Run("List transactions when there are non-zero transactions", func(t *testing.T) {
			wallet := &model.Wallet{
				UserID: userID,
			}
			transactions := []*model.Transaction{{}}

			repoSpy.On("WalletGet", id).Return(wallet, nil).Once()
			repoSpy.On("TransactionListByWallet", userID, id).Return(transactions, nil).Once()

			res := httptest.NewRecorder()
			req := NewListTransactionsByWalletRequest(id, token)

			r.ServeHTTP(res, req)

			expected := newTransactionListResponse([]*handlers.Transaction{{}})

			AssertStatusCode(t, res, http.StatusOK)
			AssertResponseBody(t, res, expected)
		})
	})
}
