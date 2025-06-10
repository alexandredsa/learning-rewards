package repository

import (
	"context"
	"time"

	"github.com/alexandredsa/learning-rewards/reward-processor/pkg/models"
	"gorm.io/gorm"
)

// UserEventRepository defines the interface for user event count operations
type UserEventRepository interface {
	// Increment increments the count for a user's event
	Increment(ctx context.Context, userID, eventType, category string) error
	// GetCount returns the current count for a user's event
	// If category is provided, it will count events with that specific category
	// If category is empty, it will count all events of that type using GROUP BY
	GetCount(ctx context.Context, userID, eventType, category string) (int, error)
}

// Ensure GormUserEventRepository implements UserEventRepository
var _ UserEventRepository = (*GormUserEventRepository)(nil)

// GormUserEventRepository implements UserEventRepository using GORM
type GormUserEventRepository struct {
	db *gorm.DB
}

// NewGormUserEventRepository creates a new GORM-based user event count repository
func NewGormUserEventRepository(db *gorm.DB) *GormUserEventRepository {
	return &GormUserEventRepository{db: db}
}

// Increment implements UserEventRepository
func (r *GormUserEventRepository) Increment(ctx context.Context, userID, eventType, category string) error {
	var count models.UserEventCount

	// Use a transaction to ensure atomicity
	err := r.db.Transaction(func(tx *gorm.DB) error {
		// Create or update the event count
		result := tx.WithContext(ctx).
			Where("user_id = ? AND event_type = ? AND category = ?", userID, eventType, category).
			First(&count)

		if result.Error == gorm.ErrRecordNotFound {
			// Create new record
			count = models.UserEventCount{
				UserID:    userID,
				EventType: eventType,
				Category:  category,
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
		return tx.Model(&models.UserEventCount{}).
			Where("user_id = ? AND event_type = ? AND category = ?", userID, eventType, category).
			Updates(map[string]interface{}{
				"count":      count.Count,
				"updated_at": count.UpdatedAt,
			}).Error
	})

	return err
}

// GetCount implements UserEventRepository
func (r *GormUserEventRepository) GetCount(ctx context.Context, userID, eventType, category string) (int, error) {
	var count int64

	query := r.db.WithContext(ctx).Model(&models.UserEventCount{}).
		Where("user_id = ? AND event_type = ?", userID, eventType)

	if category != "" {
		// For category-specific rules, get count for that category
		query = query.Where("category = ?", category)
	} else {
		// For generic rules, sum up all counts for this event type
		query = query.Select("COALESCE(SUM(count), 0)")
	}

	err := query.Count(&count).Error
	if err != nil {
		return 0, err
	}

	return int(count), nil
}
