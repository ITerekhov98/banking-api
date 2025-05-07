package handler

import (
	"encoding/json"
	"net/http"

	"banking-api/internal/cbr"
	"banking-api/internal/handler/dto"
)

type KeyRateHandler struct{}

func NewKeyRateHandler() *KeyRateHandler {
	return &KeyRateHandler{}
}

// GetKeyRate godoc
// @Summary     Get current key rate
// @Description Get current key rate from Central Bank
// @Tags        CBR
// @Produce     json
// @Success     200 {object} KeyRateResponse
// @Failure     500 {string} string
// @Router      /api/keyrate [get]
// @Security    BearerAuth
func (h *KeyRateHandler) GetKeyRate(w http.ResponseWriter, r *http.Request) {
	rate, err := cbr.GetKeyRate()
	if err != nil {
		http.Error(w, "failed to fetch key rate", http.StatusInternalServerError)
		return
	}

	resp := dto.KeyRateResponse{KeyRate: rate}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
