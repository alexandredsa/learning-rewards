package resolver

import (
	"encoding/json"

	"github.com/alexandredsa/learning-rewards/reward-processor/graph/model"
	"github.com/alexandredsa/learning-rewards/reward-processor/pkg/models"
)

func convertToGraphQLRule(rule *models.Rule) *model.Rule {
	// Convert conditions to JSON string
	conditionsJSON, _ := json.Marshal(rule.Conditions)

	// Convert count to pointer
	var countPtr *int
	if rule.Count > 0 {
		count := rule.Count
		countPtr = &count
	}

	// Convert reward amount to pointer
	var amountPtr *int
	if rule.Reward.Amount > 0 {
		amount := rule.Reward.Amount
		amountPtr = &amount
	}

	condStr := string(conditionsJSON)
	return &model.Rule{
		ID:         rule.ID,
		Type:       string(rule.Type),
		EventType:  rule.EventType,
		Count:      countPtr,
		Conditions: &condStr,
		Reward: &model.Reward{
			Type:        string(rule.Reward.Type),
			Amount:      amountPtr,
			Description: rule.Reward.Description,
		},
		Enabled: rule.Enabled,
	}
}
