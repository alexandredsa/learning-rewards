package db

import (
	"catalog-api/internal/models"
	"log"
	"os"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(dsn string) (*gorm.DB, error) {
	if dsn == "" {
		dsn = "postgres://user:pass@localhost:5432/catalog?sslmode=disable"
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to Postgres: %v", err)
	}

	log.Println("Connected to DB successfully")

	if os.Getenv("ENV") == "dev" {
		db.AutoMigrate(&models.Category{}, &models.Course{})
		log.Println("Seeding database...")
		models.SeedDB(db)
	}

	return db, nil
}
