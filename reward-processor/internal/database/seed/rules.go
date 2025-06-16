package seed

import (
	"context"
	"os"

	"github.com/alexandredsa/learning-rewards/reward-processor/internal/repository"
	"github.com/alexandredsa/learning-rewards/reward-processor/pkg/models"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// SeedRules seeds the database with initial rules if none exist
// Only runs in non-production environments
func SeedRules(ctx context.Context, db *gorm.DB, log *zap.Logger) error {
	// Skip seeding in production
	if os.Getenv("ENV") == "production" {
		log.Info("Skipping rule seeding in production environment")
		return nil
	}

	ruleRepo := repository.NewGormRuleRepository(db)

	// Check if we already have rules
	rules, err := ruleRepo.GetEnabledRules(ctx)
	if err != nil {
		return err
	}

	// If we already have rules, skip seeding
	if len(rules) > 0 {
		log.Info("Rules already exist in database, skipping seed")
		return nil
	}

	var (
		categoryMath        string = "MATH"
		categoryProgramming string = "PROGRAMMING"
	)

	// Define initial rules
	initialRules := []models.Rule{
		{
			ID:                 "rule-001",
			EventType:          "COURSE_COMPLETED",
			ConditionsCategory: &categoryMath,
			Count:              1,
			Reward: models.Reward{
				Type:        models.BadgeReward,
				Description: "Finished a Math course",
			},
			Enabled: true,
		},
		{
			ID:                 "rule-002",
			EventType:          "COURSE_COMPLETED",
			Count:              5,
			ConditionsCategory: &categoryMath,
			Reward: models.Reward{
				Type:        models.PointsReward,
				Amount:      100,
				Description: "Completed 5 math courses",
			},
			Enabled: true,
		},
		{
			ID:        "rule-003",
			EventType: "COURSE_COMPLETED",
			Count:     30,
			Reward: models.Reward{
				Type:        models.PointsReward,
				Amount:      30,
				Description: "Completed 30 courses",
			},
			Enabled: true,
		},
		{
			ID:        "rule-004",
			EventType: "CHAPTER_COMPLETED",
			Count:     10,
			Reward: models.Reward{
				Type:        models.PointsReward,
				Amount:      10,
				Description: "Completed 10 chapters",
			},
			Enabled: true,
		},
		{
			ID:                 "rule-005",
			EventType:          "COURSE_COMPLETED",
			Count:              1,
			ConditionsCategory: &categoryProgramming,
			Reward: models.Reward{
				Type:        models.BadgeReward,
				Description: "Finished a Programming course",
			},
			Enabled: true,
		},
		{
			ID:                 "rule-006",
			EventType:          "COURSE_COMPLETED",
			Count:              5,
			ConditionsCategory: &categoryProgramming,
			Reward: models.Reward{
				Type:        models.PointsReward,
				Amount:      150,
				Description: "Completed 5 programming courses",
			},
			Enabled: true,
		},
	}

	// Insert rules in a transaction
	err = db.Transaction(func(tx *gorm.DB) error {
		for _, rule := range initialRules {
			if err := tx.Create(&rule).Error; err != nil {
				return err
			}
			log.Info("Seeded rule", zap.String("rule_id", rule.ID), zap.String("description", rule.Reward.Description))
		}
		return nil
	})

	if err != nil {
		return err
	}

	log.Info("Successfully seeded rules", zap.Int("count", len(initialRules)))
	return nil
}
