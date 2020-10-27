package test

import (
	"bytes"
	"encoding/json"
	"expense-api/model"
	"expense-api/router"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

type TransactionListResponse struct {
	Count   int                  `json:"count"`
	Entries []*model.Transaction `json:"entries"`
}

func newTransactionListResponse(slice []*model.Transaction) *TransactionListResponse {
	return &TransactionListResponse{
		Count:   len(slice),
		Entries: slice,
	}
}

func TestCreateTransaction(t *testing.T) {
	spy := NewRepositorySpy()
	r := router.Setup(spy)

	newTransactionRequest := func(transaction *model.Transaction) *http.Request {
		body := createRequestBody(transaction)
		req, _ := http.NewRequest(http.MethodPost, "/transaction/", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		return req
	}

	t.Run("Create transaction with amount = 0", func(t *testing.T) {
		res := httptest.NewRecorder()
		req := newTransactionRequest(&model.Transaction{Amount: 0})

		r.ServeHTTP(res, req)

		assertStatusCode(t, res, http.StatusBadRequest)
	})

	t.Run("Create transaction with invalid transaction type", func(t *testing.T) {
		res := httptest.NewRecorder()
		req := newTransactionRequest(&model.Transaction{Amount: 1000, Type: "invalid"})

		r.ServeHTTP(res, req)

		assertStatusCode(t, res, http.StatusBadRequest)
	})

	t.Run("Create transaction with valid data", func(t *testing.T) {
		res := httptest.NewRecorder()
		req := newTransactionRequest(&model.Transaction{Amount: 1000, Type: model.Expense})

		r.ServeHTTP(res, req)

		assertStatusCode(t, res, http.StatusCreated)

		var jsonResponse model.Transaction
		parseSingleTransactionBody(t, res, &jsonResponse)
	})
}

func TestGetTransaction(t *testing.T) {
	spy := NewRepositorySpy()
	r := router.Setup(spy)

	firstTransaction, _ := spy.TransactionCreate(time.Now(), 1000, model.Expense)

	newTransactionRequest := func(id int) *http.Request {
		url := fmt.Sprintf("/transaction/%d", id)
		req, _ := http.NewRequest(http.MethodGet, url, nil)
		return req
	}

	t.Run("Get transaction with a negative id", func(t *testing.T) {
		res := httptest.NewRecorder()
		req := newTransactionRequest(-10)

		r.ServeHTTP(res, req)

		assertStatusCode(t, res, http.StatusBadRequest)
	})

	t.Run("Get transaction with id = 0", func(t *testing.T) {
		res := httptest.NewRecorder()
		req := newTransactionRequest(0)

		r.ServeHTTP(res, req)

		assertStatusCode(t, res, http.StatusBadRequest)
	})

	t.Run("Get transaction with non-existent id", func(t *testing.T) {
		res := httptest.NewRecorder()
		req := newTransactionRequest(10)

		r.ServeHTTP(res, req)

		assertStatusCode(t, res, http.StatusNotFound)
	})

	t.Run("Get transaction with valid id", func(t *testing.T) {
		res := httptest.NewRecorder()
		req := newTransactionRequest(1)

		r.ServeHTTP(res, req)

		assertStatusCode(t, res, http.StatusOK)

		var jsonResponse model.Transaction
		parseSingleTransactionBody(t, res, &jsonResponse)

		if !reflect.DeepEqual(jsonResponse, *firstTransaction) {
			t.Errorf("expected %+v, got %+v", *firstTransaction, jsonResponse)
		}
	})
}

func TestUpdateTransaction(t *testing.T) {
	spy := NewRepositorySpy()
	r := router.Setup(spy)

	transaction, _ := spy.TransactionCreate(time.Now(), 1000, model.Expense)

	newTransactionRequest := func(id int, transaction *model.Transaction) *http.Request {
		url := fmt.Sprintf("/transaction/%d", id)
		body := createRequestBody(transaction)
		req, _ := http.NewRequest(http.MethodPatch, url, bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		return req
	}

	t.Run("Update non-existent transaction", func(t *testing.T) {
		res := httptest.NewRecorder()
		req := newTransactionRequest(10, &model.Transaction{Amount: 1100})

		r.ServeHTTP(res, req)

		assertStatusCode(t, res, http.StatusNotFound)
	})

	t.Run("Update existing transaction with invalid type", func(t *testing.T) {
		res := httptest.NewRecorder()
		req := newTransactionRequest(int(transaction.ID), &model.Transaction{Amount: 2, Type: "invalid"})

		r.ServeHTTP(res, req)

		assertStatusCode(t, res, http.StatusBadRequest)
	})

	t.Run("Update existing transaction with valid arguments", func(t *testing.T) {
		got := &model.Transaction{Amount: 2000, Type: model.Income}
		res := httptest.NewRecorder()
		req := newTransactionRequest(int(transaction.ID), got)

		r.ServeHTTP(res, req)

		assertStatusCode(t, res, http.StatusOK)

		if got.Amount != transaction.Amount || got.Type != transaction.Type {
			t.Errorf("Trnsaction fields not updated properly!")
		}
	})
}

func TestDeleteTransaction(t *testing.T) {
	spy := NewRepositorySpy()
	r := router.Setup(spy)

	transaction, _ := spy.TransactionCreate(time.Now(), 1000, model.Expense)

	newTransactionRequest := func(id int) *http.Request {
		url := fmt.Sprintf("/transaction/%d", id)
		req, _ := http.NewRequest(http.MethodDelete, url, nil)
		return req
	}

	t.Run("Delete non-existent transaction", func(t *testing.T) {
		res := httptest.NewRecorder()
		req := newTransactionRequest(10)

		r.ServeHTTP(res, req)

		assertStatusCode(t, res, http.StatusNotFound)
	})

	t.Run("Delete existing transaction", func(t *testing.T) {
		res := httptest.NewRecorder()
		req := newTransactionRequest(int(transaction.ID))

		r.ServeHTTP(res, req)

		assertStatusCode(t, res, http.StatusNoContent)
	})
}

func TestListTransactions(t *testing.T) {
	spy := NewRepositorySpy()
	r := router.Setup(spy)

	newTransactionRequest := func() *http.Request {
		req, _ := http.NewRequest(http.MethodGet, "/transaction/", nil)
		return req
	}

	t.Run("List transactions when there are no transactions", func(t *testing.T) {
		res := httptest.NewRecorder()
		req := newTransactionRequest()

		r.ServeHTTP(res, req)

		assertStatusCode(t, res, http.StatusOK)

		expected := newTransactionListResponse(spy.transactionSlice())

		var got TransactionListResponse
		parseListTransactionBody(t, res, &got)

		if !reflect.DeepEqual(got, *expected) {
			t.Errorf("expected %+v ;%T, got %+v ;%T", *expected, *expected, got, got)
		}
	})

	t.Run("List transactions when there are non-zero transactions", func(t *testing.T) {
		spy.TransactionCreate(time.Now(), 1000, model.Expense)

		res := httptest.NewRecorder()
		req := newTransactionRequest()

		r.ServeHTTP(res, req)

		assertStatusCode(t, res, http.StatusOK)

		slice := spy.transactionSlice()

		expected := newTransactionListResponse(slice)

		var got TransactionListResponse
		parseListTransactionBody(t, res, &got)

		if !reflect.DeepEqual(got, *expected) {
			t.Errorf("expected %+v, got %+v", *expected, got)
		}
	})
}

func parseSingleTransactionBody(t *testing.T, res *httptest.ResponseRecorder, jsonResponse *model.Transaction) {
	t.Helper()

	if err := json.NewDecoder(res.Body).Decode(jsonResponse); err != nil {
		t.Errorf("couldn't parse json response: %v", err)
	}
}

func parseListTransactionBody(t *testing.T, res *httptest.ResponseRecorder, jsonResponse *TransactionListResponse) {
	t.Helper()

	if err := json.NewDecoder(res.Body).Decode(jsonResponse); err != nil {
		t.Errorf("couldn't parse json response: %v", err)
	}
}
