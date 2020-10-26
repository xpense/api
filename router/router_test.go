package router

import (
	"bytes"
	"encoding/json"
	"expense-api/model"
	"expense-api/repository"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type RepositorySpy struct {
	transactions map[uint]*model.Transaction
}

func NewRepositorySpy() repository.Repository {
	return &RepositorySpy{
		transactions: map[uint]*model.Transaction{},
	}
}

func (r *RepositorySpy) TransactionCreate(timestamp time.Time, amount uint64, transactionType model.TransactionType) (*model.Transaction, error) {
	id := uint(len(r.transactions))

	t := &model.Transaction{}

	t.ID = id
	if !timestamp.IsZero() {
		t.Timestamp = timestamp
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

	newCreateTransactionRequest := func(transaction *model.Transaction) *http.Request {
		body := createRequestBody(transaction)
		req, _ := http.NewRequest(http.MethodPost, "/transaction/", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		return req
	}

	t.Run("Create transaction with amount = 0", func(t *testing.T) {
		res := httptest.NewRecorder()
		req := newCreateTransactionRequest(&model.Transaction{Amount: 0})

		r.ServeHTTP(res, req)

		assertStatusCode(t, res, http.StatusBadRequest)
	})

	t.Run("Create transaction with invalid transaction type", func(t *testing.T) {
		res := httptest.NewRecorder()
		req := newCreateTransactionRequest(&model.Transaction{Amount: 1000, Type: "invalid"})

		r.ServeHTTP(res, req)

		assertStatusCode(t, res, http.StatusBadRequest)
	})

	t.Run("Create transaction with valid data", func(t *testing.T) {
		res := httptest.NewRecorder()
		req := newCreateTransactionRequest(&model.Transaction{Amount: 1000, Type: model.Expense})

		r.ServeHTTP(res, req)

		assertStatusCode(t, res, http.StatusCreated)

		var jsonResponse model.Transaction
		if err := json.NewDecoder(res.Body).Decode(&jsonResponse); err != nil {
			t.Errorf("couldn't parse json response")
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
