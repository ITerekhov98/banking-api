package handler

import (
	"banking-api/internal/handler/dto"
	"banking-api/internal/service"
	"encoding/json"
	"net/http"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Register godoc
// @Summary     Register a new user
// @Description Registers a user with unique email and username, returns user info
// @Tags        Auth
// @Accept      json
// @Produce     json
// @Param       registerRequest body registerRequest true "Registration data"
// @Success     201 {object} RegisterResponse
// @Failure     400 {string} string
// @Router      /register [post]
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	user, err := h.authService.Register(r.Context(), req.Email, req.Username, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	resp := dto.RegisterResponse{
		Id:       user.ID,
		Email:    user.Email,
		Username: user.Username,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)

}

// Login godoc
// @Summary     Authenticate user and return JWT token
// @Description Validates credentials and returns access token if successful
// @Tags        Auth
// @Accept      json
// @Produce     json
// @Param       loginRequest body LoginRequest true "Login data"
// @Success     200 {object} LoginResponse
// @Failure     400 {string} string
// @Router      /login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	token, err := h.authService.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	resp := dto.LoginResponse{Token: token}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
