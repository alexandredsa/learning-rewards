package service

import (
	"context"
	"event-processor/internal/messaging/kafka"
	"event-processor/internal/models"
	"time"

	"github.com/google/uuid"
)

type EventService interface {
	ProcessEvent(ctx context.Context, userID, eventType, courseID, category string, timestamp time.Time) error
}

type eventService struct {
	producer *kafka.Producer
}

func NewEventService(producer *kafka.Producer) EventService {
	return &eventService{producer: producer}
}

func (s *eventService) ProcessEvent(ctx context.Context, userID, eventType, courseID, category string, timestamp time.Time) error {
	event := models.LearningEvent{
		ID:        uuid.New(),
		UserID:    userID,
		EventType: eventType,
		CourseID:  courseID,
		Category:  category,
		Timestamp: timestamp,
	}

	return s.producer.PublishEvent(ctx, event)
}
