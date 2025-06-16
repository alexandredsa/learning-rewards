package resolver_test

import (
	"context"
	"testing"

	"github.com/alexandredsa/learning-rewards/reward-processor/graph/model"
	"github.com/alexandredsa/learning-rewards/reward-processor/graph/resolver"
	"github.com/alexandredsa/learning-rewards/reward-processor/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

// MockRuleRepository is a mock implementation of repository.RuleRepository
type MockRuleRepository struct {
	mock.Mock
}

func (m *MockRuleRepository) CreateRule(ctx context.Context, rule *models.Rule) error {
	args := m.Called(ctx, rule)
	return args.Error(0)
}

func (m *MockRuleRepository) UpdateRule(ctx context.Context, id string, rule *models.Rule) error {
	args := m.Called(ctx, id, rule)
	return args.Error(0)
}

func (m *MockRuleRepository) GetRuleByID(ctx context.Context, id string) (*models.Rule, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Rule), args.Error(1)
}

func (m *MockRuleRepository) GetEnabledRules(ctx context.Context) ([]models.Rule, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.Rule), args.Error(1)
}

// TestCase represents a test case with setup and assertions
type TestCase struct {
	name         string
	setupMocks   func(*MockRuleRepository)
	runTest      func(*resolver.Resolver) (interface{}, error)
	assertResult func(*testing.T, interface{}, error)
	assertMocks  func(*testing.T, *MockRuleRepository)
}

func setupTestResolver(t *testing.T) (*resolver.Resolver, *MockRuleRepository) {
	mockRepo := new(MockRuleRepository)
	logger, _ := zap.NewDevelopment()
	resolver := resolver.NewResolver(mockRepo, logger)
	return resolver, mockRepo
}

func runTestCase(t *testing.T, tc TestCase) {
	resolver, mockRepo := setupTestResolver(t)

	// Setup mocks
	tc.setupMocks(mockRepo)

	// Run test
	result, err := tc.runTest(resolver)

	// Assert results
	tc.assertResult(t, result, err)

	// Assert mock expectations
	tc.assertMocks(t, mockRepo)
}

func TestCreateRule(t *testing.T) {
	tests := []TestCase{
		{
			name: "create rule with conditions",
			setupMocks: func(m *MockRuleRepository) {
				m.On("CreateRule", mock.Anything, &models.Rule{
					EventType:          "COURSE_COMPLETED",
					Count:              5,
					ConditionsCategory: ptrString("MATH"),
					Reward: models.Reward{
						Type:        models.RewardType("POINTS"),
						Amount:      100,
						Description: "Completed 5 math courses",
					},
					Enabled: true,
				}).Return(nil)
			},
			runTest: func(r *resolver.Resolver) (interface{}, error) {
				return r.Mutation().CreateRule(context.Background(), model.CreateRuleInput{
					EventType: "COURSE_COMPLETED",
					Count:     ptrInt(5),
					Conditions: &model.RuleConditionsInput{
						Category: ptrString("MATH"),
					},
					Reward: &model.RewardInput{
						Type:        model.RewardType("POINTS"),
						Amount:      ptrInt(100),
						Description: "Completed 5 math courses",
					},
					Enabled: true,
				})
			},
			assertResult: func(t *testing.T, result interface{}, err error) {
				assert.NoError(t, err)
				rule := result.(*model.Rule)
				assert.Equal(t, "COURSE_COMPLETED", rule.EventType)
				assert.Equal(t, ptrInt(5), rule.Count)
				assert.Equal(t, ptrString("MATH"), rule.Conditions.Category)
				assert.Equal(t, model.RewardType("POINTS"), rule.Reward.Type)
				assert.Equal(t, ptrInt(100), rule.Reward.Amount)
				assert.Equal(t, "Completed 5 math courses", rule.Reward.Description)
				assert.True(t, rule.Enabled)
			},
			assertMocks: func(t *testing.T, m *MockRuleRepository) {
				m.AssertExpectations(t)
			},
		},
		{
			name: "create rule without conditions (count should default to '1')",
			setupMocks: func(m *MockRuleRepository) {
				m.On("CreateRule", mock.Anything, &models.Rule{
					EventType: "COURSE_COMPLETED",
					Count:     1,
					Reward: models.Reward{
						Type:        models.RewardType("BADGE"),
						Description: "Completed a course",
					},
					Enabled: true,
				}).Return(nil)
			},
			runTest: func(r *resolver.Resolver) (interface{}, error) {
				return r.Mutation().CreateRule(context.Background(), model.CreateRuleInput{
					EventType: "COURSE_COMPLETED",
					Reward: &model.RewardInput{
						Type:        model.RewardType("BADGE"),
						Description: "Completed a course",
					},
					Enabled: true,
				})
			},
			assertResult: func(t *testing.T, result interface{}, err error) {
				assert.NoError(t, err)
				rule := result.(*model.Rule)
				assert.Equal(t, "COURSE_COMPLETED", rule.EventType)
				assert.Equal(t, *rule.Count, 1)
				assert.Nil(t, rule.Conditions)
				assert.Equal(t, model.RewardType("BADGE"), rule.Reward.Type)
				assert.Nil(t, rule.Reward.Amount)
				assert.Equal(t, "Completed a course", rule.Reward.Description)
				assert.True(t, rule.Enabled)
			},
			assertMocks: func(t *testing.T, m *MockRuleRepository) {
				m.AssertExpectations(t)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			runTestCase(t, tc)
		})
	}
}

func TestUpdateRule(t *testing.T) {
	tests := []TestCase{
		{
			name: "update rule conditions",
			setupMocks: func(m *MockRuleRepository) {
				existingRule := &models.Rule{
					ID:                 "rule-001",
					EventType:          "COURSE_COMPLETED",
					ConditionsCategory: ptrString("MATH"),
					Reward: models.Reward{
						Type:        models.RewardType("POINTS"),
						Amount:      100,
						Description: "Math course reward",
					},
					Enabled: true,
				}
				updatedRule := &models.Rule{
					ID:                 "rule-001",
					EventType:          "COURSE_COMPLETED",
					ConditionsCategory: ptrString("SCIENCE"),
					Reward: models.Reward{
						Type:        models.RewardType("POINTS"),
						Amount:      100,
						Description: "Math course reward",
					},
					Enabled: true,
				}

				m.On("GetRuleByID", mock.Anything, "rule-001").Return(existingRule, nil)
				m.On("UpdateRule", mock.Anything, "rule-001", mock.MatchedBy(func(rule *models.Rule) bool {
					return rule.ConditionsCategory != nil && *rule.ConditionsCategory == "SCIENCE"
				})).Return(nil)
				m.On("GetRuleByID", mock.Anything, "rule-001").Return(updatedRule, nil)
			},
			runTest: func(r *resolver.Resolver) (interface{}, error) {
				return r.Mutation().UpdateRule(context.Background(), "rule-001", model.UpdateRuleInput{
					Conditions: &model.RuleConditionsInput{
						Category: ptrString("SCIENCE"),
					},
				})
			},
			assertResult: func(t *testing.T, result interface{}, err error) {
				assert.NoError(t, err)
				rule := result.(*model.Rule)
				assert.Equal(t, "rule-001", rule.ID)
				assert.Equal(t, ptrString("SCIENCE"), rule.Conditions.Category)
			},
			assertMocks: func(t *testing.T, m *MockRuleRepository) {
				m.AssertExpectations(t)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			runTestCase(t, tc)
		})
	}
}

func TestConvertToGraphQLRule(t *testing.T) {
	tests := []struct {
		name     string
		input    *models.Rule
		expected *model.Rule
	}{
		{
			name: "convert rule with conditions",
			input: &models.Rule{
				ID:                 "rule-001",
				EventType:          "COURSE_COMPLETED",
				Count:              5,
				ConditionsCategory: ptrString("MATH"),
				Reward: models.Reward{
					Type:        models.RewardType("POINTS"),
					Amount:      100,
					Description: "Math course reward",
				},
				Enabled: true,
			},
			expected: &model.Rule{
				ID:        "rule-001",
				EventType: "COURSE_COMPLETED",
				Count:     ptrInt(5),
				Conditions: &model.RuleConditions{
					Category: ptrString("MATH"),
				},
				Reward: &model.Reward{
					Type:        model.RewardType("POINTS"),
					Amount:      ptrInt(100),
					Description: "Math course reward",
				},
				Enabled: true,
			},
		},
		{
			name: "convert rule without conditions",
			input: &models.Rule{
				ID:        "rule-002",
				EventType: "COURSE_COMPLETED",
				Reward: models.Reward{
					Type:        models.RewardType("BADGE"),
					Description: "Course completion badge",
				},
				Enabled: true,
			},
			expected: &model.Rule{
				ID:         "rule-002",
				EventType:  "COURSE_COMPLETED",
				Conditions: nil,
				Reward: &model.Reward{
					Type:        model.RewardType("BADGE"),
					Description: "Course completion badge",
				},
				Enabled: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := resolver.ConvertToGraphQLRule(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Helper function to create a pointer to an int
func ptrInt(i int) *int {
	return &i
}

// Helper function to create a pointer to a string
func ptrString(s string) *string {
	return &s
}
