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

func createRequestBody(model interface{}) []byte {
	body, _ := json.Marshal(model)
	return body
}
