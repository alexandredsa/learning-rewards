package repository

import (
	"context"
	"time"

	"github.com/alexandredsa/learning-rewards/reward-processor/pkg/models"
	"gorm.io/gorm"
)

// UserEventRepository defines the interface for user event count operations
type UserEventRepository interface {
	// IncrementAndGetCount increments the count for a user's event and returns the new count
	IncrementAndGetCount(ctx context.Context, userID, eventType string) (int, error)
	// GetCount returns the current count for a user's event
	GetCount(ctx context.Context, userID, eventType string) (int, error)
}

// GormUserEventRepository implements UserEventRepository using GORM
type GormUserEventRepository struct {
	db *gorm.DB
}

// NewGormUserEventRepository creates a new GORM-based user event count repository
func NewGormUserEventRepository(db *gorm.DB) *GormUserEventRepository {
	return &GormUserEventRepository{db: db}
}

// IncrementAndGetCount implements UserEventRepository
func (r *GormUserEventRepository) IncrementAndGetCount(ctx context.Context, userID, eventType string) (int, error) {
	var count models.UserEventCount

	// Use a transaction to ensure atomicity
	err := r.db.Transaction(func(tx *gorm.DB) error {
		// Try to find existing record
		result := tx.WithContext(ctx).
			Where("user_id = ? AND event_type = ?", userID, eventType).
			First(&count)

		if result.Error == gorm.ErrRecordNotFound {
			// Create new record if not found
			count = models.UserEventCount{
				UserID:    userID,
				EventType: eventType,
				Count:     1,
				UpdatedAt: time.Now(),
			}
			return tx.Create(&count).Error
		}

		if result.Error != nil {
			return result.Error
		}

		// Increment existing record
		count.Count++
		count.UpdatedAt = time.Now()
		return tx.Save(&count).Error
	})

	if err != nil {
		return 0, err
	}

	return count.Count, nil
}

// GetCount implements UserEventRepository
func (r *GormUserEventRepository) GetCount(ctx context.Context, userID, eventType string) (int, error) {
	var count models.UserEventCount

	result := r.db.WithContext(ctx).
		Where("user_id = ? AND event_type = ?", userID, eventType).
		First(&count)

	if result.Error == gorm.ErrRecordNotFound {
		return 0, nil
	}

	if result.Error != nil {
		return 0, result.Error
	}

	return count.Count, nil
}
