package models

import "time"

// StudyActivity represents a single learning activity in a session
type StudyActivity struct {
	ID          int64     `json:"id" db:"id"`
	SessionID   int64     `json:"session_id" db:"session_id"`
	WordID      int64     `json:"word_id" db:"word_id"`
	ActivityType string   `json:"activity_type" db:"activity_type"`
	Timestamp   time.Time `json:"timestamp" db:"timestamp"`
	Duration    int       `json:"duration" db:"duration"` // Duration in seconds
	Score       float64   `json:"score" db:"score"`
}

// Validate checks if the study activity is valid
func (sa *StudyActivity) Validate() error {
	// Validate activity type
	validActivityTypes := map[string]bool{
		"translation": true,
		"pronunciation": true,
		"listening": true,
		"writing": true,
	}

	if !validActivityTypes[sa.ActivityType] {
		sa.ActivityType = "translation" // Default activity type
	}

	// Ensure positive duration
	if sa.Duration < 0 {
		sa.Duration = 0
	}

	// Validate score range
	if sa.Score < 0 {
		sa.Score = 0
	} else if sa.Score > 100 {
		sa.Score = 100
	}

	return nil
}
