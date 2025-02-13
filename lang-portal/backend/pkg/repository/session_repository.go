package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pavittarx/lang-portal/backend/pkg/models"
)

// SessionRepository handles database operations for sessions
type SessionRepository struct {
	db *sql.DB
}

// NewSessionRepository creates a new instance of SessionRepository
func NewSessionRepository(db *sql.DB) *SessionRepository {
	return &SessionRepository{db: db}
}

// Create starts a new session
func (r *SessionRepository) Create(ctx context.Context, session *models.Session) error {
	query := `
		INSERT INTO sessions (activity_id, group_id, start_time, score, created_at) 
		VALUES (?, ?, ?, ?, ?)
	`

	result, err := r.db.ExecContext(ctx, query,
		session.ActivityID,
		session.GroupID,
		session.StartTime,
		session.Score,
		session.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert ID: %w", err)
	}

	session.ID = id
	return nil
}

// GetByID retrieves a session by its ID
func (r *SessionRepository) GetByID(ctx context.Context, id int64) (*models.Session, error) {
	query := `
		SELECT id, activity_id, group_id, start_time, end_time, score, created_at 
		FROM sessions 
		WHERE id = ?
	`

	var session models.Session
	var endTime sql.NullTime
	var groupID sql.NullInt64

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&session.ID,
		&session.ActivityID,
		&groupID,
		&session.StartTime,
		&endTime,
		&session.Score,
		&session.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("session not found: %w", err)
		}
		return nil, fmt.Errorf("failed to retrieve session: %w", err)
	}

	// Handle nullable fields
	if groupID.Valid {
		session.GroupID = &groupID.Int64
	}
	if endTime.Valid {
		session.EndTime = &endTime.Time
	}

	return &session, nil
}

// Update updates an existing session
func (r *SessionRepository) Update(ctx context.Context, session *models.Session) error {
	query := `
		UPDATE sessions 
		SET activity_id = ?, group_id = ?, end_time = ?, score = ? 
		WHERE id = ?
	`

	_, err := r.db.ExecContext(ctx, query,
		session.ActivityID,
		session.GroupID,
		session.EndTime,
		session.Score,
		session.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update session: %w", err)
	}

	return nil
}

// Delete removes a session
func (r *SessionRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM sessions WHERE id = ?`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no session found with ID %d", id)
	}

	return nil
}

// List retrieves sessions with optional pagination
func (r *SessionRepository) List(ctx context.Context, limit, offset int) ([]models.Session, error) {
	query := `
		SELECT id, activity_id, group_id, start_time, end_time, score, created_at 
		FROM sessions 
		ORDER BY created_at DESC 
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query sessions: %w", err)
	}
	defer rows.Close()

	var sessions []models.Session
	for rows.Next() {
		var session models.Session
		var endTime sql.NullTime
		var groupID sql.NullInt64

		err := rows.Scan(
			&session.ID,
			&session.ActivityID,
			&groupID,
			&session.StartTime,
			&endTime,
			&session.Score,
			&session.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan session: %w", err)
		}

		// Handle nullable fields
		if groupID.Valid {
			session.GroupID = &groupID.Int64
		}
		if endTime.Valid {
			session.EndTime = &endTime.Time
		}

		sessions = append(sessions, session)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over sessions: %w", err)
	}

	return sessions, nil
}
