package main

import (
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"

	"banking-api/internal/app"
	"banking-api/internal/db"
	"banking-api/internal/repository"
	scheduler "banking-api/internal/schedluer"
	"banking-api/pkg/logger"
)

func main() {
	// load variables from .env
	logger.Init()

	err := godotenv.Load()
	if err != nil {
		logger.Fatal("No .env file found")
	}

	// init DB connection
	if err := db.Init(); err != nil {
		logger.Fatal("DB init error: ", err)
	}

	// run payment scheduler
	repo := &repository.CreditRepository{}
	scheduler.StartCreditPaymentScheduler(repo, 2*time.Hour)

	// run application
	router := app.SetupRouter()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	logger.Info("Server starting on port ", port)
	err = http.ListenAndServe(":"+port, router)
	if err != nil {
		logger.Fatal("Server failed: ", err)
	}
}
