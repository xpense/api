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

func NewCreatePartyRequest(party *handlers.Party, token string) *http.Request {
	body := createRequestBody(party)
	req, _ := http.NewRequest(http.MethodPost, BasePartiesPath, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	return req
}

func NewGetPartyRequest(id uint, token string) *http.Request {
	url := fmt.Sprintf("%s%d", BasePartiesPath, id)
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	return req
}

func NewUpdatePartyRequest(id uint, party *handlers.Party, token string) *http.Request {
	url := fmt.Sprintf("%s%d", BasePartiesPath, id)
	body := createRequestBody(party)
	req, _ := http.NewRequest(http.MethodPatch, url, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	return req
}

func NewDeletePartyRequest(id uint, token string) *http.Request {
	url := fmt.Sprintf("%s%d", BasePartiesPath, id)
	req, _ := http.NewRequest(http.MethodDelete, url, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	return req
}

func NewListPartiesRequest(token string) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, BasePartiesPath, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	return req
}

func NewListTransactionsByPartyRequest(id uint, token string) *http.Request {
	url := fmt.Sprintf("%s%d/transactions", BasePartiesPath, id)
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	return req
}

// Transactions

func NewCreateTransactionRequest(transaction *handlers.Transaction, token string) *http.Request {
	body := createRequestBody(transaction)
	req, _ := http.NewRequest(http.MethodPost, BaseTransactionsPath, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	return req
}

func NewGetTransactionRequest(id uint, token string) *http.Request {
	url := fmt.Sprintf("%s%d", BaseTransactionsPath, id)
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	return req
}

func NewUpdateTransactionRequest(id uint, transaction *handlers.Transaction, token string) *http.Request {
	url := fmt.Sprintf("%s%d", BaseTransactionsPath, id)
	body := createRequestBody(transaction)
	req, _ := http.NewRequest(http.MethodPatch, url, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	return req
}

func NewDeleteTransactionRequest(id uint, token string) *http.Request {
	url := fmt.Sprintf("%s%d", BaseTransactionsPath, id)
	req, _ := http.NewRequest(http.MethodDelete, url, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	return req
}

func NewListTransactionsRequest(token string) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, BaseTransactionsPath, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	return req
}

// Wallets

func NewCreateWalletRequest(wallet *handlers.Wallet, token string) *http.Request {
	body := createRequestBody(wallet)
	req, _ := http.NewRequest(http.MethodPost, BaseWalletsPath, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	return req
}

func NewGetWalletRequest(id uint, token string) *http.Request {
	url := fmt.Sprintf("%s%d", BaseWalletsPath, id)
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	return req
}

func NewUpdateWalletRequest(id uint, wallet *handlers.Wallet, token string) *http.Request {
	url := fmt.Sprintf("%s%d", BaseWalletsPath, id)
	body := createRequestBody(wallet)
	req, _ := http.NewRequest(http.MethodPatch, url, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	return req
}

func NewDeleteWalletRequest(id uint, token string) *http.Request {
	url := fmt.Sprintf("%s%d", BaseWalletsPath, id)
	req, _ := http.NewRequest(http.MethodDelete, url, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	return req
}

func NewListWalletsRequest(token string) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, BaseWalletsPath, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	return req
}

func NewListTransactionsByWalletRequest(id uint, token string) *http.Request {
	url := fmt.Sprintf("%s%d/transactions", BaseWalletsPath, id)
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	return req
}
