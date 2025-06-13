package resolver

import (
	"github.com/alexandredsa/learning-rewards/reward-processor/internal/repository"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	RuleRepository repository.RuleRepository
}

// NewResolver creates a new resolver with the required dependencies
func NewResolver(ruleRepo repository.RuleRepository) *Resolver {
	return &Resolver{
		RuleRepository: ruleRepo,
	}
}
