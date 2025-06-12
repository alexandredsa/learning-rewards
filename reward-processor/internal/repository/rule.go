package repository

import (
	"context"

	"github.com/alexandredsa/learning-rewards/reward-processor/pkg/models"
	"gorm.io/gorm"
)

// RuleRepository defines the interface for rule operations
type RuleRepository interface {
	// GetEnabledRules returns all enabled rules
	GetEnabledRules(ctx context.Context) ([]models.Rule, error)
}

// Ensure GormRuleRepository implements RuleRepository
var _ RuleRepository = (*GormRuleRepository)(nil)

// GormRuleRepository implements RuleRepository using GORM
type GormRuleRepository struct {
	db *gorm.DB
}

// NewGormRuleRepository creates a new GORM-based rule repository
func NewGormRuleRepository(db *gorm.DB) *GormRuleRepository {
	return &GormRuleRepository{db: db}
}

// GetEnabledRules implements RuleRepository
func (r *GormRuleRepository) GetEnabledRules(ctx context.Context) ([]models.Rule, error) {
	var rules []models.Rule
	err := r.db.WithContext(ctx).
		Where("enabled = ?", true).
		Find(&rules).Error
	return rules, err
}
