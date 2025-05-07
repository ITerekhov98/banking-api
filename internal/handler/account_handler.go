package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"banking-api/internal/middleware"
	"banking-api/internal/service"

	"github.com/gorilla/mux"
)

type AccountHandler struct {
	accountService *service.AccountService
}

func NewAccountHandler(accountService *service.AccountService) *AccountHandler {
	return &AccountHandler{accountService: accountService}
}

// CreateAccount godoc
// @Summary     Create a new bank account
// @Description Creates a new bank account for the authenticated user
// @Tags        Accounts
// @Accept      json
// @Produce     json
// @Success     201 {object} model.Account
// @Failure     400 {string} string
// @Failure     401 {string} string
// @Router      /api/accounts [post]
// @Security    BearerAuth
func (h *AccountHandler) CreateAccount(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(int64)

	account, err := h.accountService.CreateAccount(r.Context(), userID)
	if err != nil {
		http.Error(w, "failed to create account", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(account)
}

// Payload to make funds transfer
// swagger:model depositRequest
type transferRequest struct {
	// ID of sender account
	// example: 42
	FromAccountID int64 `json:"from_account_id"`
	// ID of the recipient account
	// example: 24
	ToAccountID int64 `json:"to_account_id"`
	// Amount of funds
	// example: 1000.1
	Amount float64 `json:"amount"`
}

// Transfer godoc
// @Summary     Transfer funds between accounts
// @Description Transfers money from one user-owned account to another account
// @Tags        Transactions
// @Accept      json
// @Produce     json
// @Param       transferRequest body transferRequest true "Transfer details"
// @Success     200 {string} string "ok"
// @Failure     400 {string} string
// @Failure     401 {string} string
// @Failure     403 {string} string
// @Router      /api/transfer [post]
// @Security    BearerAuth
func (h *AccountHandler) Transfer(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(int64)

	var req transferRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	err := h.accountService.Transfer(r.Context(), userID, req.FromAccountID, req.ToAccountID, req.Amount)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}

// Payload to make transactions with account - deposit or withdrawal
// swagger:model depositRequest
type depositRequest struct {
	// ID of the target account
	// example: 42
	AccountID int64 `json:"account_id"`
	// Amount of funds
	// example: 1000.1
	Amount float64 `json:"amount"`
}

// Deposit godoc
// @Summary     Deposit money to account
// @Description Adds funds to the specified user-owned account
// @Tags        Transactions
// @Accept      json
// @Produce     json
// @Param       depositRequest body depositRequest true "Deposit details"
// @Success     200 {object} model.Account
// @Failure     400 {string} string
// @Failure     401 {string} string
// @Router      /api/deposit [post]
// @Security    BearerAuth
func (h *AccountHandler) Deposit(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(int64)

	var req depositRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	account, err := h.accountService.Deposit(r.Context(), userID, req.AccountID, req.Amount)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(account)
}

// GetAccount godoc
// @Summary     Get account details
// @Description Returns information about a bank account if it belongs to the user
// @Tags        Accounts
// @Produce     json
// @Param       id path int true "Account ID"
// @Success     200 {object} model.Account
// @Failure     400 {string} string
// @Failure     403 {string} string
// @Failure     404 {string} string
// @Router      /api/accounts/{id} [get]
// @Security    BearerAuth
func (h *AccountHandler) Get(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(int64)

	vars := mux.Vars(r)
	accountID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "invalid credit id", http.StatusBadRequest)
		return
	}

	account, err := h.accountService.Get(r.Context(), userID, accountID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(account)
}

// Withdraw godoc
// @Summary     Withdraw money from account
// @Description Deducts funds from the specified user-owned account
// @Tags        Transactions
// @Accept      json
// @Produce     json
// @Param       withdrawRequest body depositRequest true "Withdraw details"
// @Success     200 {object} model.Account
// @Failure     400 {string} string
// @Failure     401 {string} string
// @Router      /api/withdraw [post]
// @Security    BearerAuth
func (h *AccountHandler) Withdraw(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(int64)

	var req depositRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	account, err := h.accountService.Withdraw(r.Context(), userID, req.AccountID, req.Amount)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(account)
}
