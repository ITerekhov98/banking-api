package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"banking-api/internal/handler/dto"
	"banking-api/internal/middleware"
	"banking-api/internal/service"
)

type CreditHandler struct {
	creditService *service.CreditService
}

func NewCreditHandler(creditService *service.CreditService) *CreditHandler {
	return &CreditHandler{creditService: creditService}
}

// CreateCredit godoc
// @Summary     Create a new credit
// @Description Issues a new credit for the user with annuity payment schedule
// @Tags        Credits
// @Accept      json
// @Produce     json
// @Param       createCreditRequest body createCreditRequest true "Credit details"
// @Success     201 {object} model.Credit
// @Failure     400 {string} string
// @Failure     401 {string} string
// @Router      /api/credits [post]
// @Security    BearerAuth
func (h *CreditHandler) CreateCredit(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(int64)

	var req dto.CreateCreditRequest
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

// GetSchedule godoc
// @Summary     Get credit payment schedule
// @Description Returns a full monthly payment schedule for a specific credit
// @Tags        Credits
// @Produce     json
// @Param       id path int true "Credit ID"
// @Success     200 {array} model.PaymentSchedule
// @Failure     400 {string} string
// @Failure     401 {string} string
// @Router      /api/credits/{id}/schedule [get]
// @Security    BearerAuth
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
