package db

import (
	"database/sql"
	"fmt"
	"os"

	"banking-api/pkg/logger"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Init() error {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	sslmode := os.Getenv("DB_SSLMODE")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode)

	var err error
	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("could not open db: %w", err)
	}

	err = DB.Ping()
	if err != nil {
		return fmt.Errorf("could not ping db: %w", err)
	}

	logger.Info("Connected to PostgreSQL")
	return nil
}
