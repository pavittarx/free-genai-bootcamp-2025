package services

import (
	"context"
	"time"

	"github.com/pavittarx/lang-portal/backend/pkg/models"
	"github.com/pavittarx/lang-portal/backend/pkg/repository"
)

// SessionActivityService handles business logic for session activities
type SessionActivityService struct {
	repo *repository.SessionActivityRepository
	sessionRepo *repository.SessionRepository
}

// NewSessionActivityService creates a new instance of SessionActivityService
func NewSessionActivityService(
	repo *repository.SessionActivityRepository, 
	sessionRepo *repository.SessionRepository,
) *SessionActivityService {
	return &SessionActivityService{
		repo: repo,
		sessionRepo: sessionRepo,
	}
}

// AddSessionActivity adds a new activity to an existing session
func (s *SessionActivityService) AddSessionActivity(
	ctx context.Context, 
	sessionID, activityID int64, 
	challenge, answer, input string,
	score int,
) (*models.SessionActivity, error) {
	// Validate session exists
	_, err := s.sessionRepo.GetByID(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	// Create session activity
	sessionActivity := &models.SessionActivity{
		SessionID:   sessionID,
		ActivityID:  activityID,
		Challenge:   challenge,
		Answer:      answer,
		Input:       input,
		Result:      "", // To be determined by scoring logic
		Score:       score,
		CreatedAt:   time.Now(),
	}

	// Validate the session activity
	if err := sessionActivity.Validate(); err != nil {
		return nil, err
	}

	// Save to repository
	if err := s.repo.Create(ctx, sessionActivity); err != nil {
		return nil, err
	}

	return sessionActivity, nil
}

// GetSessionActivities retrieves all activities for a specific session
func (s *SessionActivityService) GetSessionActivities(
	ctx context.Context, 
	sessionID int64,
) ([]models.SessionActivity, error) {
	return s.repo.ListBySessionID(ctx, sessionID)
}

// DeleteSessionActivity removes a session activity
func (s *SessionActivityService) DeleteSessionActivity(
	ctx context.Context, 
	sessionActivityID int64,
) error {
	return s.repo.Delete(ctx, sessionActivityID)
}
