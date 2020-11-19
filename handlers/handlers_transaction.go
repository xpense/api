package handlers

import (
	"expense-api/middleware"
	auth_middleware "expense-api/middleware/auth"
	"expense-api/model"
	"expense-api/repository"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type TransactionsHandler interface {
	ListTransactions(ctx *gin.Context)
	CreateTransaction(ctx *gin.Context)
	GetTransaction(ctx *gin.Context)
	UpdateTransaction(ctx *gin.Context)
	DeleteTransaction(ctx *gin.Context)
}

// Transaction is a transaction with an omitted user
type Transaction struct {
	ID          uint      `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Timestamp   time.Time `json:"timestamp"`
	Amount      uint64    `json:"amount"`
	Description string    `json:"description"`
}

func TransactionModelToResponse(t *model.Transaction) *Transaction {
	return &Transaction{
		ID:          t.ID,
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
		Timestamp:   t.Timestamp,
		Amount:      t.Amount,
		Description: t.Description,
	}
}

func (h *handler) CreateTransaction(ctx *gin.Context) {
	userID, err := auth_middleware.GetUserIDFromContext(ctx)
	if err != nil {
		ctx.Status(http.StatusUnauthorized)
		return
	}

	var tRequest model.Transaction
	if err := ctx.Bind(&tRequest); err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

	tModel, err := TransactionCreateRequestToModel(&tRequest)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	tModel.UserID = userID

	if err := h.repo.TransactionCreate(tModel); err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	tResponse := TransactionModelToResponse(tModel)

	ctx.JSON(http.StatusCreated, tResponse)
}

func (h *handler) UpdateTransaction(ctx *gin.Context) {
	id := middleware.GetIDParamFromContext(ctx)

	var tRequest model.Transaction
	if err := ctx.Bind(&tRequest); err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

	updatedTModel, err := h.repo.TransactionUpdate(id, &tRequest)
	if err != nil {
		if err == repository.ErrorRecordNotFound {
			ctx.Status(http.StatusNotFound)
			return
		}
		ctx.Status(http.StatusInternalServerError)
		return
	}

	tResponse := TransactionModelToResponse(updatedTModel)

	ctx.JSON(http.StatusOK, tResponse)
}

func (h *handler) DeleteTransaction(ctx *gin.Context) {
	idStr := ctx.Param("id")
	tID, err := strconv.Atoi(idStr)
	if err != nil || tID <= 0 {
		ctx.Status(http.StatusBadRequest)
		return
	}

	if err := h.repo.TransactionDelete(uint(tID)); err != nil {
		if err == repository.ErrorRecordNotFound {
			ctx.Status(http.StatusNotFound)
			return
		}
		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (h *handler) GetTransaction(ctx *gin.Context) {
	userID, err := auth_middleware.GetUserIDFromContext(ctx)
	if err != nil {
		ctx.Status(http.StatusUnauthorized)
		return
	}

	tIDString := ctx.Param("id")
	tID, err := strconv.Atoi(tIDString)
	if err != nil || tID <= 0 {
		ctx.Status(http.StatusBadRequest)
		return
	}

	tModel, err := h.repo.TransactionGet(uint(tID))
	if err != nil {
		if err == repository.ErrorRecordNotFound {
			ctx.Status(http.StatusNotFound)
			return
		}
		ctx.Status(http.StatusInternalServerError)
		return
	}

	if tModel.UserID != userID {
		ctx.Status(http.StatusUnauthorized)
		return
	}

	tResponse := TransactionModelToResponse(tModel)

	ctx.JSON(http.StatusOK, tResponse)
}

func (h *handler) ListTransactions(ctx *gin.Context) {
	userID, err := auth_middleware.GetUserIDFromContext(ctx)
	if err != nil {
		ctx.Status(http.StatusUnauthorized)
		return
	}

	tModels, err := h.repo.TransactionList(userID)
	if err != nil {
		if err == repository.ErrorRecordNotFound {
			ctx.Status(http.StatusNotFound)
			return
		}
		ctx.Status(http.StatusInternalServerError)
		return
	}

	tResponses := make([]*Transaction, 0, len(tModels))

	for _, t := range tModels {
		tResponses = append(tResponses, TransactionModelToResponse(t))
	}

	res := NewListResponse(tResponses)
	ctx.JSON(http.StatusOK, res)
}
