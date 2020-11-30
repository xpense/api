package router

import (
	"encoding/json"
	"expense-api/internal/handlers"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
)

type (
	PartyListResponse struct {
		Count   int               `json:"count"`
		Entries []*handlers.Party `json:"entries"`
	}

	TransactionListResponse struct {
		Count   int                     `json:"count"`
		Entries []*handlers.Transaction `json:"entries"`
	}

	WalletListResponse struct {
		Count   int                `json:"count"`
		Entries []*handlers.Wallet `json:"entries"`
	}
)

// Assertions
func AssertEqual(t *testing.T, got, expected interface{}) {
	t.Helper()

	if !cmp.Equal(got, expected) {
		t.Errorf("expected %+v, got %+v", expected, got)
	}
}

func AssertStatusCode(t *testing.T, res *httptest.ResponseRecorder, expectedStatusCode int) {
	t.Helper()

	if res.Code != expectedStatusCode {
		t.Errorf("expected status %v, got %v", expectedStatusCode, res.Code)
	}
}

func AssertErrorMessage(t *testing.T, res *httptest.ResponseRecorder, expected string) {
	t.Helper()

	jsonResponse := ParseJSON(t, res)
	got := jsonResponse["message"].(string)

	notEqualMsg := fmt.Sprintf("expected error message: '%v', instead got error message: '%v'", expected, got)

	assert.Equal(t, expected, got, notEqualMsg)
}

func AssertResponseBody(t *testing.T, res *httptest.ResponseRecorder, expected interface{}) {
	t.Helper()

	switch expected := expected.(type) {
	// Single entities
	case *handlers.LoginToken:
		var got handlers.LoginToken
		ParseLoginTokenJSONResponse(t, res, &got)
		AssertEqual(t, got, *expected)

	case *handlers.Account:
		var got handlers.Account
		ParseAccountJSONResponse(t, res, &got)
		AssertEqual(t, got, *expected)

	case *handlers.Party:
		var got handlers.Party
		ParsePartyJSONResponse(t, res, &got)
		AssertEqual(t, got, *expected)

	case *handlers.Wallet:
		var got handlers.Wallet
		ParseWalletJSONResponse(t, res, &got)
		AssertEqual(t, got, *expected)

	case *handlers.Transaction:
		var got handlers.Transaction
		ParseTransactionJSONResponse(t, res, &got)
		AssertEqual(t, got, *expected)

	// List of entities
	case *PartyListResponse:
		var got PartyListResponse
		ParsePartyListJSONResponse(t, res, &got)
		AssertEqual(t, got, *expected)

	case *WalletListResponse:
		var got WalletListResponse
		ParseWalletListJSONResponse(t, res, &got)
		AssertEqual(t, got, *expected)

	case *TransactionListResponse:
		var got TransactionListResponse
		ParseTransactionListJSONResponse(t, res, &got)
		AssertEqual(t, got, *expected)

	default:
		t.Errorf("unknown type %v", expected)
	}
}

// Parsing
func ParseJSON(t *testing.T, res *httptest.ResponseRecorder) map[string]interface{} {
	t.Helper()

	jsonResponse := map[string]interface{}{}
	if err := json.NewDecoder(res.Body).Decode(&jsonResponse); err != nil {
		t.Errorf("couldn't parse json response: %v", err)
	}

	return jsonResponse
}

func ParseLoginTokenJSONResponse(t *testing.T, res *httptest.ResponseRecorder, into *handlers.LoginToken) {
	if err := json.NewDecoder(res.Body).Decode(&into); err != nil {
		t.Errorf("couldn't parse json response: %v", err)
	}
}

func ParseAccountJSONResponse(t *testing.T, res *httptest.ResponseRecorder, into *handlers.Account) {
	if err := json.NewDecoder(res.Body).Decode(&into); err != nil {
		t.Errorf("couldn't parse json response: %v", err)
	}
}

func ParsePartyJSONResponse(t *testing.T, res *httptest.ResponseRecorder, into *handlers.Party) {
	if err := json.NewDecoder(res.Body).Decode(&into); err != nil {
		t.Errorf("couldn't parse json response: %v", err)
	}
}

func ParseWalletJSONResponse(t *testing.T, res *httptest.ResponseRecorder, into *handlers.Wallet) {
	if err := json.NewDecoder(res.Body).Decode(&into); err != nil {
		t.Errorf("couldn't parse json response: %v", err)
	}
}

func ParseTransactionJSONResponse(t *testing.T, res *httptest.ResponseRecorder, into *handlers.Transaction) {
	if err := json.NewDecoder(res.Body).Decode(&into); err != nil {
		t.Errorf("couldn't parse json response: %v", err)
	}
}

func ParsePartyListJSONResponse(t *testing.T, res *httptest.ResponseRecorder, into *PartyListResponse) {
	if err := json.NewDecoder(res.Body).Decode(&into); err != nil {
		t.Errorf("couldn't parse json response: %v", err)
	}
}

func ParseWalletListJSONResponse(t *testing.T, res *httptest.ResponseRecorder, into *WalletListResponse) {
	if err := json.NewDecoder(res.Body).Decode(&into); err != nil {
		t.Errorf("couldn't parse json response: %v", err)
	}
}

func ParseTransactionListJSONResponse(t *testing.T, res *httptest.ResponseRecorder, into *TransactionListResponse) {
	if err := json.NewDecoder(res.Body).Decode(&into); err != nil {
		t.Errorf("couldn't parse json response: %v", err)
	}
}

// Other
func createRequestBody(model interface{}) []byte {
	body, _ := json.Marshal(model)
	return body
}
