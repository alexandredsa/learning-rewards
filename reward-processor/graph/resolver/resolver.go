package resolver

import (
	"github.com/alexandredsa/learning-rewards/reward-processor/internal/repository"
	"go.uber.org/zap"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	RuleRepository repository.RuleRepository
	Logger         *zap.Logger
}

// NewResolver creates a new resolver with the required dependencies
func NewResolver(ruleRepo repository.RuleRepository, logger *zap.Logger) *Resolver {
	return &Resolver{
		RuleRepository: ruleRepo,
		Logger:         logger,
	}
}
