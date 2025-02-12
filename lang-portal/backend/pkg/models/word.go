package models

import (
	"errors"
)

// Word represents a multilingual word entry
type Word struct {
	ID          int64  `json:"id" db:"id"`
	Hindi       string `json:"hindi" db:"hindi"`
	Scrambled   string `json:"scrambled" db:"scrambled"`
	Hinglish    string `json:"hinglish" db:"hinglish"`
	English     string `json:"english" db:"english"`
	Difficulty  string `json:"difficulty" db:"difficulty"`
	Description string `json:"description,omitempty" db:"description"`
}

// Validate checks if the word meets basic validation criteria
func (w *Word) Validate() error {
	if w.Hindi == "" || w.English == "" {
		return errors.New("word must have Hindi and English translations")
	}
	
	validDifficulties := map[string]bool{
		"easy":   true,
		"medium": true,
		"hard":   true,
	}
	
	if !validDifficulties[w.Difficulty] {
		w.Difficulty = "easy" // Default to easy
	}
	
	return nil
}

// WordTranslation provides a simplified view of word translations
type WordTranslation struct {
	Hindi     string `json:"hindi"`
	Hinglish  string `json:"hinglish"`
	English   string `json:"english"`
}

// RandomWordFilter allows filtering random word selection
type RandomWordFilter struct {
	Difficulty string
	GroupID    *int64
}
