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

type Response interface {
	handlers.LoginToken |
		handlers.Account |
		handlers.Party |
		handlers.Wallet |
		handlers.Transaction |
		PartyListResponse |
		WalletListResponse |
		TransactionListResponse
}

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

func AssertResponseBody[R Response](t *testing.T, res *httptest.ResponseRecorder, expected *R) {
	t.Helper()
	var got R
	ParseJSONtoResponse(t, res, &got)
	AssertEqual(t, got, *expected)
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

func ParseJSONtoResponse[R Response](t *testing.T, res *httptest.ResponseRecorder, got *R) {
	t.Helper()
	if err := json.NewDecoder(res.Body).Decode(got); err != nil {
		t.Errorf("couldn't parse json response: %v", err)
	}
}

// Other
func createRequestBody(model interface{}) []byte {
	body, _ := json.Marshal(model)
	return body
}
