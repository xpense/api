package test

import (
	"bytes"
	"encoding/json"
	"expense-api/handlers"
	"expense-api/middleware/auth"
	"expense-api/model"
	"expense-api/repository"
	"expense-api/router"
	"expense-api/router/test/spies"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
)

const basePartiesPath = "/parties/"

func TestCreateParty(t *testing.T) {
	repoSpy := &spies.RepositorySpy{}
	jwtServiceSpy := &spies.JWTServiceSpy{}
	hasherSpy := &spies.PasswordHasherSpy{}

	r := router.Setup(repoSpy, jwtServiceSpy, hasherSpy, router.TestConfig)

	newPartyRequest := func(party *model.Party, token string) *http.Request {
		body := createRequestBody(party)
		req, _ := http.NewRequest(http.MethodPost, basePartiesPath, bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		return req
	}

	t.Run("Missing/Invalid authorization token cases", func(t *testing.T) {
		party := &model.Party{}
		token := "invalid-token"

		missingTokenReq := newPartyRequest(party, token)
		invalidTokenReq := newPartyRequest(party, token)

		unauthorizedTestCases := UnauthorizedTestCases(missingTokenReq, invalidTokenReq, r, jwtServiceSpy)
		t.Run("Unauthorized test cases", unauthorizedTestCases)
	})

	t.Run("Valid authorization token cases", func(t *testing.T) {
		token := "valid-token"
		userID := uint(1)
		claims := auth.CustomClaims{
			ID: userID,
		}
		jwtServiceSpy.On("ValidateJWT", token).Return(&claims, nil)

		t.Run("Try to create a party with already existing name, belonging to the same user", func(t *testing.T) {
			party := &model.Party{
				Name:   "Amazon",
				UserID: userID,
			}

			repoSpy.On("PartyCreate", party).Return(repository.ErrorUniqueConstaintViolation).Once()

			res := httptest.NewRecorder()
			req := newPartyRequest(party, token)

			r.ServeHTTP(res, req)

			wantErrorMessage := handlers.ErrMsgPartyNameTaken

			assertStatusCode(t, res, http.StatusConflict)
			assertErrorMessage(t, res, wantErrorMessage)
		})

		t.Run("Create party with valid data", func(t *testing.T) {
			party := &model.Party{
				Name:   "Rewe",
				UserID: userID,
			}

			repoSpy.On("PartyCreate", party).Return(nil).Once()

			res := httptest.NewRecorder()
			req := newPartyRequest(party, token)

			r.ServeHTTP(res, req)

			party.UserID = 0

			assertStatusCode(t, res, http.StatusCreated)
			assertSinglePartyResponseBody(t, res, party)
		})
	})
}

func TestGetParty(t *testing.T) {
	repoSpy := &spies.RepositorySpy{}
	jwtServiceSpy := &spies.JWTServiceSpy{}
	hasherSpy := &spies.PasswordHasherSpy{}

	r := router.Setup(repoSpy, jwtServiceSpy, hasherSpy, router.TestConfig)

	newPartyRequest := func(id uint, token string) *http.Request {
		url := fmt.Sprintf("%s%d", basePartiesPath, id)
		req, _ := http.NewRequest(http.MethodGet, url, nil)
		req.Header.Set("Authorization", "Bearer "+token)
		return req
	}

	t.Run("Missing/Invalid authorization token cases", func(t *testing.T) {
		id := uint(1)
		token := "invalid-token"

		missingTokenReq := newPartyRequest(id, token)
		invalidTokenReq := newPartyRequest(id, token)

		unauthorizedTestCases := UnauthorizedTestCases(missingTokenReq, invalidTokenReq, r, jwtServiceSpy)
		t.Run("Unauthorized test cases", unauthorizedTestCases)
	})

	t.Run("Valid authorization token cases", func(t *testing.T) {
		token := "valid-token"
		userID := uint(1)
		claims := auth.CustomClaims{
			ID: userID,
		}
		jwtServiceSpy.On("ValidateJWT", token).Return(&claims, nil)

		t.Run("Get party with id = 0", func(t *testing.T) {
			id := uint(0)

			res := httptest.NewRecorder()
			req := newPartyRequest(id, token)

			r.ServeHTTP(res, req)

			assertStatusCode(t, res, http.StatusBadRequest)
		})

		t.Run("Get party with non-existent id", func(t *testing.T) {
			id := uint(10)

			repoSpy.On("PartyGet", id).Return(nil, repository.ErrorRecordNotFound).Once()

			res := httptest.NewRecorder()
			req := newPartyRequest(id, token)

			r.ServeHTTP(res, req)

			assertStatusCode(t, res, http.StatusNotFound)
		})

		t.Run("Get party with valid id that belongs to another user", func(t *testing.T) {
			id := uint(1)
			anotherUsersID := uint(6)
			party := &model.Party{
				Name:   "new party",
				UserID: anotherUsersID,
			}

			repoSpy.On("PartyGet", id).Return(party, nil).Once()

			res := httptest.NewRecorder()
			req := newPartyRequest(id, token)

			r.ServeHTTP(res, req)

			assertStatusCode(t, res, http.StatusForbidden)
		})

		t.Run("Get party with valid id", func(t *testing.T) {
			id := uint(1)
			party := &model.Party{
				Name:   "new party",
				UserID: userID,
			}

			repoSpy.On("PartyGet", id).Return(party, nil).Twice()

			res := httptest.NewRecorder()
			req := newPartyRequest(id, token)

			r.ServeHTTP(res, req)

			party.UserID = 0

			assertStatusCode(t, res, http.StatusOK)
			assertSinglePartyResponseBody(t, res, party)
		})
	})
}

func TestUpdateParty(t *testing.T) {
	repoSpy := &spies.RepositorySpy{}
	jwtServiceSpy := &spies.JWTServiceSpy{}
	hasherSpy := &spies.PasswordHasherSpy{}

	r := router.Setup(repoSpy, jwtServiceSpy, hasherSpy, router.TestConfig)

	newPartyRequest := func(id uint, party *model.Party, token string) *http.Request {
		url := fmt.Sprintf("%s%d", basePartiesPath, id)
		body := createRequestBody(party)
		req, _ := http.NewRequest(http.MethodPatch, url, bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		return req
	}

	t.Run("Missing/Invalid authorization token cases", func(t *testing.T) {
		id := uint(1)
		party := &model.Party{}
		token := "invalid-token"

		missingTokenReq := newPartyRequest(id, party, token)
		invalidTokenReq := newPartyRequest(id, party, token)

		unauthorizedTestCases := UnauthorizedTestCases(missingTokenReq, invalidTokenReq, r, jwtServiceSpy)
		t.Run("Unauthorized test cases", unauthorizedTestCases)
	})

	t.Run("Valid authorization token cases", func(t *testing.T) {
		token := "valid-token"
		userID := uint(1)
		claims := auth.CustomClaims{
			ID: userID,
		}
		jwtServiceSpy.On("ValidateJWT", token).Return(&claims, nil)

		t.Run("Update non-existent party", func(t *testing.T) {
			id := uint(1)
			party := &model.Party{
				Name:   "new party",
				UserID: userID,
			}

			repoSpy.On("PartyGet", id).Return(nil, repository.ErrorRecordNotFound).Once()

			res := httptest.NewRecorder()
			req := newPartyRequest(id, party, token)

			r.ServeHTTP(res, req)

			assertStatusCode(t, res, http.StatusNotFound)
		})

		t.Run("Try to update party with valid id that belongs to another user", func(t *testing.T) {
			id := uint(1)
			anotherUsersID := uint(6)
			party := &model.Party{
				Name:   "new party",
				UserID: anotherUsersID,
			}

			repoSpy.On("PartyGet", id).Return(party, nil).Once()

			res := httptest.NewRecorder()
			req := newPartyRequest(id, party, token)

			r.ServeHTTP(res, req)

			assertStatusCode(t, res, http.StatusForbidden)
		})

		t.Run("Try to update a party with already existing name, belonging to the same user", func(t *testing.T) {
			id := uint(1)
			party := &model.Party{
				Name:   "Kaufland",
				UserID: userID,
			}

			repoSpy.On("PartyGet", id).Return(party, nil).Once()
			repoSpy.On("PartyUpdate", id, party).Return(nil, repository.ErrorUniqueConstaintViolation).Once()

			res := httptest.NewRecorder()
			req := newPartyRequest(id, party, token)

			r.ServeHTTP(res, req)

			wantErrorMessage := handlers.ErrMsgPartyNameTaken

			assertStatusCode(t, res, http.StatusConflict)
			assertErrorMessage(t, res, wantErrorMessage)
		})

		t.Run("Update existing party with valid arguments", func(t *testing.T) {
			id := uint(3)
			party := &model.Party{
				Name:   "new party",
				UserID: userID,
			}

			repoSpy.On("PartyGet", id).Return(party, nil).Once()
			repoSpy.On("PartyUpdate", id, party).Return(party, nil).Once()

			res := httptest.NewRecorder()
			req := newPartyRequest(id, party, token)

			r.ServeHTTP(res, req)

			party.UserID = 0

			assertStatusCode(t, res, http.StatusOK)
			assertSinglePartyResponseBody(t, res, party)
		})
	})
}

func TestDeleteParty(t *testing.T) {
	repoSpy := &spies.RepositorySpy{}
	jwtServiceSpy := &spies.JWTServiceSpy{}
	hasherSpy := &spies.PasswordHasherSpy{}

	r := router.Setup(repoSpy, jwtServiceSpy, hasherSpy, router.TestConfig)

	newPartyRequest := func(id uint, token string) *http.Request {
		url := fmt.Sprintf("%s%d", basePartiesPath, id)
		req, _ := http.NewRequest(http.MethodDelete, url, nil)
		req.Header.Set("Authorization", "Bearer "+token)
		return req
	}

	t.Run("Missing/Invalid authorization token cases", func(t *testing.T) {
		id := uint(1)
		token := "invalid-token"

		missingTokenReq := newPartyRequest(id, token)
		invalidTokenReq := newPartyRequest(id, token)

		unauthorizedTestCases := UnauthorizedTestCases(missingTokenReq, invalidTokenReq, r, jwtServiceSpy)
		t.Run("Unauthorized test cases", unauthorizedTestCases)
	})

	t.Run("Valid authorization token cases", func(t *testing.T) {
		token := "valid-token"
		userID := uint(1)
		claims := auth.CustomClaims{
			ID: userID,
		}
		jwtServiceSpy.On("ValidateJWT", token).Return(&claims, nil)

		t.Run("Delete non-existent party", func(t *testing.T) {
			id := uint(1)

			repoSpy.On("PartyGet", id).Return(nil, repository.ErrorRecordNotFound).Once()

			res := httptest.NewRecorder()
			req := newPartyRequest(id, token)

			r.ServeHTTP(res, req)

			assertStatusCode(t, res, http.StatusNotFound)
		})

		t.Run("Try to delete party with valid id that belongs to another user", func(t *testing.T) {
			id := uint(1)
			anotherUsersID := uint(6)
			party := &model.Party{
				Name:   "new party",
				UserID: anotherUsersID,
			}

			repoSpy.On("PartyGet", id).Return(party, nil).Once()

			res := httptest.NewRecorder()
			req := newPartyRequest(id, token)

			r.ServeHTTP(res, req)

			assertStatusCode(t, res, http.StatusForbidden)
		})

		t.Run("Delete existing party", func(t *testing.T) {
			id := uint(2)
			party := &model.Party{
				Name:   "new party",
				UserID: userID,
			}

			repoSpy.On("PartyGet", id).Return(party, nil).Once()
			repoSpy.On("PartyDelete", id).Return(nil).Once()

			res := httptest.NewRecorder()
			req := newPartyRequest(id, token)

			r.ServeHTTP(res, req)

			assertStatusCode(t, res, http.StatusNoContent)
		})
	})
}

func TestListParties(t *testing.T) {
	repoSpy := &spies.RepositorySpy{}
	jwtServiceSpy := &spies.JWTServiceSpy{}
	hasherSpy := &spies.PasswordHasherSpy{}

	r := router.Setup(repoSpy, jwtServiceSpy, hasherSpy, router.TestConfig)

	newPartyListResponse := func(slice []*model.Party) *partyListResponse {
		return &partyListResponse{
			Count:   len(slice),
			Entries: slice,
		}
	}

	newPartyRequest := func(token string) *http.Request {
		req, _ := http.NewRequest(http.MethodGet, basePartiesPath, nil)
		req.Header.Set("Authorization", "Bearer "+token)
		return req
	}

	t.Run("Missing/Invalid authorization token cases", func(t *testing.T) {
		token := "invalid-token"

		missingTokenReq := newPartyRequest(token)
		invalidTokenReq := newPartyRequest(token)

		unauthorizedTestCases := UnauthorizedTestCases(missingTokenReq, invalidTokenReq, r, jwtServiceSpy)
		t.Run("Unauthorized test cases", unauthorizedTestCases)
	})

	t.Run("Valid authorization token cases", func(t *testing.T) {
		token := "valid-token"
		userID := uint(1)
		claims := auth.CustomClaims{
			ID: userID,
		}
		jwtServiceSpy.On("ValidateJWT", token).Return(&claims, nil)

		t.Run("List parties when there are no parties", func(t *testing.T) {
			parties := []*model.Party{}
			repoSpy.On("PartyList", userID).Return(parties, nil).Once()

			res := httptest.NewRecorder()
			req := newPartyRequest(token)

			r.ServeHTTP(res, req)

			expected := newPartyListResponse(parties)

			assertStatusCode(t, res, http.StatusOK)
			assertListPartyResponseBody(t, res, expected)
		})

		t.Run("List parties when there are non-zero parties", func(t *testing.T) {
			parties := []*model.Party{{}}

			repoSpy.On("PartyList", userID).Return(parties, nil).Once()

			res := httptest.NewRecorder()
			req := newPartyRequest(token)

			r.ServeHTTP(res, req)

			expected := newPartyListResponse(parties)

			assertStatusCode(t, res, http.StatusOK)
			assertListPartyResponseBody(t, res, expected)
		})
	})
}

func TestListTransactionsByParty(t *testing.T) {
	repoSpy := &spies.RepositorySpy{}
	jwtServiceSpy := &spies.JWTServiceSpy{}
	hasherSpy := &spies.PasswordHasherSpy{}

	r := router.Setup(repoSpy, jwtServiceSpy, hasherSpy, router.TestConfig)

	newTransactionListResponse := func(slice []*model.Transaction) *transactionListResponse {
		return &transactionListResponse{
			Count:   len(slice),
			Entries: slice,
		}
	}

	newPartyRequest := func(id uint, token string) *http.Request {
		url := fmt.Sprintf("%s%d/transactions", basePartiesPath, id)
		req, _ := http.NewRequest(http.MethodGet, url, nil)
		req.Header.Set("Authorization", "Bearer "+token)
		return req
	}

	t.Run("Missing/Invalid authorization token cases", func(t *testing.T) {
		id := uint(1)
		token := "invalid-token"

		missingTokenReq := newPartyRequest(id, token)
		invalidTokenReq := newPartyRequest(id, token)

		unauthorizedTestCases := UnauthorizedTestCases(missingTokenReq, invalidTokenReq, r, jwtServiceSpy)
		t.Run("Unauthorized test cases", unauthorizedTestCases)
	})

	t.Run("Valid authorization token cases", func(t *testing.T) {
		id := uint(1)
		token := "valid-token"
		userID := uint(1)
		claims := auth.CustomClaims{
			ID: userID,
		}
		jwtServiceSpy.On("ValidateJWT", token).Return(&claims, nil)

		t.Run("List transactions of a non-existent party", func(t *testing.T) {
			repoSpy.On("PartyGet", id).Return(nil, repository.ErrorRecordNotFound).Once()

			res := httptest.NewRecorder()
			req := newPartyRequest(id, token)

			r.ServeHTTP(res, req)

			assertStatusCode(t, res, http.StatusNotFound)
		})

		t.Run("List transactions of a party that belongs to another user", func(t *testing.T) {
			party := &model.Party{
				UserID: userID + 1,
			}
			repoSpy.On("PartyGet", id).Return(party, nil).Once()

			res := httptest.NewRecorder()
			req := newPartyRequest(id, token)

			r.ServeHTTP(res, req)

			assertStatusCode(t, res, http.StatusForbidden)
		})

		t.Run("List transactions when there are no transactions", func(t *testing.T) {
			party := &model.Party{
				UserID: userID,
			}
			transactions := []*model.Transaction{}

			repoSpy.On("PartyGet", id).Return(party, nil).Once()
			repoSpy.On("TransactionListByParty", userID, id).Return(transactions, nil).Once()

			res := httptest.NewRecorder()
			req := newPartyRequest(id, token)

			r.ServeHTTP(res, req)

			expected := newTransactionListResponse(transactions)

			assertStatusCode(t, res, http.StatusOK)
			assertListTransactionResponseBody(t, res, expected)
		})

		t.Run("List transactions when there are non-zero transactions", func(t *testing.T) {
			party := &model.Party{
				UserID: userID,
			}
			transactions := []*model.Transaction{{}}

			repoSpy.On("PartyGet", id).Return(party, nil).Once()
			repoSpy.On("TransactionListByParty", userID, id).Return(transactions, nil).Once()

			res := httptest.NewRecorder()
			req := newPartyRequest(id, token)

			r.ServeHTTP(res, req)

			expected := newTransactionListResponse(transactions)

			assertStatusCode(t, res, http.StatusOK)
			assertListTransactionResponseBody(t, res, expected)
		})
	})
}

type partyListResponse struct {
	Count   int            `json:"count"`
	Entries []*model.Party `json:"entries"`
}

func assertSinglePartyResponseBody(t *testing.T, res *httptest.ResponseRecorder, party *model.Party) {
	t.Helper()

	var got model.Party
	if err := json.NewDecoder(res.Body).Decode(&got); err != nil {
		t.Errorf("couldn't parse json response: %v", err)
	}

	if !cmp.Equal(got, *party) {
		t.Errorf("expected %+v, got %+v", *party, got)
	}
}

func assertListPartyResponseBody(t *testing.T, res *httptest.ResponseRecorder, expected *partyListResponse) {
	t.Helper()

	var got partyListResponse
	if err := json.NewDecoder(res.Body).Decode(&got); err != nil {
		t.Errorf("couldn't parse json response: %v", err)
	}

	if !cmp.Equal(got, *expected) {
		t.Errorf("expected %+v, got %+v", *expected, got)
	}
}
