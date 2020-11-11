package test

import (
	"encoding/json"
	"net/http/httptest"
	"testing"
)

func assertStatusCode(t *testing.T, res *httptest.ResponseRecorder, expectedStatusCode int) {
	t.Helper()

	if res.Code != expectedStatusCode {
		t.Errorf("expected status %v, got %v", expectedStatusCode, res.Code)
	}
}

func assertErrorMessage(t *testing.T, expected, got string) {
	t.Helper()

	if expected != got {
		t.Errorf("expected error message: '%v', instead got error message: '%v'", expected, got)
	}
}

func parseJSON(t *testing.T, res *httptest.ResponseRecorder) map[string]interface{} {
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
