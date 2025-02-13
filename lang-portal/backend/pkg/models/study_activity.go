package models

import (
	"errors"
	"strings"
	"time"
)

// StudyActivity represents a learning activity in the language portal
type StudyActivity struct {
	ID          int64     `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	Image       string    `json:"image" db:"image"`
	Score       int       `json:"score" db:"score"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

// Validate performs validation checks on the StudyActivity struct
func (sa *StudyActivity) Validate() error {
	// Validate Name
	if sa.Name == "" {
		return errors.New("activity name cannot be empty")
	}

	// Validate Description
	if sa.Description == "" {
		return errors.New("activity description cannot be empty")
	}

	// Validate Score
	if sa.Score < 0 {
		return errors.New("score cannot be negative")
	}

	return nil
}

// Sanitize performs necessary sanitization on the StudyActivity struct
func (sa *StudyActivity) Sanitize() {
	// Trim whitespace from name and description
	sa.Name = trimString(sa.Name)
	sa.Description = trimString(sa.Description)
}

// Helper function to trim whitespace
func trimString(s string) string {
	return strings.TrimSpace(s)
}
