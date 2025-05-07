package handler

import (
	"banking-api/internal/handler/dto"
	"banking-api/internal/middleware"
	"banking-api/internal/service"
	"encoding/json"
	"net/http"
)

type CardHandler struct {
	cardService *service.CardService
}

func NewCardHandler(cardService *service.CardService) *CardHandler {
	return &CardHandler{cardService: cardService}
}

// CreateCard godoc
// @Summary     Issue a new virtual card
// @Description Generates a new encrypted virtual card for the specified account
// @Tags        Cards
// @Accept      json
// @Produce     json
// @Param       createCardRequest body createCardRequest true "Target account ID"
// @Success     201 {object} model.CardPlain
// @Failure     400 {string} string
// @Failure     401 {string} string
// @Router      /api/cards [post]
// @Security    BearerAuth
func (h *CardHandler) CreateCard(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(int64)

	var req dto.CreateCardRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	card, err := h.cardService.CreateCard(r.Context(), userID, req.AccountID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(card)
}

// GetUserCards godoc
// @Summary     Get user's virtual cards
// @Description Returns decrypted virtual card data (number, expiry) for the authenticated user
// @Tags        Cards
// @Produce     json
// @Success     200 {array} model.CardRaw
// @Failure     401 {string} string
// @Failure     500 {string} string
// @Router      /api/cards [get]
// @Security    BearerAuth
func (h *CardHandler) GetUserCards(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(int64)

	cards, err := h.cardService.GetUserCards(r.Context(), userID)
	if err != nil {
		http.Error(w, "failed to fetch cards", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cards)
}
