package models

import (
	"errors"
)

// Group represents a category of words
type Group struct {
	ID          int64  `json:"id" db:"id"`
	Name        string `json:"name" db:"name"`
	Description string `json:"description" db:"description"`
}

// GroupFilter allows filtering and searching groups
type GroupFilter struct {
	Name string
}

// Validate checks if the group meets basic validation criteria
func (g *Group) Validate() error {
	// Ensure name is not empty
	if g.Name == "" {
		return errors.New("group name cannot be empty")
	}

	// Optionally, add more validation rules
	if len(g.Name) > 100 {
		return errors.New("group name cannot exceed 100 characters")
	}

	return nil
}

// Custom errors
var (
	ErrInvalidGroup = &GroupError{message: "invalid group: name is required"}
)

// GroupError represents a custom error for group-related operations
type GroupError struct {
	message string
}

func (e *GroupError) Error() string {
	return e.message
}
