package handlers

import (
	"expense-api/model"
	"expense-api/repository"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreateTransaction(r repository.Repository) func(*gin.Context) {
	return func(ctx *gin.Context) {
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

		tResponse, err := r.TransactionCreate(tRequest.Timestamp, tRequest.Amount, tRequest.Type)
		if err != nil {
			ctx.Status(http.StatusInternalServerError)
			return
		}

		ctx.JSON(http.StatusCreated, tResponse)
	}
}

func UpdateTransaction(r repository.Repository) func(*gin.Context) {
	return func(ctx *gin.Context) {
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

		tResponse, err := r.TransactionUpdate(uint(id), tRequest.Timestamp, tRequest.Amount, tRequest.Type)
		if err != nil {
			ctx.Status(http.StatusInternalServerError)
			return
		}

		ctx.JSON(http.StatusOK, tResponse)
	}
}

func DeleteTransaction(r repository.Repository) func(*gin.Context) {
	return func(ctx *gin.Context) {
		idStr := ctx.Param("id")

		id, err := strconv.Atoi(idStr)
		if err != nil || id <= 0 {
			ctx.Status(http.StatusBadRequest)
			return
		}

		if err := r.TransactionDelete(uint(id)); err != nil {
			if err == repository.ErrorRecordNotFound {
				ctx.Status(http.StatusNotFound)
				return
			}
			ctx.Status(http.StatusInternalServerError)
			return
		}

		ctx.Status(http.StatusNoContent)
	}
}

func GetTransaction(r repository.Repository) func(*gin.Context) {
	return func(ctx *gin.Context) {
		idStr := ctx.Param("id")

		id, err := strconv.Atoi(idStr)
		if err != nil || id <= 0 {
			ctx.Status(http.StatusBadRequest)
			return
		}

		transaction, err := r.TransactionGet(uint(id))
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
}

func ListTransactions(r repository.Repository) func(*gin.Context) {
	return func(ctx *gin.Context) {
		transactions, err := r.TransactionList()
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
}
