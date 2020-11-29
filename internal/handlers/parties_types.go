package handlers

import (
	"expense-api/internal/model"
	"time"
)

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
