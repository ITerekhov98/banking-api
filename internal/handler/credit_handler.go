package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"banking-api/internal/middleware"
	"banking-api/internal/service"
)

type CreditHandler struct {
	creditService *service.CreditService
}

func NewCreditHandler(creditService *service.CreditService) *CreditHandler {
	return &CreditHandler{creditService: creditService}
}

type createCreditRequest struct {
	AccountID    int64   `json:"account_id"`
	Principal    float64 `json:"principal"`
	InterestRate float64 `json:"interest_rate"`
	TermMonths   int     `json:"term_months"`
}

func (h *CreditHandler) CreateCredit(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(int64)

	var req createCreditRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	credit, err := h.creditService.CreateCredit(r.Context(), userID, req.AccountID, req.Principal, req.InterestRate, req.TermMonths)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(credit)
}

func (h *CreditHandler) GetSchedule(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(int64)

	vars := mux.Vars(r)
	creditID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "invalid credit id", http.StatusBadRequest)
		return
	}

	schedule, err := h.creditService.GetSchedule(r.Context(), userID, creditID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(schedule)
}
