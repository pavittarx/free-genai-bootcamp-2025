package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pavittarx/lang-portal/backend/pkg/models"
)

// SessionActivityRepository handles database operations for session activities
type SessionActivityRepository struct {
	db *sql.DB
}

// NewSessionActivityRepository creates a new instance of SessionActivityRepository
func NewSessionActivityRepository(db *sql.DB) *SessionActivityRepository {
	return &SessionActivityRepository{db: db}
}

// Create adds a new session activity to the database
func (r *SessionActivityRepository) Create(ctx context.Context, sessionActivity *models.SessionActivity) error {
	query := `
		INSERT INTO session_activities 
		(session_id, activity_id, challenge, answer, input, result, score, created_at) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := r.db.ExecContext(ctx, query,
		sessionActivity.SessionID,
		sessionActivity.ActivityID,
		sessionActivity.Challenge,
		sessionActivity.Answer,
		sessionActivity.Input,
		sessionActivity.Result,
		sessionActivity.Score,
		sessionActivity.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create session activity: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert ID: %w", err)
	}

	sessionActivity.ID = id
	return nil
}

// GetByID retrieves a specific session activity
func (r *SessionActivityRepository) GetByID(ctx context.Context, id int64) (*models.SessionActivity, error) {
	query := `
		SELECT id, session_id, activity_id, challenge, answer, input, result, score, created_at 
		FROM session_activities 
		WHERE id = ?
	`

	var sessionActivity models.SessionActivity
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&sessionActivity.ID,
		&sessionActivity.SessionID,
		&sessionActivity.ActivityID,
		&sessionActivity.Challenge,
		&sessionActivity.Answer,
		&sessionActivity.Input,
		&sessionActivity.Result,
		&sessionActivity.Score,
		&sessionActivity.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("session activity not found: %w", err)
		}
		return nil, fmt.Errorf("failed to retrieve session activity: %w", err)
	}

	return &sessionActivity, nil
}

// ListBySessionID retrieves all session activities for a given session
func (r *SessionActivityRepository) ListBySessionID(ctx context.Context, sessionID int64) ([]models.SessionActivity, error) {
	query := `
		SELECT id, session_id, activity_id, challenge, answer, input, result, score, created_at 
		FROM session_activities 
		WHERE session_id = ? 
		ORDER BY created_at
	`

	rows, err := r.db.QueryContext(ctx, query, sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to query session activities: %w", err)
	}
	defer rows.Close()

	var sessionActivities []models.SessionActivity
	for rows.Next() {
		var sa models.SessionActivity
		err := rows.Scan(
			&sa.ID,
			&sa.SessionID,
			&sa.ActivityID,
			&sa.Challenge,
			&sa.Answer,
			&sa.Input,
			&sa.Result,
			&sa.Score,
			&sa.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan session activity: %w", err)
		}
		sessionActivities = append(sessionActivities, sa)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating session activities: %w", err)
	}

	return sessionActivities, nil
}

// Update modifies an existing session activity
func (r *SessionActivityRepository) Update(ctx context.Context, sessionActivity *models.SessionActivity) error {
	query := `
		UPDATE session_activities 
		SET result = ?, score = ? 
		WHERE id = ?
	`

	_, err := r.db.ExecContext(ctx, query,
		sessionActivity.Result,
		sessionActivity.Score,
		sessionActivity.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update session activity: %w", err)
	}

	return nil
}

// Delete removes a session activity
func (r *SessionActivityRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM session_activities WHERE id = ?`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete session activity: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no session activity found with ID %d", id)
	}

	return nil
}
