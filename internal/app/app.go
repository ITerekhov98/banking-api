package app

import (
	"banking-api/internal/handler"
	"banking-api/internal/middleware"
	"banking-api/internal/repository"
	"banking-api/internal/service"
	"banking-api/pkg/logger"

	"github.com/gorilla/mux"
)

func SetupRouter() *mux.Router {
	router := mux.NewRouter()

	// for swagger doc
	router.Use(middleware.CORS("http://localhost:8080", "http://127.0.0.1:8080"))

	userRepo := &repository.UserRepository{}
	authService := service.NewAuthService(userRepo)
	authHandler := handler.NewAuthHandler(authService)

	router.HandleFunc("/register", authHandler.Register).Methods("POST", "OPTIONS")
	router.HandleFunc("/login", authHandler.Login).Methods("POST", "OPTIONS")

	api := router.PathPrefix("/api").Subrouter()
	api.Use(middleware.AuthMiddleware)

	accountRepo := &repository.AccountRepository{}
	accountService := service.NewAccountService(accountRepo)
	accountHandler := handler.NewAccountHandler(accountService)

	api.HandleFunc("/accounts", accountHandler.CreateAccount).Methods("POST")
	api.HandleFunc("/accounts/{id:[0-9]+}", accountHandler.Get).Methods("GET")
	api.HandleFunc("/transfer", accountHandler.Transfer).Methods("POST")
	api.HandleFunc("/deposit", accountHandler.Deposit).Methods("POST")
	api.HandleFunc("/withdraw", accountHandler.Withdraw).Methods("POST")

	cardRepo := &repository.CardRepository{}
	cardService := service.NewCardervice(cardRepo, accountRepo)
	cardHandler := handler.NewCardHandler(cardService)

	api.HandleFunc("/cards", cardHandler.CreateCard).Methods("POST")
	api.HandleFunc("/cards", cardHandler.GetUserCards).Methods("GET")

	creditRepo := &repository.CreditRepository{}
	creditService := service.NewCreditService(accountRepo, creditRepo)
	creditHandler := handler.NewCreditHandler(creditService)

	api.HandleFunc("/credits", creditHandler.CreateCredit).Methods("POST")
	api.HandleFunc("/credits/{id:[0-9]+}/schedule", creditHandler.GetSchedule).Methods("GET")

	analyticsRepo := &repository.AnalyticsRepository{}
	analyticsService := service.NewAnalyticsService(analyticsRepo)
	analyticsHandler := handler.NewAnalyticsHandler(analyticsService)

	api.HandleFunc("/analytics", analyticsHandler.GetMonthlyStats).Methods("GET")
	api.HandleFunc("/analytics/predict", analyticsHandler.GetPredictedBalance).Methods("GET")

	keyRateHandler := handler.NewKeyRateHandler()
	api.HandleFunc("/keyrate", keyRateHandler.GetKeyRate).Methods("GET")

	logger.Info("Router initialized")

	return router
}
