package models

import (
	"github.com/google/uuid"
)

type Category struct {
	ID   uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name string    `gorm:"not null"`
}
