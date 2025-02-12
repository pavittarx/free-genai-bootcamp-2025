package models

import "time"

// Session represents a user's learning session
type Session struct {
	ID            int64     `json:"id" db:"id"`
	StartTime     time.Time `json:"start_time" db:"start_time"`
	EndTime       time.Time `json:"end_time" db:"end_time"`
	TotalDuration int       `json:"total_duration" db:"total_duration"` // Duration in seconds
	WordsLearned  int       `json:"words_learned" db:"words_learned"`
	AverageScore  float64   `json:"average_score" db:"average_score"`
	Status        string    `json:"status" db:"status"`
}

// Validate checks if the session is valid
func (s *Session) Validate() error {
	// Validate session status
	validStatuses := map[string]bool{
		"active": true,
		"completed": true,
		"interrupted": true,
	}

	if !validStatuses[s.Status] {
		s.Status = "active" // Default status
	}

	// Ensure end time is not before start time
	if s.EndTime.Before(s.StartTime) {
		s.EndTime = s.StartTime
	}

	// Calculate total duration
	s.TotalDuration = int(s.EndTime.Sub(s.StartTime).Seconds())

	// Validate words learned
	if s.WordsLearned < 0 {
		s.WordsLearned = 0
	}

	// Validate average score
	if s.AverageScore < 0 {
		s.AverageScore = 0
	} else if s.AverageScore > 100 {
		s.AverageScore = 100
	}

	return nil
}
