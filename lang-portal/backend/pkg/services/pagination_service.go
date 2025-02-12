package services

import (
	"context"

	"github.com/pavittarx/lang-portal/pkg/models"
	"github.com/pavittarx/lang-portal/pkg/repository"
)

// WordPaginationService provides paginated word retrieval methods
type WordPaginationService struct {
	wordRepo repository.PaginatedWordRepository
}

// NewWordPaginationService creates a new instance of WordPaginationService
func NewWordPaginationService(wordRepo repository.PaginatedWordRepository) *WordPaginationService {
	return &WordPaginationService{
		wordRepo: wordRepo,
	}
}

// ListWords retrieves a paginated list of words
func (s *WordPaginationService) ListWords(ctx context.Context, req models.PaginationRequest) (*models.PaginationResponse, error) {
	req.Validate()
	
	words, total, err := s.wordRepo.ListWords(ctx, req)
	if err != nil {
		return nil, err
	}
	
	return &models.PaginationResponse{
		Data:       words,
		Total:      total,
		Page:       req.Page,
		Limit:      req.Limit,
		TotalPages: int((total + int64(req.Limit) - 1) / int64(req.Limit)),
	}, nil
}

// ListWordsByGroup retrieves paginated words for a specific group
func (s *WordPaginationService) ListWordsByGroup(ctx context.Context, groupID int64, req models.PaginationRequest) (*models.PaginationResponse, error) {
	req.Validate()
	
	words, total, err := s.wordRepo.ListWordsByGroup(ctx, groupID, req)
	if err != nil {
		return nil, err
	}
	
	return &models.PaginationResponse{
		Data:       words,
		Total:      total,
		Page:       req.Page,
		Limit:      req.Limit,
		TotalPages: int((total + int64(req.Limit) - 1) / int64(req.Limit)),
	}, nil
}

// GroupPaginationService provides paginated group retrieval methods
type GroupPaginationService struct {
	groupRepo repository.PaginatedGroupRepository
}

// NewGroupPaginationService creates a new instance of GroupPaginationService
func NewGroupPaginationService(groupRepo repository.PaginatedGroupRepository) *GroupPaginationService {
	return &GroupPaginationService{
		groupRepo: groupRepo,
	}
}

// ListGroups retrieves a paginated list of groups
func (s *GroupPaginationService) ListGroups(ctx context.Context, req models.PaginationRequest) (*models.PaginationResponse, error) {
	req.Validate()
	
	groups, total, err := s.groupRepo.ListGroups(ctx, req)
	if err != nil {
		return nil, err
	}
	
	return &models.PaginationResponse{
		Data:       groups,
		Total:      total,
		Page:       req.Page,
		Limit:      req.Limit,
		TotalPages: int((total + int64(req.Limit) - 1) / int64(req.Limit)),
	}, nil
}

// StudyActivityPaginationService provides paginated study activity retrieval methods
type StudyActivityPaginationService struct {
	activityRepo repository.PaginatedStudyActivityRepository
}

// NewStudyActivityPaginationService creates a new instance of StudyActivityPaginationService
func NewStudyActivityPaginationService(activityRepo repository.PaginatedStudyActivityRepository) *StudyActivityPaginationService {
	return &StudyActivityPaginationService{
		activityRepo: activityRepo,
	}
}

// ListStudyActivities retrieves a paginated list of study activities
func (s *StudyActivityPaginationService) ListStudyActivities(ctx context.Context, req models.PaginationRequest) (*models.PaginationResponse, error) {
	req.Validate()
	
	activities, total, err := s.activityRepo.ListStudyActivities(ctx, req)
	if err != nil {
		return nil, err
	}
	
	return &models.PaginationResponse{
		Data:       activities,
		Total:      total,
		Page:       req.Page,
		Limit:      req.Limit,
		TotalPages: int((total + int64(req.Limit) - 1) / int64(req.Limit)),
	}, nil
}

// ListSessionActivities retrieves paginated activities for a specific session
func (s *StudyActivityPaginationService) ListSessionActivities(ctx context.Context, sessionID int64, req models.PaginationRequest) (*models.PaginationResponse, error) {
	req.Validate()
	
	activities, total, err := s.activityRepo.ListSessionActivities(ctx, sessionID, req)
	if err != nil {
		return nil, err
	}
	
	return &models.PaginationResponse{
		Data:       activities,
		Total:      total,
		Page:       req.Page,
		Limit:      req.Limit,
		TotalPages: int((total + int64(req.Limit) - 1) / int64(req.Limit)),
	}, nil
}

// SessionPaginationService provides paginated session retrieval methods
type SessionPaginationService struct {
	sessionRepo repository.PaginatedSessionRepository
}

// NewSessionPaginationService creates a new instance of SessionPaginationService
func NewSessionPaginationService(sessionRepo repository.PaginatedSessionRepository) *SessionPaginationService {
	return &SessionPaginationService{
		sessionRepo: sessionRepo,
	}
}

// ListSessions retrieves a paginated list of sessions
func (s *SessionPaginationService) ListSessions(ctx context.Context, req models.PaginationRequest) (*models.PaginationResponse, error) {
	req.Validate()
	
	sessions, total, err := s.sessionRepo.ListSessions(ctx, req)
	if err != nil {
		return nil, err
	}
	
	return &models.PaginationResponse{
		Data:       sessions,
		Total:      total,
		Page:       req.Page,
		Limit:      req.Limit,
		TotalPages: int((total + int64(req.Limit) - 1) / int64(req.Limit)),
	}, nil
}
