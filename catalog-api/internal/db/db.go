package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Init() {
	var err error

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s sslmode=disable password=%s",
		getEnv("DB_HOST", "localhost"),
		getEnv("DB_PORT", "5432"),
		getEnv("DB_USER", "user"),
		getEnv("DB_NAME", "catalog"),
		getEnv("DB_PASSWORD", "pass"),
	)

	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("failed to connect to DB: %v", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatalf("cannot reach DB: %v", err)
	}

	log.Println("Connected to DB")
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
