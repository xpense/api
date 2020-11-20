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

	"github.com/google/go-cmp/cmp"
)

func TestCreateWallet(t *testing.T) {
	repoSpy := &spies.RepositorySpy{}
	jwtServiceSpy := &spies.JWTServiceSpy{}
	hasherSpy := &spies.PasswordHasherSpy{}

	r := router.Setup(repoSpy, jwtServiceSpy, hasherSpy, router.TestConfig)

	newWalletRequest := func(wallet *model.Wallet, token string) *http.Request {
		body := createRequestBody(wallet)
		req, _ := http.NewRequest(http.MethodPost, "/wallet/", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		return req
	}

	t.Run("Missing/Invalid authorization token cases", func(t *testing.T) {
		wallet := &model.Wallet{}
		token := "invalid-token"

		missingTokenReq := newWalletRequest(wallet, token)
		invalidTokenReq := newWalletRequest(wallet, token)

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
			req := newWalletRequest(wallet, token)

			r.ServeHTTP(res, req)

			wantErrorMessage := handlers.ErrMsgWalletNameTaken

			assertStatusCode(t, res, http.StatusBadRequest)
			assertErrorMessage(t, res, wantErrorMessage)
		})

		t.Run("Create wallet with valid data", func(t *testing.T) {
			wallet := &model.Wallet{
				UserID: userID,
			}

			repoSpy.On("WalletCreate", wallet).Return(nil).Once()

			res := httptest.NewRecorder()
			req := newWalletRequest(wallet, token)

			r.ServeHTTP(res, req)

			wallet.UserID = 0

			assertStatusCode(t, res, http.StatusCreated)
			assertSingleWalletResponseBody(t, res, wallet)
		})
	})
}

func TestGetWallet(t *testing.T) {
	repoSpy := &spies.RepositorySpy{}
	jwtServiceSpy := &spies.JWTServiceSpy{}
	hasherSpy := &spies.PasswordHasherSpy{}

	r := router.Setup(repoSpy, jwtServiceSpy, hasherSpy, router.TestConfig)

	newWalletRequest := func(id uint, token string) *http.Request {
		url := fmt.Sprintf("/wallet/%d", id)
		req, _ := http.NewRequest(http.MethodGet, url, nil)
		req.Header.Set("Authorization", "Bearer "+token)
		return req
	}

	t.Run("Missing/Invalid authorization token cases", func(t *testing.T) {
		id := uint(1)
		token := "invalid-token"

		missingTokenReq := newWalletRequest(id, token)
		invalidTokenReq := newWalletRequest(id, token)

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
			req := newWalletRequest(id, token)

			r.ServeHTTP(res, req)

			assertStatusCode(t, res, http.StatusBadRequest)
		})

		t.Run("Get wallet with non-existent id", func(t *testing.T) {
			id := uint(10)

			repoSpy.On("WalletGet", id).Return(nil, repository.ErrorRecordNotFound).Once()

			res := httptest.NewRecorder()
			req := newWalletRequest(id, token)

			r.ServeHTTP(res, req)

			assertStatusCode(t, res, http.StatusNotFound)
		})

		t.Run("Get wallet with valid id that belongs to another user", func(t *testing.T) {
			id := uint(1)
			wallet := &model.Wallet{
				Name:   "new wallet",
				UserID: userID + 1,
			}

			repoSpy.On("WalletGet", id).Return(wallet, nil).Once()

			res := httptest.NewRecorder()
			req := newWalletRequest(id, token)

			r.ServeHTTP(res, req)

			assertStatusCode(t, res, http.StatusUnauthorized)
		})

		t.Run("Get wallet with valid id", func(t *testing.T) {
			id := uint(1)
			wallet := &model.Wallet{
				Name:   "new wallet",
				UserID: userID,
			}

			repoSpy.On("WalletGet", id).Return(wallet, nil).Twice()

			res := httptest.NewRecorder()
			req := newWalletRequest(id, token)

			r.ServeHTTP(res, req)

			wallet.UserID = 0

			assertStatusCode(t, res, http.StatusOK)
			assertSingleWalletResponseBody(t, res, wallet)
		})
	})
}

func TestUpdateWallet(t *testing.T) {
	repoSpy := &spies.RepositorySpy{}
	jwtServiceSpy := &spies.JWTServiceSpy{}
	hasherSpy := &spies.PasswordHasherSpy{}

	r := router.Setup(repoSpy, jwtServiceSpy, hasherSpy, router.TestConfig)

	newWalletRequest := func(id uint, wallet *model.Wallet, token string) *http.Request {
		url := fmt.Sprintf("/wallet/%d", id)
		body := createRequestBody(wallet)
		req, _ := http.NewRequest(http.MethodPatch, url, bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		return req
	}

	t.Run("Missing/Invalid authorization token cases", func(t *testing.T) {
		id := uint(1)
		wallet := &model.Wallet{}
		token := "invalid-token"

		missingTokenReq := newWalletRequest(id, wallet, token)
		invalidTokenReq := newWalletRequest(id, wallet, token)

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
			wallet := &model.Wallet{
				Name:   "new wallet",
				UserID: userID,
			}

			repoSpy.On("WalletGet", id).Return(nil, repository.ErrorRecordNotFound).Once()

			res := httptest.NewRecorder()
			req := newWalletRequest(id, wallet, token)

			r.ServeHTTP(res, req)

			assertStatusCode(t, res, http.StatusNotFound)
		})

		t.Run("Try to update wallet with valid id that belongs to another user", func(t *testing.T) {
			id := uint(1)
			wallet := &model.Wallet{
				Name:   "new wallet",
				UserID: userID + 1,
			}

			repoSpy.On("WalletGet", id).Return(wallet, nil).Once()

			res := httptest.NewRecorder()
			req := newWalletRequest(id, wallet, token)

			r.ServeHTTP(res, req)

			assertStatusCode(t, res, http.StatusUnauthorized)
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
			req := newWalletRequest(id, wallet, token)

			r.ServeHTTP(res, req)

			wantErrorMessage := handlers.ErrMsgWalletNameTaken

			assertStatusCode(t, res, http.StatusBadRequest)
			assertErrorMessage(t, res, wantErrorMessage)
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
			req := newWalletRequest(id, wallet, token)

			r.ServeHTTP(res, req)

			wallet.UserID = 0

			assertStatusCode(t, res, http.StatusOK)
			assertSingleWalletResponseBody(t, res, wallet)
		})
	})
}

func TestDeleteWallet(t *testing.T) {
	repoSpy := &spies.RepositorySpy{}
	jwtServiceSpy := &spies.JWTServiceSpy{}
	hasherSpy := &spies.PasswordHasherSpy{}

	r := router.Setup(repoSpy, jwtServiceSpy, hasherSpy, router.TestConfig)

	newWalletRequest := func(id uint, token string) *http.Request {
		url := fmt.Sprintf("/wallet/%d", id)
		req, _ := http.NewRequest(http.MethodDelete, url, nil)
		req.Header.Set("Authorization", "Bearer "+token)
		return req
	}

	t.Run("Missing/Invalid authorization token cases", func(t *testing.T) {
		id := uint(1)
		token := "invalid-token"

		missingTokenReq := newWalletRequest(id, token)
		invalidTokenReq := newWalletRequest(id, token)

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
			req := newWalletRequest(id, token)

			r.ServeHTTP(res, req)

			assertStatusCode(t, res, http.StatusNotFound)
		})

		t.Run("Try to delete wallet with valid id that belongs to another user", func(t *testing.T) {
			id := uint(1)
			wallet := &model.Wallet{
				Name:   "new wallet",
				UserID: userID + 1,
			}

			repoSpy.On("WalletGet", id).Return(wallet, nil).Once()

			res := httptest.NewRecorder()
			req := newWalletRequest(id, token)

			r.ServeHTTP(res, req)

			assertStatusCode(t, res, http.StatusUnauthorized)
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
			req := newWalletRequest(id, token)

			r.ServeHTTP(res, req)

			assertStatusCode(t, res, http.StatusNoContent)
		})
	})
}

func TestListWallets(t *testing.T) {
	repoSpy := &spies.RepositorySpy{}
	jwtServiceSpy := &spies.JWTServiceSpy{}
	hasherSpy := &spies.PasswordHasherSpy{}

	r := router.Setup(repoSpy, jwtServiceSpy, hasherSpy, router.TestConfig)

	newWalletListResponse := func(slice []*model.Wallet) *walletListResponse {
		return &walletListResponse{
			Count:   len(slice),
			Entries: slice,
		}
	}

	newWalletRequest := func(token string) *http.Request {
		req, _ := http.NewRequest(http.MethodGet, "/wallet/", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		return req
	}

	t.Run("Missing/Invalid authorization token cases", func(t *testing.T) {
		token := "invalid-token"

		missingTokenReq := newWalletRequest(token)
		invalidTokenReq := newWalletRequest(token)

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
			req := newWalletRequest(token)

			r.ServeHTTP(res, req)

			expected := newWalletListResponse(wallets)

			assertStatusCode(t, res, http.StatusOK)
			assertListWalletResponseBody(t, res, expected)
		})

		t.Run("List wallets when there are non-zero wallets", func(t *testing.T) {
			wallets := []*model.Wallet{{}}

			repoSpy.On("WalletList", userID).Return(wallets, nil).Once()

			res := httptest.NewRecorder()
			req := newWalletRequest(token)

			r.ServeHTTP(res, req)

			expected := newWalletListResponse(wallets)

			assertStatusCode(t, res, http.StatusOK)
			assertListWalletResponseBody(t, res, expected)
		})
	})
}

func TestListTransactionsByWallet(t *testing.T) {
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

	newWalletRequest := func(id uint, token string) *http.Request {
		url := fmt.Sprintf("/wallet/%d/transaction", id)
		req, _ := http.NewRequest(http.MethodGet, url, nil)
		req.Header.Set("Authorization", "Bearer "+token)
		return req
	}

	t.Run("Missing/Invalid authorization token cases", func(t *testing.T) {
		id := uint(1)
		token := "invalid-token"

		missingTokenReq := newWalletRequest(id, token)
		invalidTokenReq := newWalletRequest(id, token)

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
			req := newWalletRequest(id, token)

			r.ServeHTTP(res, req)

			assertStatusCode(t, res, http.StatusNotFound)
		})

		t.Run("List transactions of a wallet that belongs to another user", func(t *testing.T) {
			wallet := &model.Wallet{
				UserID: userID + 1,
			}
			repoSpy.On("WalletGet", id).Return(wallet, nil).Once()

			res := httptest.NewRecorder()
			req := newWalletRequest(id, token)

			r.ServeHTTP(res, req)

			assertStatusCode(t, res, http.StatusUnauthorized)
		})

		t.Run("List transactions when there are no transactions", func(t *testing.T) {
			wallet := &model.Wallet{
				UserID: userID,
			}
			transactions := []*model.Transaction{}

			repoSpy.On("WalletGet", id).Return(wallet, nil).Once()
			repoSpy.On("TransactionListByWallet", userID, id).Return(transactions, nil).Once()

			res := httptest.NewRecorder()
			req := newWalletRequest(id, token)

			r.ServeHTTP(res, req)

			expected := newTransactionListResponse(transactions)

			assertStatusCode(t, res, http.StatusOK)
			assertListTransactionResponseBody(t, res, expected)
		})

		t.Run("List transactions when there are non-zero transactions", func(t *testing.T) {
			wallet := &model.Wallet{
				UserID: userID,
			}
			transactions := []*model.Transaction{{}}

			repoSpy.On("WalletGet", id).Return(wallet, nil).Once()
			repoSpy.On("TransactionListByWallet", userID, id).Return(transactions, nil).Once()

			res := httptest.NewRecorder()
			req := newWalletRequest(id, token)

			r.ServeHTTP(res, req)

			expected := newTransactionListResponse(transactions)

			assertStatusCode(t, res, http.StatusOK)
			assertListTransactionResponseBody(t, res, expected)
		})
	})
}

type walletListResponse struct {
	Count   int             `json:"count"`
	Entries []*model.Wallet `json:"entries"`
}

func assertSingleWalletResponseBody(t *testing.T, res *httptest.ResponseRecorder, wallet *model.Wallet) {
	t.Helper()

	var got model.Wallet
	if err := json.NewDecoder(res.Body).Decode(&got); err != nil {
		t.Errorf("couldn't parse json response: %v", err)
	}

	if !cmp.Equal(got, *wallet) {
		t.Errorf("expected %+v, got %+v", *wallet, got)
	}
}

func assertListWalletResponseBody(t *testing.T, res *httptest.ResponseRecorder, expected *walletListResponse) {
	t.Helper()

	var got walletListResponse
	if err := json.NewDecoder(res.Body).Decode(&got); err != nil {
		t.Errorf("couldn't parse json response: %v", err)
	}

	if !cmp.Equal(got, *expected) {
		t.Errorf("expected %+v, got %+v", *expected, got)
	}
}
