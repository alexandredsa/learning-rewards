package resolver

import (
	"github.com/alexandredsa/learning-rewards/reward-processor/graph/model"
	"github.com/alexandredsa/learning-rewards/reward-processor/pkg/models"
)

const (
	minCountValue = 1
)

func ConvertToGraphQLRule(rule *models.Rule) *model.Rule {
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

	var conditions *model.RuleConditions
	if rule.ConditionsCategory != nil {
		conditions = &model.RuleConditions{
			Category: rule.ConditionsCategory,
		}
	}

	return &model.Rule{
		ID:         rule.ID,
		EventType:  rule.EventType,
		Count:      countPtr,
		Conditions: conditions,
		Reward: &model.Reward{
			Type:        model.RewardType(rule.Reward.Type),
			Amount:      amountPtr,
			Description: rule.Reward.Description,
		},
		Enabled: rule.Enabled,
	}
}

func ConvertGraphQLRuleToModel(rule interface{}) *models.Rule {
	switch r := rule.(type) {
	case *model.CreateRuleInput:
		// Convert count from pointer to value
		count := 0
		if r.Count != nil {
			count = *r.Count
		} else {
			count = minCountValue
		}

		// Convert conditions
		var conditionsCategory *string
		if r.Conditions != nil && r.Conditions.Category != nil {
			conditionsCategory = r.Conditions.Category
		}

		// Convert reward amount from pointer to value
		rewardAmount := 0
		if r.Reward.Amount != nil {
			rewardAmount = *r.Reward.Amount
		}

		return &models.Rule{
			EventType:          r.EventType,
			Count:              count,
			ConditionsCategory: conditionsCategory,
			Reward: models.Reward{
				Type:        models.RewardType(r.Reward.Type),
				Amount:      rewardAmount,
				Description: r.Reward.Description,
			},
			Enabled: r.Enabled,
		}

	case *model.UpdateRuleInput:
		rule := &models.Rule{}

		if r.EventType != nil {
			rule.EventType = *r.EventType
		}
		if r.Count != nil {
			rule.Count = *r.Count
		}
		if r.Conditions != nil && r.Conditions.Category != nil {
			rule.ConditionsCategory = r.Conditions.Category
		}
		if r.Reward != nil {
			if r.Reward.Type != "" {
				rule.Reward.Type = models.RewardType(r.Reward.Type)
			}
			if r.Reward.Amount != nil {
				rule.Reward.Amount = *r.Reward.Amount
			}
			if r.Reward.Description != "" {
				rule.Reward.Description = r.Reward.Description
			}
		}
		if r.Enabled != nil {
			rule.Enabled = *r.Enabled
		}

		return rule

	default:
		return nil
	}
}
