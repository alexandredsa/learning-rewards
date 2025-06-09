package database

import (
	"fmt"
	"log"
	"time"

	"github.com/alexandredsa/learning-rewards/reward-processor/pkg/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
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
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
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

// Connect creates a new database connection and runs auto-migrations
func Connect(dsn string) (*gorm.DB, error) {
	if dsn == "" {
		dsn = "postgres://user:pass@localhost:5433/reward_processor?sslmode=disable"
	}

	db, err := connectWithRetry(dsn)
	if err != nil {
		return nil, err
	}

	log.Println("Connected to DB successfully")

	// Auto-migrate the schema
	if err := db.AutoMigrate(&models.UserEventCount{}); err != nil {
		return nil, fmt.Errorf("failed to auto-migrate database: %w", err)
	}

	return db, nil
}
