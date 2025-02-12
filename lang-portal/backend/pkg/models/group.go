package models

// Group represents a category of words
type Group struct {
	ID          int64  `json:"id" db:"id"`
	Name        string `json:"name" db:"name"`
	Description string `json:"description" db:"description"`
}

// GroupFilter allows filtering and searching groups
type GroupFilter struct {
	Name        *string
	Description *string
}

// Validate checks if the group meets basic validation criteria
func (g *Group) Validate() error {
	if g.Name == "" {
		return ErrInvalidGroup
	}
	return nil
}

// WordGroup represents the many-to-many relationship between words and groups
type WordGroup struct {
	WordID  int64 `json:"word_id" db:"word_id"`
	GroupID int64 `json:"group_id" db:"group_id"`
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
