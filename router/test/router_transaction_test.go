package test

import (
	"bytes"
	"encoding/json"
	"expense-api/model"
	"expense-api/repository"
	"expense-api/router"
	"expense-api/router/test/spies"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
)

func TestCreateTransaction(t *testing.T) {
	repoSpy := &spies.RepositorySpy{}
	hasherSpy := &spies.PasswordHasherSpy{}

	r := router.Setup(repoSpy, hasherSpy)

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
		transaction := &model.Transaction{
			Timestamp: time.Now().Round(0),
			Amount:    1000,
			Type:      model.Expense,
		}

		repoSpy.On("TransactionCreate", transaction.Timestamp, transaction.Amount, transaction.Type).Return(transaction, nil).Once()

		res := httptest.NewRecorder()
		req := newTransactionRequest(transaction)

		r.ServeHTTP(res, req)

		assertStatusCode(t, res, http.StatusCreated)
		assertSingleTransactionResponseBody(t, res, transaction)
	})
}

func TestGetTransaction(t *testing.T) {
	repoSpy := &spies.RepositorySpy{}
	hasherSpy := &spies.PasswordHasherSpy{}

	r := router.Setup(repoSpy, hasherSpy)

	newTransactionRequest := func(id uint) *http.Request {
		url := fmt.Sprintf("/transaction/%d", id)
		req, _ := http.NewRequest(http.MethodGet, url, nil)
		return req
	}

	t.Run("Get transaction with id = 0", func(t *testing.T) {
		id := uint(0)

		res := httptest.NewRecorder()
		req := newTransactionRequest(id)

		r.ServeHTTP(res, req)

		assertStatusCode(t, res, http.StatusBadRequest)
	})

	t.Run("Get transaction with non-existent id", func(t *testing.T) {
		id := uint(10)

		repoSpy.On("TransactionGet", id).Return(nil, repository.ErrorRecordNotFound).Once()

		res := httptest.NewRecorder()
		req := newTransactionRequest(id)

		r.ServeHTTP(res, req)

		assertStatusCode(t, res, http.StatusNotFound)
	})

	t.Run("Get transaction with valid id", func(t *testing.T) {
		id := uint(1)
		transaction := &model.Transaction{
			Timestamp: time.Now().Round(0),
			Amount:    1000,
			Type:      model.Expense,
		}

		repoSpy.On("TransactionGet", id).Return(transaction, nil).Once()

		res := httptest.NewRecorder()
		req := newTransactionRequest(id)

		r.ServeHTTP(res, req)

		assertStatusCode(t, res, http.StatusOK)
		assertSingleTransactionResponseBody(t, res, transaction)
	})
}

func TestUpdateTransaction(t *testing.T) {
	repoSpy := &spies.RepositorySpy{}
	hasherSpy := &spies.PasswordHasherSpy{}

	r := router.Setup(repoSpy, hasherSpy)

	newTransactionRequest := func(id uint, transaction *model.Transaction) *http.Request {
		url := fmt.Sprintf("/transaction/%d", id)
		body := createRequestBody(transaction)
		req, _ := http.NewRequest(http.MethodPatch, url, bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		return req
	}

	t.Run("Update non-existent transaction", func(t *testing.T) {
		id := uint(1)
		transaction := &model.Transaction{Amount: 1000}

		repoSpy.On("TransactionUpdate", id, mock.Anything, transaction.Amount, mock.Anything).Return(nil, repository.ErrorRecordNotFound).Once()

		res := httptest.NewRecorder()
		req := newTransactionRequest(id, transaction)

		r.ServeHTTP(res, req)

		assertStatusCode(t, res, http.StatusNotFound)
	})

	t.Run("Update existing transaction with invalid type", func(t *testing.T) {
		id := uint(2)
		transaction := &model.Transaction{Amount: 1000, Type: "invalid"}

		res := httptest.NewRecorder()
		req := newTransactionRequest(id, transaction)

		r.ServeHTTP(res, req)

		assertStatusCode(t, res, http.StatusBadRequest)
	})

	t.Run("Update existing transaction with valid arguments", func(t *testing.T) {
		id := uint(3)
		transaction := &model.Transaction{Amount: 2000, Type: model.Income}

		repoSpy.On("TransactionUpdate", id, mock.Anything, transaction.Amount, transaction.Type).Return(transaction, nil).Once()

		res := httptest.NewRecorder()
		req := newTransactionRequest(id, transaction)

		r.ServeHTTP(res, req)

		assertStatusCode(t, res, http.StatusOK)
		assertSingleTransactionResponseBody(t, res, transaction)
	})
}

func TestDeleteTransaction(t *testing.T) {
	repoSpy := &spies.RepositorySpy{}
	hasherSpy := &spies.PasswordHasherSpy{}

	r := router.Setup(repoSpy, hasherSpy)

	newTransactionRequest := func(id uint) *http.Request {
		url := fmt.Sprintf("/transaction/%d", id)
		req, _ := http.NewRequest(http.MethodDelete, url, nil)
		return req
	}

	t.Run("Delete non-existent transaction", func(t *testing.T) {
		id := uint(1)

		repoSpy.On("TransactionDelete", id).Return(repository.ErrorRecordNotFound).Once()

		res := httptest.NewRecorder()
		req := newTransactionRequest(id)

		r.ServeHTTP(res, req)

		assertStatusCode(t, res, http.StatusNotFound)
	})

	t.Run("Delete existing transaction", func(t *testing.T) {
		id := uint(2)

		repoSpy.On("TransactionDelete", id).Return(nil).Once()

		res := httptest.NewRecorder()
		req := newTransactionRequest(id)

		r.ServeHTTP(res, req)

		assertStatusCode(t, res, http.StatusNoContent)
	})
}

func TestListTransactions(t *testing.T) {
	repoSpy := &spies.RepositorySpy{}
	hasherSpy := &spies.PasswordHasherSpy{}

	r := router.Setup(repoSpy, hasherSpy)

	newTransactionListResponse := func(slice []*model.Transaction) *transactionListResponse {
		return &transactionListResponse{
			Count:   len(slice),
			Entries: slice,
		}
	}

	newTransactionRequest := func() *http.Request {
		req, _ := http.NewRequest(http.MethodGet, "/transaction/", nil)
		return req
	}

	t.Run("List transactions when there are no transactions", func(t *testing.T) {
		transactions := []*model.Transaction{}

		repoSpy.On("TransactionList").Return(transactions, nil).Once()

		res := httptest.NewRecorder()
		req := newTransactionRequest()

		r.ServeHTTP(res, req)

		expected := newTransactionListResponse(transactions)

		assertStatusCode(t, res, http.StatusOK)
		assertListTransactionResponseBody(t, res, expected)
	})

	t.Run("List transactions when there are non-zero transactions", func(t *testing.T) {
		transactions := []*model.Transaction{{}, {}}

		repoSpy.On("TransactionList").Return(transactions, nil).Once()

		res := httptest.NewRecorder()
		req := newTransactionRequest()

		r.ServeHTTP(res, req)

		expected := newTransactionListResponse(transactions)

		assertStatusCode(t, res, http.StatusOK)
		assertListTransactionResponseBody(t, res, expected)
	})
}

type transactionListResponse struct {
	Count   int                  `json:"count"`
	Entries []*model.Transaction `json:"entries"`
}

func assertSingleTransactionResponseBody(t *testing.T, res *httptest.ResponseRecorder, transaction *model.Transaction) {
	t.Helper()

	var got model.Transaction
	if err := json.NewDecoder(res.Body).Decode(&got); err != nil {
		t.Errorf("couldn't parse json response: %v", err)
	}

	if !reflect.DeepEqual(got, *transaction) {
		t.Errorf("expected %+v, got %+v", *transaction, got)
	}
}

func assertListTransactionResponseBody(t *testing.T, res *httptest.ResponseRecorder, expected *transactionListResponse) {
	t.Helper()

	var got transactionListResponse
	if err := json.NewDecoder(res.Body).Decode(&got); err != nil {
		t.Errorf("couldn't parse json response: %v", err)
	}

	if !reflect.DeepEqual(got, *expected) {
		t.Errorf("expected %+v ;%T, got %+v ;%T", *expected, *expected, got, got)
	}
}
