package handlers

import (
	"expense-api/model"
	"expense-api/repository"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TransactionsHandler interface {
	ListTransactions(ctx *gin.Context)
	CreateTransaction(ctx *gin.Context)
	GetTransaction(ctx *gin.Context)
	UpdateTransaction(ctx *gin.Context)
	DeleteTransaction(ctx *gin.Context)
}

func (h *handler) CreateTransaction(ctx *gin.Context) {
	var tRequest model.Transaction
	if err := ctx.Bind(&tRequest); err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

	err := model.TransactionValidateCreateBody(tRequest.Timestamp, tRequest.Amount, tRequest.Type)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	tResponse, err := h.repo.TransactionCreate(tRequest.Timestamp, tRequest.Amount, tRequest.Type)
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusCreated, tResponse)
}

func (h *handler) UpdateTransaction(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		ctx.Status(http.StatusBadRequest)
		return
	}

	var tRequest model.Transaction
	if err := ctx.Bind(&tRequest); err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

	if err := model.TransactionValidateUpdateBody(&tRequest); err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

	tResponse, err := h.repo.TransactionUpdate(uint(id), tRequest.Timestamp, tRequest.Amount, tRequest.Type)
	if err != nil {
		if err == repository.ErrorRecordNotFound {
			ctx.Status(http.StatusNotFound)
			return
		}

		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, tResponse)
}

func (h *handler) DeleteTransaction(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		ctx.Status(http.StatusBadRequest)
		return
	}

	if err := h.repo.TransactionDelete(uint(id)); err != nil {
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
	idStr := ctx.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		ctx.Status(http.StatusBadRequest)
		return
	}

	transaction, err := h.repo.TransactionGet(uint(id))
	if err != nil {
		if err == repository.ErrorRecordNotFound {
			ctx.Status(http.StatusNotFound)
			return
		}
		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, transaction)
}

func (h *handler) ListTransactions(ctx *gin.Context) {
	transactions, err := h.repo.TransactionList()
	if err != nil {
		if err == repository.ErrorRecordNotFound {
			ctx.Status(http.StatusNotFound)
			return
		}
		ctx.Status(http.StatusInternalServerError)
		return
	}

	res := NewListResponse(transactions)
	ctx.JSON(http.StatusOK, res)
}
