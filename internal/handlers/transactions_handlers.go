package handlers

import (
	"expense-api/internal/middleware"
	auth_middleware "expense-api/internal/middleware/auth"
	"expense-api/internal/repository"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)

type TransactionsHandler interface {
	ListTransactions(ctx *gin.Context)
	CreateTransaction(ctx *gin.Context)
	GetTransaction(ctx *gin.Context)
	UpdateTransaction(ctx *gin.Context)
	DeleteTransaction(ctx *gin.Context)
}

func (h *handler) CreateTransaction(ctx *gin.Context) {
	userID, err := auth_middleware.GetUserIDFromContext(ctx)
	if err != nil {
		ctx.Status(http.StatusForbidden)
		return
	}

	var tRequest Transaction
	if err := ctx.Bind(&tRequest); err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

	if tRequest.Amount.Cmp(decimal.Zero) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": ErrorRequiredAmount.Error(),
		})
		return
	}

	tModel := TransactionRequestToModel(&tRequest, userID)

	{ // Validate wallet ownership
		if tModel.WalletID == 0 {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": ErrorRequiredWalletID.Error(),
			})
			return
		}

		wallet, err := h.repo.WalletGet(tModel.WalletID)
		if err != nil {
			if err == repository.ErrorRecordNotFound {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"message": ErrorWalletNotFound.Error(),
				})
				return
			}
			ctx.Status(http.StatusInternalServerError)
			return
		}

		if wallet.UserID != userID {
			ctx.JSON(http.StatusForbidden, gin.H{
				"message": ErrorBadWalletID.Error(),
			})
			return
		}
	}

	{ // Validate party ownership
		if tModel.PartyID == 0 {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": ErrorRequiredPartyID.Error(),
			})
			return
		}

		party, err := h.repo.PartyGet(tModel.PartyID)
		if err != nil {
			if err == repository.ErrorRecordNotFound {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"message": ErrorPartyNotFound.Error(),
				})
				return
			}
			ctx.Status(http.StatusInternalServerError)
			return
		}

		if party.UserID != userID {
			ctx.JSON(http.StatusForbidden, gin.H{
				"message": ErrorBadPartyID.Error(),
			})
			return
		}
	}

	if err := h.repo.TransactionCreate(tModel); err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	tResponse := TransactionModelToResponse(tModel)
	ctx.JSON(http.StatusCreated, tResponse)
}

func (h *handler) GetTransaction(ctx *gin.Context) {
	id := middleware.GetIDParamFromContext(ctx)

	tModel, err := h.repo.TransactionGet(id)
	if err != nil {
		if err == repository.ErrorRecordNotFound {
			ctx.Status(http.StatusNotFound)
			return
		}
		ctx.Status(http.StatusInternalServerError)
		return
	}

	tResponse := TransactionModelToResponse(tModel)
	ctx.JSON(http.StatusOK, tResponse)
}

func (h *handler) UpdateTransaction(ctx *gin.Context) {
	userID, err := auth_middleware.GetUserIDFromContext(ctx)
	if err != nil {
		ctx.Status(http.StatusForbidden)
		return
	}

	id := middleware.GetIDParamFromContext(ctx)

	var tRequest Transaction
	if err := ctx.Bind(&tRequest); err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

	tModel := TransactionRequestToModel(&tRequest, userID)

	// Validate wallet ownership
	if tModel.WalletID != 0 {
		wallet, err := h.repo.WalletGet(tModel.WalletID)
		if err != nil {
			if err == repository.ErrorRecordNotFound {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"message": ErrorWalletNotFound.Error(),
				})
				return
			}
			ctx.Status(http.StatusInternalServerError)
			return
		}

		if wallet.UserID != userID {
			ctx.JSON(http.StatusForbidden, gin.H{
				"message": ErrorBadWalletID.Error(),
			})
			return
		}
	}

	// Validate party ownership
	if tModel.PartyID != 0 {
		party, err := h.repo.PartyGet(tModel.PartyID)
		if err != nil {
			if err == repository.ErrorRecordNotFound {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"message": ErrorPartyNotFound.Error(),
				})
				return
			}
			ctx.Status(http.StatusInternalServerError)
			return
		}

		if party.UserID != userID {
			ctx.JSON(http.StatusForbidden, gin.H{
				"message": ErrorBadPartyID.Error(),
			})
			return
		}
	}

	updatedTModel, err := h.repo.TransactionUpdate(id, tModel)
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	tResponse := TransactionModelToResponse(updatedTModel)
	ctx.JSON(http.StatusOK, tResponse)
}

func (h *handler) DeleteTransaction(ctx *gin.Context) {
	id := middleware.GetIDParamFromContext(ctx)

	if err := h.repo.TransactionDelete(id); err != nil {
		if err == repository.ErrorRecordNotFound {
			ctx.Status(http.StatusNotFound)
			return
		}
		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (h *handler) ListTransactions(ctx *gin.Context) {
	userID, err := auth_middleware.GetUserIDFromContext(ctx)
	if err != nil {
		ctx.Status(http.StatusForbidden)
		return
	}

	tModels, err := h.repo.TransactionList(userID)
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	tResponse := make([]*Transaction, 0, len(tModels))
	for _, t := range tModels {
		tResponse = append(tResponse, TransactionModelToResponse(t))
	}

	res := NewListResponse(tResponse)
	ctx.JSON(http.StatusOK, res)
}
