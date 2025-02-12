package repository

import (
	"context"

	"github.com/pavittarx/lang-portal/pkg/models"
)

// PaginatedWordRepository defines methods for paginated word retrieval
type PaginatedWordRepository interface {
	// ListWords retrieves a paginated list of words
	ListWords(ctx context.Context, req models.PaginationRequest) ([]models.Word, int64, error)
	
	// ListWordsByGroup retrieves paginated words for a specific group
	ListWordsByGroup(ctx context.Context, groupID int64, req models.PaginationRequest) ([]models.Word, int64, error)
}

// PaginatedGroupRepository defines methods for paginated group retrieval
type PaginatedGroupRepository interface {
	// ListGroups retrieves a paginated list of groups
	ListGroups(ctx context.Context, req models.PaginationRequest) ([]models.Group, int64, error)
}

// PaginatedStudyActivityRepository defines methods for paginated study activity retrieval
type PaginatedStudyActivityRepository interface {
	// ListStudyActivities retrieves a paginated list of study activities
	ListStudyActivities(ctx context.Context, req models.PaginationRequest) ([]models.StudyActivity, int64, error)
	
	// ListSessionActivities retrieves paginated activities for a specific session
	ListSessionActivities(ctx context.Context, sessionID int64, req models.PaginationRequest) ([]models.StudyActivity, int64, error)
}

// PaginatedSessionRepository defines methods for paginated session retrieval
type PaginatedSessionRepository interface {
	// ListSessions retrieves a paginated list of sessions
	ListSessions(ctx context.Context, req models.PaginationRequest) ([]models.Session, int64, error)
	
	// GetSessionDetails retrieves details of a specific session
	GetSessionDetails(ctx context.Context, sessionID int64) (*models.Session, error)
}
