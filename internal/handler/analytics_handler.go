package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"banking-api/internal/middleware"
	"banking-api/internal/service"
)

type AnalyticsHandler struct {
	service *service.AnalyticsService
}

func NewAnalyticsHandler(service *service.AnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{service: service}
}

func (h *AnalyticsHandler) GetMonthlyStats(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(int64)

	stats, err := h.service.GetMonthlyStats(r.Context(), userID)
	if err != nil {
		http.Error(w, "failed to get stats", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

func (h *AnalyticsHandler) GetPredictedBalance(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(int64)
	daysStr := r.URL.Query().Get("days")
	if daysStr == "" {
		http.Error(w, "missing days parameter", http.StatusBadRequest)
		return
	}
	days, err := strconv.Atoi(daysStr)
	if err != nil {
		http.Error(w, "invalid days value", http.StatusBadRequest)
		return
	}

	pred, err := h.service.GetPredictedBalance(r.Context(), userID, days)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(map[string]float64{
		"predicted_balance": pred,
	})
}
