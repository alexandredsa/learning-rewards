package service

import (
	"context"
	"event-processor/internal/repository"
	"time"

	"github.com/google/uuid"
)

type EventService interface {
	ProcessEvent(ctx context.Context, userID, eventType, courseID string, timestamp time.Time) error
}

type eventService struct {
	repo repository.EventRepository
}

func NewEventService(repo repository.EventRepository) EventService {
	return &eventService{repo: repo}
}

func (s *eventService) ProcessEvent(ctx context.Context, userID, eventType, courseID string, timestamp time.Time) error {
	event := repository.Event{
		ID:        uuid.New(),
		UserID:    userID,
		EventType: eventType,
		CourseID:  courseID,
		Timestamp: timestamp,
	}
	return s.repo.SaveEvent(ctx, event)
}
