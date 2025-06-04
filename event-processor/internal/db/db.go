package db

import (
	"event-processor/internal/models"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	maxRetries = 5
	baseDelay  = time.Second
)

func connectWithRetry(dsn string) (*gorm.DB, error) {
	var db *gorm.DB
	var err error
	delay := baseDelay

	for attempt := 1; attempt <= maxRetries; attempt++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			return db, nil
		}

		if attempt == maxRetries {
			return nil, fmt.Errorf("failed to connect after %d attempts: %v", maxRetries, err)
		}

		log.Printf("Database connection attempt %d failed: %v. Retrying in %v...", attempt, err, delay)
		time.Sleep(delay)
		delay *= 2 // exponential backoff
	}

	return nil, fmt.Errorf("unexpected error in retry logic")
}

func Connect(dsn string) (*gorm.DB, error) {
	if dsn == "" {
		dsn = "postgres://user:pass@localhost:5433/event_processor?sslmode=disable"
	}

	db, err := connectWithRetry(dsn)
	if err != nil {
		return nil, err
	}

	log.Println("Connected to DB successfully")

	if os.Getenv("ENV") == "dev" {
		if err := Migrate(db); err != nil {
			return nil, err
		}
	}

	return db, nil
}

func Migrate(db *gorm.DB) error {
	// Enable uuid extension (once)
	db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`)

	// Run auto migration
	if err := db.AutoMigrate(&models.Event{}); err != nil {
		log.Printf("failed to auto-migrate: %v", err)
		return err
	}

	return nil
}
