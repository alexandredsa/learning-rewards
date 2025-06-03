package models

import (
	"github.com/google/uuid"
)

type Course struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey"`
	Title      string    `gorm:"not null"`
	CategoryID uuid.UUID
	Category   Category `gorm:"foreignKey:CategoryID"`
}
