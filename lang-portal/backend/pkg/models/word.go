package models

import (
	"errors"
	"math/rand"
	"strings"
	"time"
	"unicode"
)

// Word represents the structure of a word in the language portal
type Word struct {
	ID        int64     `json:"id" db:"id"`
	Hindi     string    `json:"hindi" db:"hindi"`
	Scrambled string    `json:"scrambled" db:"scrambled"`
	Hinglish  string    `json:"hinglish" db:"hinglish"`
	English   string    `json:"english" db:"english"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// Validate performs validation checks on the Word struct
func (w *Word) Validate() error {
	// Trim whitespace
	w.Hindi = strings.TrimSpace(w.Hindi)
	w.Scrambled = strings.TrimSpace(w.Scrambled)
	w.Hinglish = strings.TrimSpace(w.Hinglish)
	w.English = strings.TrimSpace(w.English)

	// Check for empty fields
	if w.Hindi == "" {
		return errors.New("hindi word cannot be empty")
	}
	if w.English == "" {
		return errors.New("english word cannot be empty")
	}

	// Validate Hindi characters
	for _, r := range w.Hindi {
		if !unicode.Is(unicode.Devanagari, r) && !unicode.IsSpace(r) {
			return errors.New("hindi word must contain only Devanagari characters")
		}
	}

	// Validate English characters
	for _, r := range w.English {
		if !unicode.IsLetter(r) && !unicode.IsSpace(r) {
			return errors.New("english word must contain only alphabetic characters")
		}
	}

	// Ensure scrambled word is not longer than original
	if len(w.Scrambled) > len(w.Hindi) {
		return errors.New("scrambled word cannot be longer than original word")
	}

	// Set created_at if not already set
	if w.CreatedAt.IsZero() {
		w.CreatedAt = time.Now()
	}

	return nil
}

// Sanitize removes any potentially harmful content
func (w *Word) Sanitize() {
	// Remove any leading/trailing whitespace
	w.Hindi = strings.TrimSpace(w.Hindi)
	w.Scrambled = strings.TrimSpace(w.Scrambled)
	w.Hinglish = strings.TrimSpace(w.Hinglish)
	w.English = strings.TrimSpace(w.English)

	// Optionally, you could add more sanitization logic here
	// For example, converting to lowercase, removing special characters, etc.
}

// GenerateScrambledWord creates a scrambled version of the Hindi word if not provided
func (w *Word) GenerateScrambledWord() {
	if w.Scrambled == "" && w.Hindi != "" {
		// Simple scrambling algorithm
		runes := []rune(w.Hindi)
		for i := len(runes) - 1; i > 0; i-- {
			j := rand.Intn(i + 1)
			runes[i], runes[j] = runes[j], runes[i]
		}
		w.Scrambled = string(runes)
	}
}
