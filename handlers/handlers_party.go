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

type PartyHandler interface {
	ListParties(ctx *gin.Context)
	CreateParty(ctx *gin.Context)
	GetParty(ctx *gin.Context)
	UpdateParty(ctx *gin.Context)
	DeleteParty(ctx *gin.Context)
	ListTransactionsByParty(ctx *gin.Context)
}

// Party is a list of transactions belonging to an account
type Party struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
}

func PartyModelToResponse(p *model.Party) *Party {
	return &Party{
		ID:        p.ID,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
		Name:      p.Name,
	}
}

func PartyRequestToModel(p *Party, userID uint) *model.Party {
	return &model.Party{
		Name:   p.Name,
		UserID: userID,
	}
}

const ErrMsgPartyNameTaken = "party with the same name, belonging to the same user already exists"

func (h *handler) CreateParty(ctx *gin.Context) {
	userID, err := auth_middleware.GetUserIDFromContext(ctx)
	if err != nil {
		ctx.Status(http.StatusForbidden)
		return
	}

	var wRequest Party
	if err := ctx.Bind(&wRequest); err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

	wModel := PartyRequestToModel(&wRequest, userID)

	if err := h.repo.PartyCreate(wModel); err != nil {
		if err == repository.ErrorUniqueConstaintViolation {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": ErrMsgPartyNameTaken,
			})
			return
		}
		ctx.Status(http.StatusInternalServerError)
		return
	}

	wResponse := PartyModelToResponse(wModel)
	ctx.JSON(http.StatusCreated, wResponse)
}

func (h *handler) UpdateParty(ctx *gin.Context) {
	userID, err := auth_middleware.GetUserIDFromContext(ctx)
	if err != nil {
		ctx.Status(http.StatusForbidden)
		return
	}

	id := middleware.GetIDParamFromContext(ctx)

	var wRequest Party
	if err := ctx.Bind(&wRequest); err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

	wModel := PartyRequestToModel(&wRequest, userID)
	updatedWModel, err := h.repo.PartyUpdate(id, wModel)
	if err != nil {
		if err == repository.ErrorRecordNotFound {
			ctx.Status(http.StatusNotFound)
			return
		}
		if err == repository.ErrorUniqueConstaintViolation {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": ErrMsgPartyNameTaken,
			})
			return
		}
		ctx.Status(http.StatusInternalServerError)
		return
	}

	wResponse := PartyModelToResponse(updatedWModel)
	ctx.JSON(http.StatusOK, wResponse)
}

func (h *handler) DeleteParty(ctx *gin.Context) {
	id := middleware.GetIDParamFromContext(ctx)

	if err := h.repo.PartyDelete(id); err != nil {
		if err == repository.ErrorRecordNotFound {
			ctx.Status(http.StatusNotFound)
			return
		}
		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (h *handler) GetParty(ctx *gin.Context) {
	id := middleware.GetIDParamFromContext(ctx)

	wModel, err := h.repo.PartyGet(id)
	if err != nil {
		if err == repository.ErrorRecordNotFound {
			ctx.Status(http.StatusNotFound)
			return
		}
		ctx.Status(http.StatusInternalServerError)
		return
	}

	wResponse := PartyModelToResponse(wModel)
	ctx.JSON(http.StatusOK, wResponse)
}

func (h *handler) ListParties(ctx *gin.Context) {
	userID, err := auth_middleware.GetUserIDFromContext(ctx)
	if err != nil {
		ctx.Status(http.StatusForbidden)
		return
	}

	wModels, err := h.repo.PartyList(userID)
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	wResponse := make([]*Party, 0, len(wModels))

	for _, p := range wModels {
		wResponse = append(wResponse, PartyModelToResponse(p))
	}

	res := NewListResponse(wResponse)
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) ListTransactionsByParty(ctx *gin.Context) {
	userID, err := auth_middleware.GetUserIDFromContext(ctx)
	if err != nil {
		ctx.Status(http.StatusForbidden)
		return
	}

	id := middleware.GetIDParamFromContext(ctx)

	tModels, err := h.repo.TransactionListByParty(userID, id)
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
