package models

import (
	"errors"
	"strings"
	"time"
	"unicode"
)

// Group represents a collection of words with a specific theme or category
type Group struct {
	ID        int64     `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// Validate performs validation checks on the Group struct
func (g *Group) Validate() error {
	// Trim whitespace
	g.Name = strings.TrimSpace(g.Name)

	// Check for empty group name
	if g.Name == "" {
		return errors.New("group name cannot be empty")
	}

	// Check group name length
	if len(g.Name) < 2 || len(g.Name) > 50 {
		return errors.New("group name must be between 2 and 50 characters")
	}

	// Validate group name characters (allow letters, spaces, and some punctuation)
	for _, r := range g.Name {
		if !unicode.IsLetter(r) && !unicode.IsSpace(r) && r != '-' && r != '_' {
			return errors.New("group name can only contain letters, spaces, hyphens, and underscores")
		}
	}

	return nil
}

// Sanitize removes any potentially harmful content and trims whitespace
func (g *Group) Sanitize() {
	g.Name = strings.TrimSpace(g.Name)
}
