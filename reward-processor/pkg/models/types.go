package models

import "time"

// RuleType represents the type of rule
type RuleType string

// RewardType represents the type of reward
type RewardType string

const (
	BadgeReward  RewardType = "BADGE"
	PointsReward RewardType = "POINTS"
)

// Rule represents a reward rule
type Rule struct {
	ID         string            `json:"id"`
	Type       RuleType          `json:"type"`
	EventType  string            `json:"event_type"`
	Count      int               `json:"count,omitempty"`
	Conditions map[string]string `json:"conditions"`
	Reward     Reward            `json:"reward"`
	Enabled    bool              `json:"enabled"`
}

// Reward represents a reward definition
type Reward struct {
	Type        RewardType `json:"type"`
	Amount      int        `json:"amount,omitempty"` // Only for POINTS rewards
	Description string     `json:"description"`
}

// UserEvent represents an incoming user event
type UserEvent struct {
	UserID    string    `json:"user_id"`
	EventType string    `json:"event_type"`
	CourseID  string    `json:"course_id"`
	Category  string    `json:"category"`
	Timestamp time.Time `json:"timestamp"`
}

// RewardTriggered represents a triggered reward event
type RewardTriggered struct {
	UserID    string    `json:"user_id"`
	RuleID    string    `json:"rule_id"`
	Reward    Reward    `json:"reward"`
	Timestamp time.Time `json:"timestamp"`
}

// UserEventCount represents a user's event count in the database
type UserEventCount struct {
	UserID    string    `json:"user_id" db:"user_id"`
	EventType string    `json:"event_type" db:"event_type"`
	Category  string    `json:"category" db:"category"`
	Count     int       `json:"count" db:"count"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
