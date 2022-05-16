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

// Account
func NewAccountRequest(method string, token string, user *handlers.Account) *http.Request {
	body := createRequestBody(user)
	req, _ := http.NewRequest(method, BaseAccountPath, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	return req
}

func NewGetAccountRequest(token string) *http.Request {
	return NewAccountRequest(http.MethodGet, token, nil)
}

func NewUpdateAccountRequest(user *handlers.Account, token string) *http.Request {
	return NewAccountRequest(http.MethodPatch, token, user)
}

func NewDeleteAccountRequest(token string) *http.Request {
	return NewAccountRequest(http.MethodDelete, token, nil)
}

// Auth
func NewAuthRequest(handler interface{}) *http.Request {
	body := createRequestBody(handler)
	path := ""
	switch handler.(type) {
	case *handlers.SignUpInfo:
		path += "/signup"
	case *handlers.LoginInfo:
		path += "/login"
	}
	req, _ := http.NewRequest(http.MethodPost, BaseAuthPath+path, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	return req
}

// Parties
func NewPartyRequest(method string, path string, token string, party *handlers.Party) *http.Request {
	body := createRequestBody(party)
	req, _ := http.NewRequest(method, BasePartiesPath+path, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	return req
}

func NewCreatePartyRequest(party *handlers.Party, token string) *http.Request {
	return NewPartyRequest(http.MethodPost, "", token, party)
}

func NewGetPartyRequest(id uint, token string) *http.Request {
	return NewPartyRequest(http.MethodGet, fmt.Sprintf("%d", id), token, nil)
}

func NewUpdatePartyRequest(id uint, party *handlers.Party, token string) *http.Request {
	return NewPartyRequest(http.MethodPatch, fmt.Sprintf("%d", id), token, party)
}

func NewDeletePartyRequest(id uint, token string) *http.Request {
	return NewPartyRequest(http.MethodDelete, fmt.Sprintf("%d", id), token, nil)
}

func NewListPartiesRequest(token string) *http.Request {
	return NewPartyRequest(http.MethodGet, "", token, nil)
}

func NewListTransactionsByPartyRequest(id uint, token string) *http.Request {
	return NewPartyRequest(http.MethodGet, fmt.Sprintf("%d/transactions", id), token, nil)
}

// Transactions
func NewTransactionRequest(method string, path string, token string, transaction *handlers.Transaction) *http.Request {
	body := createRequestBody(transaction)
	req, _ := http.NewRequest(method, BaseTransactionsPath+path, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	return req
}

func NewCreateTransactionRequest(transaction *handlers.Transaction, token string) *http.Request {
	return NewTransactionRequest(http.MethodPost, "", token, transaction)
}

func NewGetTransactionRequest(id uint, token string) *http.Request {
	return NewTransactionRequest(http.MethodGet, fmt.Sprintf("%d", id), token, nil)
}

func NewUpdateTransactionRequest(id uint, transaction *handlers.Transaction, token string) *http.Request {
	return NewTransactionRequest(http.MethodPatch, fmt.Sprintf("%d", id), token, transaction)
}

func NewDeleteTransactionRequest(id uint, token string) *http.Request {
	return NewTransactionRequest(http.MethodDelete, fmt.Sprintf("%d", id), token, nil)
}

func NewListTransactionsRequest(token string) *http.Request {
	return NewTransactionRequest(http.MethodGet, "", token, nil)
}

// Wallets
func NewWalletRequest(method, path string, token string, wallet *handlers.Wallet) *http.Request {
	body := createRequestBody(wallet)
	req, _ := http.NewRequest(method, BaseWalletsPath+path, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	return req
}

func NewCreateWalletRequest(wallet *handlers.Wallet, token string) *http.Request {
	return NewWalletRequest(http.MethodPost, "", token, wallet)
}

func NewGetWalletRequest(id uint, token string) *http.Request {
	return NewWalletRequest(http.MethodGet, fmt.Sprintf("%d", id), token, nil)
}

func NewUpdateWalletRequest(id uint, wallet *handlers.Wallet, token string) *http.Request {
	return NewWalletRequest(http.MethodPatch, fmt.Sprintf("%d", id), token, wallet)
}

func NewDeleteWalletRequest(id uint, token string) *http.Request {
	return NewWalletRequest(http.MethodDelete, fmt.Sprintf("%d", id), token, nil)
}

func NewListWalletsRequest(token string) *http.Request {
	return NewWalletRequest(http.MethodGet, "", token, nil)
}

func NewListTransactionsByWalletRequest(id uint, token string) *http.Request {
	return NewWalletRequest(http.MethodGet, fmt.Sprintf("%d/transactions", id), token, nil)
}
