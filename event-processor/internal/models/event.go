package models

import (
	"time"

	"github.com/google/uuid"
)

type LearningEvent struct {
	ID        uuid.UUID `json:"id"`
	UserID    string    `json:"user_id"`
	EventType string    `json:"event_type"`
	Category  string    `json:"category"`
	CourseID  string    `json:"course_id"`
	Timestamp time.Time `json:"timestamp"`
	CreatedAt time.Time `json:"created_at"`
}
