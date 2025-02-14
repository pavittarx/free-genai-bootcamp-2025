package services

import (
	"context"
	"fmt"
	"time"

	"github.com/pavittarx/lang-portal/backend/pkg/models"
	"github.com/pavittarx/lang-portal/backend/pkg/repository"
)

// SessionService handles business logic for sessions
type SessionService struct {
	repo *repository.SessionRepository
}

// NewSessionService creates a new instance of SessionService
func NewSessionService(repo *repository.SessionRepository) *SessionService {
	return &SessionService{repo: repo}
}

// CreateSession starts a new learning session
func (s *SessionService) CreateSession(ctx context.Context, activityID int64) (*models.Session, error) {
	// Create a new session with the current time as start_time
	session := &models.Session{
		ActivityID: activityID,
		StartTime:  time.Now(),
	}

	// Save the session to the database
	err := s.repo.Create(ctx, session)
	if err != nil {
		return nil, err
	}

	return session, nil
}

// GetSessionByID retrieves a specific session
func (s *SessionService) GetSessionByID(ctx context.Context, id int64) (*models.Session, error) {
	return s.repo.GetByID(ctx, id)
}

// GetSessionByIDWithActivities retrieves a specific session with its activities
func (s *SessionService) GetSessionByIDWithActivities(ctx context.Context, id int64) (*models.SessionWithActivities, error) {
	// Retrieve session with activities from repository
	return s.repo.GetByIDWithActivities(ctx, id)
}

// EndSession completes a session by setting the end time and calculating the score
func (s *SessionService) EndSession(ctx context.Context, id int64, score int) error {
	// Retrieve the existing session
	session, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Set end time and score
	now := time.Now()
	session.EndTime = &now
	session.Score = score

	// Validate the updated session
	if err := session.Validate(); err != nil {
		return err
	}

	// Update the session in the repository
	return s.repo.Update(ctx, session)
}

// ListSessions retrieves a list of sessions with pagination
func (s *SessionService) ListSessions(ctx context.Context, page, pageSize int) ([]models.Session, error) {
	// Calculate offset based on page and page size
	offset := (page - 1) * pageSize

	return s.repo.List(ctx, pageSize, offset)
}

// DeleteAllSessions removes all sessions and their associated session activities
func (s *SessionService) DeleteAllSessions(ctx context.Context) (int64, error) {
	// First, delete all session activities
	_, err := s.repo.DeleteAllSessionActivities(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to delete session activities: %w", err)
	}

	// Then, delete all sessions
	sessionsDeleted, err := s.repo.DeleteAll(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to delete sessions: %w", err)
	}

	return sessionsDeleted, nil
}
