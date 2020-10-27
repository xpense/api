package router

import (
	"bytes"
	"encoding/json"
	"expense-api/model"
	"expense-api/repository"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"sort"
	"testing"
	"time"
)

type RepositorySpy struct {
	transactions map[uint]*model.Transaction
}

func NewRepositorySpy() *RepositorySpy {
	return &RepositorySpy{
		transactions: map[uint]*model.Transaction{},
	}
}

func (r *RepositorySpy) transactionSlice() []*model.Transaction {
	res := make([]*model.Transaction, 0, len(r.transactions))

	for _, t := range r.transactions {
		res = append(res, t)
	}

	sort.Slice(res, func(i, j int) bool {
		return res[i].ID < res[j].ID
	})

	return res
}

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

func (r *RepositorySpy) TransactionCreate(timestamp time.Time, amount uint64, transactionType model.TransactionType) (*model.Transaction, error) {
	id := uint(len(r.transactions)) + 1

	t := &model.Transaction{}

	t.ID = id
	if !timestamp.IsZero() {
		t.Timestamp = timestamp.Round(0)
	}

	if amount > 0 {
		t.Amount = amount
	}

	if transactionType == model.Income || t.Type == model.Expense {
		t.Type = transactionType
	}

	r.transactions[id] = t

	return t, nil
}

func (r *RepositorySpy) TransactionUpdate(id uint, timestamp time.Time, amount uint64, transactionType model.TransactionType) (*model.Transaction, error) {
	t, err := r.TransactionGet(id)
	if err != nil {
		return nil, err
	}

	if !timestamp.IsZero() {
		t.Timestamp = timestamp
	}

	if amount > 0 {
		t.Amount = amount
	}

	if transactionType == model.Income || t.Type == model.Expense {
		t.Type = transactionType
	}

	return t, nil
}

func (r *RepositorySpy) TransactionGet(id uint) (*model.Transaction, error) {
	t, ok := r.transactions[id]
	if !ok {
		return nil, repository.ErrorRecordNotFound
	}

	return t, nil
}

func (r *RepositorySpy) TransactionDelete(id uint) error {
	if _, err := r.TransactionGet(id); err != nil {
		return err
	}

	delete(r.transactions, id)
	return nil
}

func (r *RepositorySpy) TransactionList() ([]*model.Transaction, error) {
	transactions := make([]*model.Transaction, 0, len(r.transactions))

	for _, t := range r.transactions {
		transactions = append(transactions, t)
	}

	return transactions, nil
}

func TestCreateTransaction(t *testing.T) {
	spy := NewRepositorySpy()
	r := Setup(spy)

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
	r := Setup(spy)

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

func TestListTransactions(t *testing.T) {
	spy := NewRepositorySpy()
	r := Setup(spy)

	// firstTransaction, _ := spy.TransactionCreate(time.Now(), 1000, model.Expense)

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
