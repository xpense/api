package handlers

import (
	"expense-api/middleware"
	auth_middleware "expense-api/middleware/auth"
	"expense-api/model"
	"expense-api/repository"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type WalletsHandler interface {
	ListWallets(ctx *gin.Context)
	CreateWallet(ctx *gin.Context)
	GetWallet(ctx *gin.Context)
	UpdateWallet(ctx *gin.Context)
	DeleteWallet(ctx *gin.Context)
	ListTransactionsByWallet(ctx *gin.Context)
}

// Wallet is a list of transactions belonging to an account
type Wallet struct {
	ID          uint      `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
}

func WalletModelToResponse(w *model.Wallet) *Wallet {
	return &Wallet{
		ID:          w.ID,
		CreatedAt:   w.CreatedAt,
		UpdatedAt:   w.UpdatedAt,
		Name:        w.Name,
		Description: w.Description,
	}
}

func WalletRequestToModel(w *Wallet, userID uint) *model.Wallet {
	return &model.Wallet{
		Name:        w.Name,
		Description: w.Description,
		UserID:      userID,
	}
}

const ErrMsgWalletNameTaken = "wallet with the same name, belonging to the same user already exists"

func (h *handler) CreateWallet(ctx *gin.Context) {
	userID, err := auth_middleware.GetUserIDFromContext(ctx)
	if err != nil {
		ctx.Status(http.StatusForbidden)
		return
	}

	var wRequest Wallet
	if err := ctx.Bind(&wRequest); err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

	wModel := WalletRequestToModel(&wRequest, userID)

	if err := h.repo.WalletCreate(wModel); err != nil {
		if err == repository.ErrorUniqueConstaintViolation {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": ErrMsgWalletNameTaken,
			})
			return
		}
		ctx.Status(http.StatusInternalServerError)
		return
	}

	wResponse := WalletModelToResponse(wModel)
	ctx.JSON(http.StatusCreated, wResponse)
}

func (h *handler) UpdateWallet(ctx *gin.Context) {
	userID, err := auth_middleware.GetUserIDFromContext(ctx)
	if err != nil {
		ctx.Status(http.StatusForbidden)
		return
	}

	id := middleware.GetIDParamFromContext(ctx)

	var wRequest Wallet
	if err := ctx.Bind(&wRequest); err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

	wModel := WalletRequestToModel(&wRequest, userID)
	updatedWModel, err := h.repo.WalletUpdate(id, wModel)
	if err != nil {
		if err == repository.ErrorRecordNotFound {
			ctx.Status(http.StatusNotFound)
			return
		}
		if err == repository.ErrorUniqueConstaintViolation {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": ErrMsgWalletNameTaken,
			})
			return
		}
		ctx.Status(http.StatusInternalServerError)
		return
	}

	wResponse := WalletModelToResponse(updatedWModel)
	ctx.JSON(http.StatusOK, wResponse)
}

func (h *handler) DeleteWallet(ctx *gin.Context) {
	id := middleware.GetIDParamFromContext(ctx)

	if err := h.repo.WalletDelete(id); err != nil {
		if err == repository.ErrorRecordNotFound {
			ctx.Status(http.StatusNotFound)
			return
		}
		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (h *handler) GetWallet(ctx *gin.Context) {
	id := middleware.GetIDParamFromContext(ctx)

	wModel, err := h.repo.WalletGet(id)
	if err != nil {
		if err == repository.ErrorRecordNotFound {
			ctx.Status(http.StatusNotFound)
			return
		}
		ctx.Status(http.StatusInternalServerError)
		return
	}

	wResponse := WalletModelToResponse(wModel)
	ctx.JSON(http.StatusOK, wResponse)
}

func (h *handler) ListWallets(ctx *gin.Context) {
	userID, err := auth_middleware.GetUserIDFromContext(ctx)
	if err != nil {
		ctx.Status(http.StatusForbidden)
		return
	}

	wModels, err := h.repo.WalletList(userID)
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	wResponse := make([]*Wallet, 0, len(wModels))

	for _, w := range wModels {
		wResponse = append(wResponse, WalletModelToResponse(w))
	}

	res := NewListResponse(wResponse)
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) ListTransactionsByWallet(ctx *gin.Context) {
	userID, err := auth_middleware.GetUserIDFromContext(ctx)
	if err != nil {
		ctx.Status(http.StatusForbidden)
		return
	}

	id := middleware.GetIDParamFromContext(ctx)

	tModels, err := h.repo.TransactionListByWallet(userID, id)
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
