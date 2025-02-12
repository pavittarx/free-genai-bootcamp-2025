package models

// WordGroup represents the relationship between a word and a group
type WordGroup struct {
	ID      int64 `json:"id" db:"id"`
	WordID  int64 `json:"word_id" db:"word_id"`
	GroupID int64 `json:"group_id" db:"group_id"`
}

// Validate checks the integrity of the WordGroup relationship
func (wg *WordGroup) Validate() error {
	// Ensure both word and group IDs are positive
	if wg.WordID <= 0 {
		return ErrInvalidWordID
	}

	if wg.GroupID <= 0 {
		return ErrInvalidGroupID
	}

	return nil
}

// ValidationError represents a custom validation error
type ValidationError struct {
	Field   string
	Message string
}

// Error implements the error interface
func (ve ValidationError) Error() string {
	return ve.Message
}

// Errors specific to WordGroup validation
var (
	// ErrInvalidWordID is returned when a word ID is invalid
	ErrInvalidWordID = ValidationError{
		Field:   "word_id",
		Message: "Word ID must be a positive integer",
	}

	// ErrInvalidGroupID is returned when a group ID is invalid
	ErrInvalidGroupID = ValidationError{
		Field:   "group_id",
		Message: "Group ID must be a positive integer",
	}
)
