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
	// GetRuleByID returns a rule by its ID
	GetRuleByID(ctx context.Context, id string) (*models.Rule, error)
	// CreateRule creates a new rule
	CreateRule(ctx context.Context, rule *models.Rule) error
	// UpdateRule updates an existing rule
	UpdateRule(ctx context.Context, id string, rule *models.Rule) error
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

// GetRuleByID implements RuleRepository
func (r *GormRuleRepository) GetRuleByID(ctx context.Context, id string) (*models.Rule, error) {
	var rule models.Rule
	err := r.db.WithContext(ctx).First(&rule, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &rule, nil
}

// CreateRule implements RuleRepository
func (r *GormRuleRepository) CreateRule(ctx context.Context, rule *models.Rule) error {
	return r.db.WithContext(ctx).Create(rule).Error
}

// UpdateRule implements RuleRepository
func (r *GormRuleRepository) UpdateRule(ctx context.Context, id string, rule *models.Rule) error {
	result := r.db.WithContext(ctx).Model(&models.Rule{}).Where("id = ?", id).Updates(rule)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
