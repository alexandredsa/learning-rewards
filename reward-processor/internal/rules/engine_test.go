package rules_test

import (
	"context"
	"testing"
	"time"

	"github.com/alexandredsa/learning-rewards/reward-processor/internal/rules"
	"github.com/alexandredsa/learning-rewards/reward-processor/pkg/models"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

// stubUserEventRepository is a simple stub implementation
type stubUserEventRepository struct {
	getCount int
	err      error
}

func (s *stubUserEventRepository) Increment(ctx context.Context, userID, eventType, category string) error {
	return s.err
}

func (s *stubUserEventRepository) GetCount(ctx context.Context, userID, eventType, category string) (int, error) {
	return s.getCount, s.err
}

func TestEvaluateEvent_SingleEventRule(t *testing.T) {
	logger, err := zap.NewDevelopment()
	assert.NoError(t, err)

	// Create engine with a stub that won't be used for single event rules
	engine := rules.NewEngine([]models.Rule{}, &stubUserEventRepository{}, logger)

	// Define a single event rule
	rule := models.Rule{
		ID:        "rule-001",
		EventType: "COURSE_COMPLETED",
		Conditions: map[string]string{
			"category": "MATH",
		},
		Reward: models.Reward{
			Type:        models.BadgeReward,
			Description: "Math Course Completed",
		},
		Enabled: true,
	}
	engine.SetRules([]models.Rule{rule})

	// Test cases
	tests := []struct {
		name          string
		event         models.UserEvent
		expectedCount int
	}{
		{
			name: "matching event triggers reward",
			event: models.UserEvent{
				UserID:    "user-001",
				EventType: "COURSE_COMPLETED",
				Category:  "MATH",
				Timestamp: time.Now(),
			},
			expectedCount: 1,
		},
		{
			name: "non-matching event type doesn't trigger",
			event: models.UserEvent{
				UserID:    "user-001",
				EventType: "COURSE_STARTED",
				Category:  "MATH",
				Timestamp: time.Now(),
			},
			expectedCount: 0,
		},
		{
			name: "non-matching category doesn't trigger",
			event: models.UserEvent{
				UserID:    "user-001",
				EventType: "COURSE_COMPLETED",
				Category:  "SCIENCE",
				Timestamp: time.Now(),
			},
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			triggered, err := engine.EvaluateEvent(context.Background(), tt.event)
			assert.NoError(t, err)
			assert.Len(t, triggered, tt.expectedCount)

			if tt.expectedCount > 0 {
				assert.Equal(t, rule.ID, triggered[0].RuleID)
				assert.Equal(t, tt.event.UserID, triggered[0].UserID)
				assert.Equal(t, rule.Reward, triggered[0].Reward)
			}
		})
	}
}

func TestEvaluateEvent_MilestoneRule(t *testing.T) {
	logger, err := zap.NewDevelopment()
	assert.NoError(t, err)

	// Define a milestone rule
	rule := models.Rule{
		ID:        "rule-002",
		EventType: "COURSE_COMPLETED",
		Count:     3,
		Conditions: map[string]string{
			"category": "MATH",
		},
		Reward: models.Reward{
			Type:        models.PointsReward,
			Amount:      100,
			Description: "Math Course Milestone",
		},
		Enabled: true,
	}

	// Test cases
	tests := []struct {
		name          string
		event         models.UserEvent
		stubCount     int
		expectedCount int
	}{
		{
			name: "milestone reached triggers reward",
			event: models.UserEvent{
				UserID:    "user-001",
				EventType: "COURSE_COMPLETED",
				Category:  "MATH",
				Timestamp: time.Now(),
			},
			stubCount:     3,
			expectedCount: 1,
		},
		{
			name: "milestone not reached doesn't trigger",
			event: models.UserEvent{
				UserID:    "user-001",
				EventType: "COURSE_COMPLETED",
				Category:  "MATH",
				Timestamp: time.Now(),
			},
			stubCount:     2,
			expectedCount: 0,
		},
		{
			name: "non-matching event type doesn't trigger",
			event: models.UserEvent{
				UserID:    "user-001",
				EventType: "COURSE_STARTED",
				Category:  "MATH",
				Timestamp: time.Now(),
			},
			stubCount:     3,
			expectedCount: 0,
		},
		{
			name: "non-matching category doesn't trigger",
			event: models.UserEvent{
				UserID:    "user-001",
				EventType: "COURSE_COMPLETED",
				Category:  "SCIENCE",
				Timestamp: time.Now(),
			},
			stubCount:     3,
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new stub for each test case
			stubRepo := &stubUserEventRepository{getCount: tt.stubCount}
			engine := rules.NewEngine([]models.Rule{rule}, stubRepo, logger)

			triggered, err := engine.EvaluateEvent(context.Background(), tt.event)
			assert.NoError(t, err)
			assert.Len(t, triggered, tt.expectedCount)

			if tt.expectedCount > 0 {
				assert.Equal(t, rule.ID, triggered[0].RuleID)
				assert.Equal(t, tt.event.UserID, triggered[0].UserID)
				assert.Equal(t, rule.Reward, triggered[0].Reward)
			}
		})
	}
}

func TestEvaluateEvent_DisabledRule(t *testing.T) {
	logger, err := zap.NewDevelopment()
	assert.NoError(t, err)

	// Define a disabled rule
	rule := models.Rule{
		ID:        "rule-003",
		EventType: "COURSE_COMPLETED",
		Conditions: map[string]string{
			"category": "MATH",
		},
		Reward: models.Reward{
			Type:        models.BadgeReward,
			Description: "Disabled Rule",
		},
		Enabled: false,
	}

	engine := rules.NewEngine([]models.Rule{rule}, &stubUserEventRepository{}, logger)

	event := models.UserEvent{
		UserID:    "user-001",
		EventType: "COURSE_COMPLETED",
		Category:  "MATH",
		Timestamp: time.Now(),
	}

	triggered, err := engine.EvaluateEvent(context.Background(), event)
	assert.NoError(t, err)
	assert.Empty(t, triggered)
}

func TestEvaluateEvent_RepositoryError(t *testing.T) {
	logger, err := zap.NewDevelopment()
	assert.NoError(t, err)

	// Define a milestone rule
	rule := models.Rule{
		ID:        "rule-004",
		EventType: "COURSE_COMPLETED",
		Count:     3,
		Conditions: map[string]string{
			"category": "MATH",
		},
		Reward: models.Reward{
			Type:        models.PointsReward,
			Amount:      100,
			Description: "Math Course Milestone",
		},
		Enabled: true,
	}

	// Create a stub that returns an error
	stubRepo := &stubUserEventRepository{err: assert.AnError}
	engine := rules.NewEngine([]models.Rule{rule}, stubRepo, logger)

	event := models.UserEvent{
		UserID:    "user-001",
		EventType: "COURSE_COMPLETED",
		Category:  "MATH",
		Timestamp: time.Now(),
	}

	triggered, err := engine.EvaluateEvent(context.Background(), event)
	assert.Error(t, err)
	assert.Empty(t, triggered)
}

func TestGetMilestoneCount(t *testing.T) {
	logger, err := zap.NewDevelopment()
	assert.NoError(t, err)

	expectedCount := 5
	stubRepo := &stubUserEventRepository{getCount: expectedCount}
	engine := rules.NewEngine([]models.Rule{}, stubRepo, logger)

	userID := "user-001"
	eventType := "COURSE_COMPLETED"
	category := "MATH"

	count, err := engine.GetMilestoneCount(context.Background(), userID, eventType, category)
	assert.NoError(t, err)
	assert.Equal(t, expectedCount, count)
}
