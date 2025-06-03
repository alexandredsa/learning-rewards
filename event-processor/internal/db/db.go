package db

import (
	"event-processor/internal/models"
	"log"
	"os"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(dsn string) (*gorm.DB, error) {
	if dsn == "" {
		dsn = "postgres://user:pass@localhost:5433/event_processor?sslmode=disable"
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
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
