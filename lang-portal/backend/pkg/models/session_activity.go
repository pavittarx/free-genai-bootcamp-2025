package models

import (
	"time"
)

// SessionActivity represents an individual activity within a learning session
type SessionActivity struct {
	ID         int64     `json:"id" db:"id"`
	SessionID  int64     `json:"session_id" db:"session_id"`
	ActivityID int64     `json:"activity_id" db:"activity_id"`
	Answer     string    `json:"answer" db:"answer"`
	Result     string    `json:"result" db:"result"`
	Score      int       `json:"score" db:"score"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}

// Validate performs validation checks on the SessionActivity struct
func (sa *SessionActivity) Validate() error {
	// Validate required fields
	if sa.SessionID <= 0 {
		return ErrInvalidID
	}

	if sa.ActivityID <= 0 {
		return ErrInvalidID
	}

	// Score validation (can be 0 or positive)
	if sa.Score < 0 {
		return ErrInvalidScore
	}

	return nil
}

// IsSuccessful checks if the session activity was completed successfully
func (sa *SessionActivity) IsSuccessful() bool {
	return sa.Result == "success"
}
