package test

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func assertStatusCode(t *testing.T, res *httptest.ResponseRecorder, expectedStatusCode int) {
	t.Helper()

	if res.Code != expectedStatusCode {
		t.Errorf("expected status %v, got %v", expectedStatusCode, res.Code)
	}
}

func assertErrorMessage(t *testing.T, res *httptest.ResponseRecorder, expected string) {
	t.Helper()

	jsonResponse := parseJSON(t, res)
	got := jsonResponse["message"].(string)

	notEqualMsg := fmt.Sprintf("expected error message: '%v', instead got error message: '%v'", expected, got)

	assert.Equal(t, expected, got, notEqualMsg)
}

func parseJSON(t *testing.T, res *httptest.ResponseRecorder) map[string]interface{} {
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
