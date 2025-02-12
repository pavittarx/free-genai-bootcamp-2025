package repository

import (
	"context"

	"github.com/pavittarx/lang-portal/pkg/models"
)

// WordQueryRepository defines the interface for querying words
type WordQueryRepository interface {
	// GetByID retrieves a word by its unique identifier
	GetByID(ctx context.Context, id int64) (*models.Word, error)

	// ListByGroup retrieves all words belonging to a specific group
	ListByGroup(ctx context.Context, groupID int64) ([]models.Word, error)

	// GetRandom retrieves a random word based on optional filters
	GetRandom(ctx context.Context, filter *models.RandomWordFilter) (*models.Word, error)

	// Search finds words matching a given search criteria
	Search(ctx context.Context, query string) ([]models.Word, error)

	// Count returns the total number of words in the repository
	Count(ctx context.Context) (int64, error)
}
