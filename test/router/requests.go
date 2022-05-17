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

var AccountRequestFactory = map[string]func(token string, user *handlers.Account) *http.Request{
	"get": func(token string, _ *handlers.Account) *http.Request {
		return NewRequest(http.MethodGet, BaseAccountPath, token, nil)
	},
	"update": func(token string, user *handlers.Account) *http.Request {
		return NewRequest(http.MethodPatch, BaseAccountPath, token, user)
	},
	"delete": func(token string, _ *handlers.Account) *http.Request {
		return NewRequest(http.MethodDelete, BaseAccountPath, token, nil)
	},
}

var AuthRequestFactory = map[string]func(handler interface{}) *http.Request{
	"sign_up": func(handler interface{}) *http.Request {
		return NewRequest(http.MethodPost, BaseAuthPath+"/signup", "", handler)
	},
	"login": func(handler interface{}) *http.Request {
		return NewRequest(http.MethodPost, BaseAuthPath+"/login", "", handler)
	},
}

// Parties
var PartyRequestFactory = map[string]func(token string, id uint, party *handlers.Party) *http.Request{
	"create": func(token string, _ uint, party *handlers.Party) *http.Request {
		return NewRequest(http.MethodPost, BasePartiesPath, token, party)
	},
	"get": func(token string, id uint, _ *handlers.Party) *http.Request {
		return NewRequest(http.MethodGet, fmt.Sprintf("%s%d", BasePartiesPath, id), token, nil)
	},
	"update": func(token string, id uint, party *handlers.Party) *http.Request {
		return NewRequest(http.MethodPatch, fmt.Sprintf("%s%d", BasePartiesPath, id), token, party)
	},
	"delete": func(token string, id uint, _ *handlers.Party) *http.Request {
		return NewRequest(http.MethodDelete, fmt.Sprintf("%s%d", BasePartiesPath, id), token, nil)
	},
	"list_all": func(token string, _ uint, _ *handlers.Party) *http.Request {
		return NewRequest(http.MethodGet, BasePartiesPath, token, nil)
	},
	"list_by_party_request": func(token string, id uint, _ *handlers.Party) *http.Request {
		return NewRequest(http.MethodGet, fmt.Sprintf("%s%d/transactions", BasePartiesPath, id), token, nil)
	},
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
