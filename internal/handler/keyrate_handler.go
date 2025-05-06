package handler

import (
	"encoding/json"
	"net/http"

	"banking-api/internal/cbr"
)

type KeyRateHandler struct{}

func NewKeyRateHandler() *KeyRateHandler {
	return &KeyRateHandler{}
}

func (h *KeyRateHandler) GetKeyRate(w http.ResponseWriter, r *http.Request) {
	rate, err := cbr.GetKeyRate()
	if err != nil {
		http.Error(w, "failed to fetch key rate", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]float64{
		"key_rate": rate,
	})
}
