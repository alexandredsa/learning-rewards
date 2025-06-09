package models

import (
	"time"

	"github.com/google/uuid"
)

type LearningEvent struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	UserID    string    `gorm:"not null"`
	EventType string    `gorm:"not null"`
	CourseID  string    `gorm:"not null"`
	Timestamp time.Time `gorm:"not null"`
	CreatedAt time.Time
}
