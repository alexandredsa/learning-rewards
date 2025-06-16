package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// RewardType represents the type of reward
type RewardType string

const (
	BadgeReward  RewardType = "BADGE"
	PointsReward RewardType = "POINTS"
)

// Rule represents a reward rule
type Rule struct {
	ID         string         `json:"id" gorm:"primaryKey"`
	EventType  string         `json:"event_type"`
	Count      int            `json:"count,omitempty"`
	Conditions RuleConditions `json:"conditions" gorm:"type:jsonb"`
	Reward     Reward         `json:"reward" gorm:"embedded"`
	Enabled    bool           `json:"enabled"`
}

// RuleConditions is a custom type to handle map[string]string in GORM
type RuleConditions map[string]string

// Value implements the driver.Valuer interface for RuleConditions
func (rc RuleConditions) Value() (driver.Value, error) {
	if rc == nil {
		return nil, nil
	}
	return json.Marshal(rc)
}

// Scan implements the sql.Scanner interface for RuleConditions
func (rc *RuleConditions) Scan(value interface{}) error {
	if value == nil {
		*rc = make(RuleConditions)
		return nil
	}

	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, rc)
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
