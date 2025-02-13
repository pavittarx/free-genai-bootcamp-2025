package models

import (
	"time"
)

// Session represents a learning session in the language portal
type Session struct {
	ID         int64      `json:"id" db:"id"`
	ActivityID int64      `json:"activity_id" db:"activity_id"`
	GroupID    *int64     `json:"group_id,omitempty" db:"group_id"`
	StartTime  time.Time  `json:"start_time" db:"start_time"`
	EndTime    *time.Time `json:"end_time,omitempty" db:"end_time"`
	Score      int        `json:"score" db:"score"`
	CreatedAt  time.Time  `json:"created_at" db:"created_at"`
}

// Validate performs validation checks on the Session struct
func (s *Session) Validate() error {
	// Validate required fields
	if s.ActivityID <= 0 {
		return ErrInvalidID
	}

	// Optional group validation
	if s.GroupID != nil && *s.GroupID <= 0 {
		return ErrInvalidID
	}

	// Validate time constraints
	if s.StartTime.IsZero() {
		return ErrInvalidTime
	}

	// EndTime can be nil, but if set, it should be after StartTime
	if s.EndTime != nil && s.EndTime.Before(s.StartTime) {
		return ErrInvalidTimeRange
	}

	// Score validation (can be 0 or positive)
	if s.Score < 0 {
		return ErrInvalidScore
	}

	return nil
}

// Duration calculates the duration of the session
func (s *Session) Duration() time.Duration {
	if s.EndTime == nil {
		return 0
	}
	return s.EndTime.Sub(s.StartTime)
}

// IsCompleted checks if the session has been completed
func (s *Session) IsCompleted() bool {
	return s.EndTime != nil
}
