package router

import (
	"bytes"
	"expense-api/internal/handlers"
	"fmt"
	"net/http"
)

var (
	BasePath             = "/api/v1"
	BaseAccountPath      = BasePath + "/account/"
	BaseAuthPath         = BasePath + "/auth"
	BasePartiesPath      = BasePath + "/parties/"
	BaseTransactionsPath = BasePath + "/transactions/"
	BaseWalletsPath      = BasePath + "/wallets/"
)

func NewRequest(method, path, token string, handler interface{}) *http.Request {
	body := createRequestBody(handler)
	req, _ := http.NewRequest(method, path, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	return req
}

// Account
func NewGetAccountRequest(token string) *http.Request {
	return NewRequest(http.MethodGet, BaseAccountPath, token, nil)
}

func NewUpdateAccountRequest(user *handlers.Account, token string) *http.Request {
	return NewRequest(http.MethodPatch, BaseAccountPath, token, user)
}

func NewDeleteAccountRequest(token string) *http.Request {
	return NewRequest(http.MethodDelete, BaseAccountPath, token, nil)
}

// Auth
func NewAuthRequest(handler interface{}) *http.Request {
	path := BaseAuthPath
	switch handler.(type) {
	case *handlers.SignUpInfo:
		path += "/signup"
	case *handlers.LoginInfo:
		path += "/login"
	}
	return NewRequest(http.MethodPost, path, "", handler)
}

// Parties
func NewCreatePartyRequest(party *handlers.Party, token string) *http.Request {
	return NewRequest(http.MethodPost, BasePartiesPath, token, party)
}

func NewGetPartyRequest(id uint, token string) *http.Request {
	return NewRequest(http.MethodGet, fmt.Sprintf("%s%d", BasePartiesPath, id), token, nil)
}

func NewUpdatePartyRequest(id uint, party *handlers.Party, token string) *http.Request {
	return NewRequest(http.MethodPatch, fmt.Sprintf("%s%d", BasePartiesPath, id), token, party)
}

func NewDeletePartyRequest(id uint, token string) *http.Request {
	return NewRequest(http.MethodDelete, fmt.Sprintf("%s%d", BasePartiesPath, id), token, nil)
}

func NewListPartiesRequest(token string) *http.Request {
	return NewRequest(http.MethodGet, BasePartiesPath, token, nil)
}

func NewListTransactionsByPartyRequest(id uint, token string) *http.Request {
	return NewRequest(http.MethodGet, fmt.Sprintf("%s%d/transactions", BasePartiesPath, id), token, nil)
}

// Transactions
func NewCreateTransactionRequest(transaction *handlers.Transaction, token string) *http.Request {
	return NewRequest(http.MethodPost, BaseTransactionsPath, token, transaction)
}

func NewGetTransactionRequest(id uint, token string) *http.Request {
	return NewRequest(http.MethodGet, fmt.Sprintf("%s%d", BaseTransactionsPath, id), token, nil)
}

func NewUpdateTransactionRequest(id uint, transaction *handlers.Transaction, token string) *http.Request {
	return NewRequest(http.MethodPatch, fmt.Sprintf("%s%d", BaseTransactionsPath, id), token, transaction)
}

func NewDeleteTransactionRequest(id uint, token string) *http.Request {
	return NewRequest(http.MethodDelete, fmt.Sprintf("%s%d", BaseTransactionsPath, id), token, nil)
}

func NewListTransactionsRequest(token string) *http.Request {
	return NewRequest(http.MethodGet, BaseTransactionsPath, token, nil)
}

// Wallets
func NewCreateWalletRequest(wallet *handlers.Wallet, token string) *http.Request {
	return NewRequest(http.MethodPost, BaseWalletsPath, token, wallet)
}

func NewGetWalletRequest(id uint, token string) *http.Request {
	return NewRequest(http.MethodGet, fmt.Sprintf("%s%d", BaseWalletsPath, id), token, nil)
}

func NewUpdateWalletRequest(id uint, wallet *handlers.Wallet, token string) *http.Request {
	return NewRequest(http.MethodPatch, fmt.Sprintf("%s%d", BaseWalletsPath, id), token, wallet)
}

func NewDeleteWalletRequest(id uint, token string) *http.Request {
	return NewRequest(http.MethodDelete, fmt.Sprintf("%s%d", BaseWalletsPath, id), token, nil)
}

func NewListWalletsRequest(token string) *http.Request {
	return NewRequest(http.MethodGet, BaseWalletsPath, token, nil)
}

func NewListTransactionsByWalletRequest(id uint, token string) *http.Request {
	return NewRequest(http.MethodGet, fmt.Sprintf("%s%d/transactions", BaseWalletsPath, id), token, nil)
}
