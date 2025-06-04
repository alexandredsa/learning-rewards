package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Event struct {
	ID        uuid.UUID
	UserID    string
	EventType string
	CourseID  string
	Timestamp time.Time
	CreatedAt time.Time
}

type EventRepository interface {
	SaveEvent(ctx context.Context, event Event) error
	GetEventStats(ctx context.Context) ([]EventStat, error)
}

type eventRepo struct {
	db *gorm.DB
}

func NewEventRepository(db *gorm.DB) EventRepository {
	return &eventRepo{db: db}
}

func (r *eventRepo) SaveEvent(ctx context.Context, event Event) error {
	return r.db.WithContext(ctx).Create(&event).Error
}

type EventStat struct {
	EventType string `json:"event_type"`
	Count     int    `json:"count"`
}

func (r *eventRepo) GetEventStats(ctx context.Context) ([]EventStat, error) {
	var stats []EventStat
	err := r.db.WithContext(ctx).
		Table("events").
		Select("event_type, COUNT(*) as count").
		Group("event_type").
		Scan(&stats).Error

	return stats, err
}
