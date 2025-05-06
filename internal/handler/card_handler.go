package handler

import (
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

type createCardRequest struct {
	AccountID int64 `json:"account_id"`
}

func (h *CardHandler) CreateCard(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(int64)

	var req createCardRequest
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
