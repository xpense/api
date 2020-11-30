package router

import (
	"expense-api/internal/handlers"
	"expense-api/internal/middleware/auth"
	"expense-api/internal/model"
	"expense-api/internal/repository"
	"expense-api/internal/router"
	"expense-api/test/spies"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateParty(t *testing.T) {
	repoSpy := &spies.RepositorySpy{}
	jwtServiceSpy := &spies.JWTServiceSpy{}
	hasherSpy := &spies.PasswordHasherSpy{}

	r := router.Setup(repoSpy, jwtServiceSpy, hasherSpy, router.TestConfig)

	t.Run("Missing/Invalid authorization token cases", func(t *testing.T) {
		party := &handlers.Party{}
		token := "invalid-token"

		missingTokenReq := NewCreatePartyRequest(party, token)
		invalidTokenReq := NewCreatePartyRequest(party, token)

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
			req := NewCreatePartyRequest(&handlers.Party{Name: party.Name}, token)

			r.ServeHTTP(res, req)

			wantErrorMessage := handlers.ErrorPartyNameTaken.Error()

			AssertStatusCode(t, res, http.StatusConflict)
			AssertErrorMessage(t, res, wantErrorMessage)
		})

		t.Run("Create party with valid data", func(t *testing.T) {
			party := &model.Party{
				Name:   "Rewe",
				UserID: userID,
			}

			repoSpy.On("PartyCreate", party).Return(nil).Once()

			res := httptest.NewRecorder()
			req := NewCreatePartyRequest(&handlers.Party{Name: party.Name}, token)

			r.ServeHTTP(res, req)

			resBody := handlers.PartyModelToResponse(party)

			AssertStatusCode(t, res, http.StatusCreated)
			AssertResponseBody(t, res, resBody)
		})
	})
}

func TestGetParty(t *testing.T) {
	repoSpy := &spies.RepositorySpy{}
	jwtServiceSpy := &spies.JWTServiceSpy{}
	hasherSpy := &spies.PasswordHasherSpy{}

	r := router.Setup(repoSpy, jwtServiceSpy, hasherSpy, router.TestConfig)

	t.Run("Missing/Invalid authorization token cases", func(t *testing.T) {
		id := uint(1)
		token := "invalid-token"

		missingTokenReq := NewGetPartyRequest(id, token)
		invalidTokenReq := NewGetPartyRequest(id, token)

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
			req := NewGetPartyRequest(id, token)

			r.ServeHTTP(res, req)

			AssertStatusCode(t, res, http.StatusBadRequest)
		})

		t.Run("Get party with non-existent id", func(t *testing.T) {
			id := uint(10)

			repoSpy.On("PartyGet", id).Return(nil, repository.ErrorRecordNotFound).Once()

			res := httptest.NewRecorder()
			req := NewGetPartyRequest(id, token)

			r.ServeHTTP(res, req)

			AssertStatusCode(t, res, http.StatusNotFound)
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
			req := NewGetPartyRequest(id, token)

			r.ServeHTTP(res, req)

			AssertStatusCode(t, res, http.StatusForbidden)
		})

		t.Run("Get party with valid id", func(t *testing.T) {
			id := uint(1)
			party := &model.Party{
				Name:   "new party",
				UserID: userID,
			}

			repoSpy.On("PartyGet", id).Return(party, nil).Twice()

			res := httptest.NewRecorder()
			req := NewGetPartyRequest(id, token)

			r.ServeHTTP(res, req)

			resBody := handlers.PartyModelToResponse(party)

			AssertStatusCode(t, res, http.StatusOK)
			AssertResponseBody(t, res, resBody)
		})
	})
}

func TestUpdateParty(t *testing.T) {
	repoSpy := &spies.RepositorySpy{}
	jwtServiceSpy := &spies.JWTServiceSpy{}
	hasherSpy := &spies.PasswordHasherSpy{}

	r := router.Setup(repoSpy, jwtServiceSpy, hasherSpy, router.TestConfig)

	t.Run("Missing/Invalid authorization token cases", func(t *testing.T) {
		id := uint(1)
		party := &handlers.Party{}
		token := "invalid-token"

		missingTokenReq := NewUpdatePartyRequest(id, party, token)
		invalidTokenReq := NewUpdatePartyRequest(id, party, token)

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
			party := &handlers.Party{
				Name: "new party",
			}

			repoSpy.On("PartyGet", id).Return(nil, repository.ErrorRecordNotFound).Once()

			res := httptest.NewRecorder()
			req := NewUpdatePartyRequest(id, party, token)

			r.ServeHTTP(res, req)

			AssertStatusCode(t, res, http.StatusNotFound)
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
			req := NewUpdatePartyRequest(id, &handlers.Party{Name: party.Name}, token)

			r.ServeHTTP(res, req)

			AssertStatusCode(t, res, http.StatusForbidden)
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
			req := NewUpdatePartyRequest(id, &handlers.Party{Name: party.Name}, token)

			r.ServeHTTP(res, req)

			wantErrorMessage := handlers.ErrorPartyNameTaken.Error()

			AssertStatusCode(t, res, http.StatusConflict)
			AssertErrorMessage(t, res, wantErrorMessage)
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
			req := NewUpdatePartyRequest(id, &handlers.Party{Name: party.Name}, token)

			r.ServeHTTP(res, req)

			resBody := handlers.PartyModelToResponse(party)

			AssertStatusCode(t, res, http.StatusOK)
			AssertResponseBody(t, res, resBody)
		})
	})
}

func TestDeleteParty(t *testing.T) {
	repoSpy := &spies.RepositorySpy{}
	jwtServiceSpy := &spies.JWTServiceSpy{}
	hasherSpy := &spies.PasswordHasherSpy{}

	r := router.Setup(repoSpy, jwtServiceSpy, hasherSpy, router.TestConfig)

	t.Run("Missing/Invalid authorization token cases", func(t *testing.T) {
		id := uint(1)
		token := "invalid-token"

		missingTokenReq := NewDeletePartyRequest(id, token)
		invalidTokenReq := NewDeletePartyRequest(id, token)

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
			req := NewDeletePartyRequest(id, token)

			r.ServeHTTP(res, req)

			AssertStatusCode(t, res, http.StatusNotFound)
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
			req := NewDeletePartyRequest(id, token)

			r.ServeHTTP(res, req)

			AssertStatusCode(t, res, http.StatusForbidden)
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
			req := NewDeletePartyRequest(id, token)

			r.ServeHTTP(res, req)

			AssertStatusCode(t, res, http.StatusNoContent)
		})
	})
}

func TestListParties(t *testing.T) {
	repoSpy := &spies.RepositorySpy{}
	jwtServiceSpy := &spies.JWTServiceSpy{}
	hasherSpy := &spies.PasswordHasherSpy{}

	r := router.Setup(repoSpy, jwtServiceSpy, hasherSpy, router.TestConfig)

	newPartyListResponse := func(slice []*handlers.Party) *PartyListResponse {
		return &PartyListResponse{
			Count:   len(slice),
			Entries: slice,
		}
	}

	t.Run("Missing/Invalid authorization token cases", func(t *testing.T) {
		token := "invalid-token"

		missingTokenReq := NewListPartiesRequest(token)
		invalidTokenReq := NewListPartiesRequest(token)

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
			req := NewListPartiesRequest(token)

			r.ServeHTTP(res, req)

			expected := newPartyListResponse([]*handlers.Party{})

			AssertStatusCode(t, res, http.StatusOK)
			AssertResponseBody(t, res, expected)
		})

		t.Run("List parties when there are non-zero parties", func(t *testing.T) {
			parties := []*model.Party{{}}

			repoSpy.On("PartyList", userID).Return(parties, nil).Once()

			res := httptest.NewRecorder()
			req := NewListPartiesRequest(token)

			r.ServeHTTP(res, req)

			expected := newPartyListResponse([]*handlers.Party{{}})

			AssertStatusCode(t, res, http.StatusOK)
			AssertResponseBody(t, res, expected)
		})
	})
}

func TestListTransactionsByParty(t *testing.T) {
	repoSpy := &spies.RepositorySpy{}
	jwtServiceSpy := &spies.JWTServiceSpy{}
	hasherSpy := &spies.PasswordHasherSpy{}

	r := router.Setup(repoSpy, jwtServiceSpy, hasherSpy, router.TestConfig)

	newTransactionListResponse := func(slice []*handlers.Transaction) *TransactionListResponse {
		return &TransactionListResponse{
			Count:   len(slice),
			Entries: slice,
		}
	}

	t.Run("Missing/Invalid authorization token cases", func(t *testing.T) {
		id := uint(1)
		token := "invalid-token"

		missingTokenReq := NewListTransactionsByPartyRequest(id, token)
		invalidTokenReq := NewListTransactionsByPartyRequest(id, token)

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
			req := NewListTransactionsByPartyRequest(id, token)

			r.ServeHTTP(res, req)

			AssertStatusCode(t, res, http.StatusNotFound)
		})

		t.Run("List transactions of a party that belongs to another user", func(t *testing.T) {
			party := &model.Party{
				UserID: userID + 1,
			}
			repoSpy.On("PartyGet", id).Return(party, nil).Once()

			res := httptest.NewRecorder()
			req := NewListTransactionsByPartyRequest(id, token)

			r.ServeHTTP(res, req)

			AssertStatusCode(t, res, http.StatusForbidden)
		})

		t.Run("List transactions when there are no transactions", func(t *testing.T) {
			party := &model.Party{
				UserID: userID,
			}
			transactions := []*model.Transaction{}

			repoSpy.On("PartyGet", id).Return(party, nil).Once()
			repoSpy.On("TransactionListByParty", userID, id).Return(transactions, nil).Once()

			res := httptest.NewRecorder()
			req := NewListTransactionsByPartyRequest(id, token)

			r.ServeHTTP(res, req)

			expected := newTransactionListResponse([]*handlers.Transaction{})

			AssertStatusCode(t, res, http.StatusOK)
			AssertResponseBody(t, res, expected)
		})

		t.Run("List transactions when there are non-zero transactions", func(t *testing.T) {
			party := &model.Party{
				UserID: userID,
			}
			transactions := []*model.Transaction{{}}

			repoSpy.On("PartyGet", id).Return(party, nil).Once()
			repoSpy.On("TransactionListByParty", userID, id).Return(transactions, nil).Once()

			res := httptest.NewRecorder()
			req := NewListTransactionsByPartyRequest(id, token)

			r.ServeHTTP(res, req)

			expected := newTransactionListResponse([]*handlers.Transaction{{}})

			AssertStatusCode(t, res, http.StatusOK)
			AssertResponseBody(t, res, expected)
		})
	})
}
