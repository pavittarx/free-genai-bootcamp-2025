package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"

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
	log.Printf("Retrieving session with ID: %d", id)

	query := `
		SELECT id, activity_id, group_id, start_time, end_time, score, created_at 
		FROM sessions 
		WHERE id = ?
	`

	var session models.Session
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&session.ID,
		&session.ActivityID,
		&session.GroupID,
		&session.StartTime,
		&session.EndTime,
		&session.Score,
		&session.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("no session available with provided session id: %d", id)
	}

	if err != nil {
		log.Printf("Error retrieving session: %v", err)
		return nil, fmt.Errorf("failed to retrieve session: %w", err)
	}

	return &session, nil
}

// GetByIDWithActivities retrieves a session with its associated activities
func (r *SessionRepository) GetByIDWithActivities(ctx context.Context, id int64) (*models.SessionWithActivities, error) {
	log.Printf("Retrieving session with activities for ID: %d", id)

	// First, retrieve the session
	session, err := r.GetByID(ctx, id)
	if err != nil {
		log.Printf("Failed to retrieve session: %v", err)
		return nil, err
	}

	// Then, retrieve session activities
	activitiesQuery := `
		SELECT id, session_id, activity_id, challenge, answer, input, result, score, created_at 
		FROM session_activities 
		WHERE session_id = ?
		ORDER BY created_at ASC
	`

	log.Printf("Executing activities query: %s with session ID: %d", activitiesQuery, id)

	rows, err := r.db.QueryContext(ctx, activitiesQuery, id)
	if err != nil {
		log.Printf("Failed to query session activities: %v", err)
		return nil, fmt.Errorf("failed to query session activities: %w", err)
	}
	defer rows.Close()

	var activities []models.SessionActivity
	for rows.Next() {
		var activity models.SessionActivity
		err := rows.Scan(
			&activity.ID,
			&activity.SessionID,
			&activity.ActivityID,
			&activity.Challenge,
			&activity.Answer,
			&activity.Input,
			&activity.Result,
			&activity.Score,
			&activity.CreatedAt,
		)
		if err != nil {
			log.Printf("Failed to scan session activity: %v", err)
			return nil, fmt.Errorf("failed to scan session activity: %w", err)
		}

		activities = append(activities, activity)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error iterating over session activities: %v", err)
		return nil, fmt.Errorf("error iterating over session activities: %w", err)
	}

	log.Printf("Retrieved %d session activities", len(activities))

	// Combine session and activities
	return &models.SessionWithActivities{
		Session:     *session,
		Activities: activities,
	}, nil
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

// DeleteAllSessionActivities removes all session activities from the database
func (r *SessionRepository) DeleteAllSessionActivities(ctx context.Context) (int64, error) {
	query := `DELETE FROM session_activities`

	result, err := r.db.ExecContext(ctx, query)
	if err != nil {
		return 0, fmt.Errorf("failed to delete all session activities: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("error checking rows affected: %w", err)
	}

	return rowsAffected, nil
}

// DeleteAll removes all sessions from the database
func (r *SessionRepository) DeleteAll(ctx context.Context) (int64, error) {
	query := `DELETE FROM sessions`

	result, err := r.db.ExecContext(ctx, query)
	if err != nil {
		return 0, fmt.Errorf("failed to delete all sessions: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("error checking rows affected: %w", err)
	}

	return rowsAffected, nil
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
