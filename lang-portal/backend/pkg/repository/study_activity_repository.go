package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pavittarx/lang-portal/backend/pkg/models"
)

// StudyActivityRepository handles database operations for study activities
type StudyActivityRepository struct {
	db *sql.DB
}

// NewStudyActivityRepository creates a new instance of StudyActivityRepository
func NewStudyActivityRepository(db *sql.DB) *StudyActivityRepository {
	return &StudyActivityRepository{db: db}
}

// GetAll retrieves all study activities
func (r *StudyActivityRepository) GetAll(ctx context.Context) ([]models.StudyActivity, error) {
	query := `
		SELECT id, name, description, image, score, created_at 
		FROM study_activities 
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query study activities: %w", err)
	}
	defer rows.Close()

	var activities []models.StudyActivity
	for rows.Next() {
		var activity models.StudyActivity
		err := rows.Scan(
			&activity.ID,
			&activity.Name,
			&activity.Description,
			&activity.Image,
			&activity.Score,
			&activity.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan study activity: %w", err)
		}
		activities = append(activities, activity)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over study activities: %w", err)
	}

	return activities, nil
}
