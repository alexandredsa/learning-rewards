package rules

import (
	"context"
	"time"

	"github.com/alexandredsa/learning-rewards/reward-processor/internal/repository"
	"github.com/alexandredsa/learning-rewards/reward-processor/pkg/models"
	"go.uber.org/zap"
)

// Engine handles rule evaluation and milestone tracking
type Engine struct {
	rules     []models.Rule
	eventRepo repository.UserEventRepository
	logger    *zap.Logger
}

// NewEngine creates a new rules engine with the given rules
func NewEngine(rules []models.Rule, eventRepo repository.UserEventRepository, logger *zap.Logger) *Engine {
	return &Engine{
		rules:     rules,
		eventRepo: eventRepo,
		logger:    logger,
	}
}

// SetRules sets the rules for the engine (used for testing)
func (e *Engine) SetRules(rules []models.Rule) {
	e.rules = rules
}

// EvaluateEvent processes a user event against all rules
func (e *Engine) EvaluateEvent(ctx context.Context, event models.UserEvent) ([]models.RewardTriggered, error) {
	var triggered []models.RewardTriggered

	e.logger.Info("Starting event evaluation",
		zap.String("user_id", event.UserID),
		zap.String("event_type", event.EventType),
		zap.Int("total_rules", len(e.rules)))

	for _, rule := range e.rules {
		if !rule.Enabled {
			e.logger.Debug("Skipping disabled rule",
				zap.String("rule_id", rule.ID),
				zap.String("user_id", event.UserID))
			continue
		}

		if rule.EventType != event.EventType {
			e.logger.Debug("Skipping rule due to event type mismatch",
				zap.String("rule_id", rule.ID),
				zap.String("rule_event_type", rule.EventType),
				zap.String("event_type", event.EventType),
				zap.String("user_id", event.UserID))
			continue
		}

		// Check conditions
		if !e.matchesConditions(event, rule.Conditions) {
			e.logger.Debug("Rule conditions not met",
				zap.String("rule_id", rule.ID),
				zap.String("user_id", event.UserID),
				zap.Any("conditions", rule.Conditions),
				zap.Any("event_data", event))
			continue
		}

		e.logger.Debug("Rule conditions met, evaluating rule type",
			zap.String("rule_id", rule.ID),
			zap.String("rule_type", string(rule.Type)),
			zap.String("user_id", event.UserID))

		switch rule.Type {
		case models.SingleEventRule:
			triggered = append(triggered, models.RewardTriggered{
				UserID:    event.UserID,
				RuleID:    rule.ID,
				Reward:    rule.Reward,
				Timestamp: time.Now(),
			})
			e.logger.Info("Single event rule triggered",
				zap.String("user_id", event.UserID),
				zap.String("rule_id", rule.ID),
				zap.String("event", event.EventType),
				zap.Any("reward", rule.Reward))

		case models.MilestoneRule:
			count, err := e.eventRepo.IncrementAndGetCount(ctx, event.UserID, event.EventType)
			if err != nil {
				e.logger.Error("Failed to increment milestone count",
					zap.String("user_id", event.UserID),
					zap.String("rule_id", rule.ID),
					zap.String("event_type", event.EventType),
					zap.Error(err))
				return nil, err
			}

			e.logger.Debug("Current milestone count",
				zap.String("user_id", event.UserID),
				zap.String("rule_id", rule.ID),
				zap.Int("current_count", count),
				zap.Int("required_count", rule.Count))

			if count == rule.Count {
				triggered = append(triggered, models.RewardTriggered{
					UserID:    event.UserID,
					RuleID:    rule.ID,
					Reward:    rule.Reward,
					Timestamp: time.Now(),
				})
				e.logger.Info("Milestone rule triggered",
					zap.String("user_id", event.UserID),
					zap.String("rule_id", rule.ID),
					zap.String("event", event.EventType),
					zap.Int("count", rule.Count),
					zap.Any("reward", rule.Reward))
			}
		}
	}

	e.logger.Info("Completed event evaluation",
		zap.String("user_id", event.UserID),
		zap.String("event_type", event.EventType),
		zap.Int("rules_triggered", len(triggered)))

	return triggered, nil
}

// matchesConditions checks if an event matches all conditions in a rule
func (e *Engine) matchesConditions(event models.UserEvent, conditions map[string]string) bool {
	for key, value := range conditions {
		switch key {
		case "category":
			if event.Category != value {
				return false
			}
			// Add more condition types here as needed
		}
	}
	return true
}

// GetMilestoneCount returns the current count for a user's milestone
func (e *Engine) GetMilestoneCount(ctx context.Context, userID, eventType string) (int, error) {
	return e.eventRepo.GetCount(ctx, userID, eventType)
}
