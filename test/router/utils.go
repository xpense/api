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
	partyListResponse struct {
		Count   int               `json:"count"`
		Entries []*handlers.Party `json:"entries"`
	}

	transactionListResponse struct {
		Count   int                     `json:"count"`
		Entries []*handlers.Transaction `json:"entries"`
	}

	walletListResponse struct {
		Count   int                `json:"count"`
		Entries []*handlers.Wallet `json:"entries"`
	}
)

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

	// Single entity
	switch expected := expected.(type) {
	case *handlers.Account:
		var got handlers.Account
		if err := json.NewDecoder(res.Body).Decode(&got); err != nil {
			t.Errorf("couldn't parse json response: %v", err)
		}

		AssertEqual(t, got, *expected)

	case *handlers.Party:
		var got handlers.Party
		if err := json.NewDecoder(res.Body).Decode(&got); err != nil {
			t.Errorf("couldn't parse json response: %v", err)
		}

		AssertEqual(t, got, *expected)

	case *handlers.Wallet:
		var got handlers.Wallet
		if err := json.NewDecoder(res.Body).Decode(&got); err != nil {
			t.Errorf("couldn't parse json response: %v", err)
		}

		AssertEqual(t, got, *expected)

	case *handlers.Transaction:
		var got handlers.Transaction
		if err := json.NewDecoder(res.Body).Decode(&got); err != nil {
			t.Errorf("couldn't parse json response: %v", err)
		}

		AssertEqual(t, got, *expected)
	}

	// List of entities
	switch expected := expected.(type) {
	case *partyListResponse:
		var got partyListResponse
		if err := json.NewDecoder(res.Body).Decode(&got); err != nil {
			t.Errorf("couldn't parse json response: %v", err)
		}

		AssertEqual(t, got, *expected)

	case *walletListResponse:
		var got walletListResponse
		if err := json.NewDecoder(res.Body).Decode(&got); err != nil {
			t.Errorf("couldn't parse json response: %v", err)
		}

		AssertEqual(t, got, *expected)

	case *transactionListResponse:
		var got transactionListResponse
		if err := json.NewDecoder(res.Body).Decode(&got); err != nil {
			t.Errorf("couldn't parse json response: %v", err)
		}

		AssertEqual(t, got, *expected)
	}
}

func ParseJSON(t *testing.T, res *httptest.ResponseRecorder) map[string]interface{} {
	t.Helper()

	jsonResponse := map[string]interface{}{}
	if err := json.NewDecoder(res.Body).Decode(&jsonResponse); err != nil {
		t.Errorf("couldn't parse json response: %v", err)
	}

	return jsonResponse
}

func createRequestBody(model interface{}) []byte {
	body, _ := json.Marshal(model)
	return body
}
