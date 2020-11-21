package handlers

import (
	"expense-api/middleware"
	auth_middleware "expense-api/middleware/auth"
	"expense-api/model"
	"expense-api/repository"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)

type TransactionHandler interface {
	ListTransactions(ctx *gin.Context)
	CreateTransaction(ctx *gin.Context)
	GetTransaction(ctx *gin.Context)
	UpdateTransaction(ctx *gin.Context)
	DeleteTransaction(ctx *gin.Context)
}

// Transaction is a transaction with an omitted user
type Transaction struct {
	ID          uint            `json:"id"`
	WalletID    uint            `json:"wallet_id"`
	PartyID     uint            `json:"party_id"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	Timestamp   time.Time       `json:"timestamp"`
	Amount      decimal.Decimal `json:"amount"`
	Description string          `json:"description"`
}

func TransactionModelToResponse(t *model.Transaction) *Transaction {
	return &Transaction{
		ID:          t.ID,
		WalletID:    t.WalletID,
		PartyID:     t.PartyID,
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
		Timestamp:   t.Timestamp,
		Amount:      t.Amount,
		Description: t.Description,
	}
}

func TransactionRequestToModel(t *Transaction, userID uint) *model.Transaction {
	return &model.Transaction{
		Amount:      t.Amount,
		Timestamp:   t.Timestamp,
		Description: t.Description,
		WalletID:    t.WalletID,
		PartyID:     t.PartyID,
		UserID:      userID,
	}
}

const (
	ErrMsgRequiredAmount   = "cannot create new transaction with an amount of 0"
	ErrMsgRequiredWalletID = "a valid wallet id must be specified to register a new transaction"
	ErrMsgRequiredPartyID  = "a valid wallet id must be specified to register a new transaction"
	ErrMsgWalletNotFound   = "wallet with specified id not found"
	ErrMsgBadWalletID      = "wallet with specified id belongs to another user"
	ErrMsgPartyNotFound    = "party with specified id not found"
	ErrMsgBadPartyID       = "party with specified id belongs to another user"
)

func (h *handler) CreateTransaction(ctx *gin.Context) {
	userID, err := auth_middleware.GetUserIDFromContext(ctx)
	if err != nil {
		ctx.Status(http.StatusUnauthorized)
		return
	}

	var tRequest Transaction
	if err := ctx.Bind(&tRequest); err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

	if tRequest.Amount.Cmp(decimal.Zero) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": ErrMsgRequiredAmount,
		})
		return
	}

	tModel := TransactionRequestToModel(&tRequest, userID)

	{ // Validate wallet ownership
		if tModel.WalletID == 0 {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": ErrMsgRequiredWalletID,
			})
			return
		}

		wallet, err := h.repo.WalletGet(tModel.WalletID)
		if err != nil {
			if err == repository.ErrorRecordNotFound {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"message": ErrMsgWalletNotFound,
				})
				return
			}
			ctx.Status(http.StatusInternalServerError)
			return
		}

		if wallet.UserID != userID {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"message": ErrMsgBadWalletID,
			})
			return
		}
	}

	{ // Validate party ownership
		if tModel.PartyID == 0 {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": ErrMsgRequiredPartyID,
			})
			return
		}

		party, err := h.repo.PartyGet(tModel.PartyID)
		if err != nil {
			if err == repository.ErrorRecordNotFound {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"message": ErrMsgPartyNotFound,
				})
				return
			}
			ctx.Status(http.StatusInternalServerError)
			return
		}

		if party.UserID != userID {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"message": ErrMsgBadPartyID,
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

func (h *handler) UpdateTransaction(ctx *gin.Context) {
	userID, err := auth_middleware.GetUserIDFromContext(ctx)
	if err != nil {
		ctx.Status(http.StatusUnauthorized)
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
					"message": ErrMsgWalletNotFound,
				})
				return
			}
			ctx.Status(http.StatusInternalServerError)
			return
		}

		if wallet.UserID != userID {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"message": ErrMsgBadWalletID,
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
					"message": ErrMsgPartyNotFound,
				})
				return
			}
			ctx.Status(http.StatusInternalServerError)
			return
		}

		if party.UserID != userID {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"message": ErrMsgBadPartyID,
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

func (h *handler) ListTransactions(ctx *gin.Context) {
	userID, err := auth_middleware.GetUserIDFromContext(ctx)
	if err != nil {
		ctx.Status(http.StatusUnauthorized)
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
