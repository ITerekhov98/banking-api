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

// GetMonthlyStats godoc
// @Summary     Get monthly income and expenses
// @Description Returns total income and expenses for the current month
// @Tags        Analytics
// @Produce     json
// @Success     200 {object} model.MonthlyStats
// @Failure     401 {string} string
// @Failure     500 {string} string
// @Router      /api/analytics [get]
// @Security    BearerAuth
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

// swagger:model predictedBalanceResponse
type predictedBalanceResponse struct {
	// Predicted account balance
	// example: 123456.78
	PredictedBalance float64 `json:"predicted_balance"`
}

// GetPredictedBalance godoc
// @Summary     Predict account balance
// @Description Returns projected account balance based on upcoming credit payments
// @Tags        Analytics
// @Produce     json
// @Param       days query int true "Prediction period in days (max 365)"
// @Success 200 {object} predictedBalanceResponse
// @Failure     400 {string} string
// @Failure     401 {string} string
// @Router      /api/analytics/predict [get]
// @Security    BearerAuth
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
	response := predictedBalanceResponse{PredictedBalance: pred}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
