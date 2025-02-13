package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pavittarx/lang-portal/backend/pkg/models"
)

// SQLiteGroupRepository handles database operations for groups using SQLite
type SQLiteGroupRepository struct {
	db *sql.DB
}

// NewSQLiteGroupRepository creates a new instance of SQLiteGroupRepository
func NewSQLiteGroupRepository(db *sql.DB) *SQLiteGroupRepository {
	return &SQLiteGroupRepository{db: db}
}

// Create adds a new group to the database
func (r *SQLiteGroupRepository) Create(ctx context.Context, group *models.Group) error {
	// Validate the group before creating
	if err := group.Validate(); err != nil {
		return err
	}

	// Sanitize the group name and description
	group.Sanitize()

	// Prepare the SQL statement
	query := `INSERT INTO groups (name, description, created_at) VALUES (?, ?, datetime('now'))`
	result, err := r.db.ExecContext(ctx, query, group.Name, group.Description)
	if err != nil {
		return fmt.Errorf("failed to create group: %w", err)
	}

	// Get the last inserted ID
	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert ID: %w", err)
	}

	group.ID = id
	return nil
}

// GetByID retrieves a group by its ID
func (r *SQLiteGroupRepository) GetByID(ctx context.Context, id int64) (*models.Group, error) {
	query := `SELECT id, name, description, created_at FROM groups WHERE id = ?`
	group := &models.Group{}

	err := r.db.QueryRowContext(ctx, query, id).Scan(&group.ID, &group.Name, &group.Description, &group.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("group with ID %d not found", id)
		}
		return nil, fmt.Errorf("failed to retrieve group: %w", err)
	}

	return group, nil
}

// Update modifies an existing group
func (r *SQLiteGroupRepository) Update(ctx context.Context, group *models.Group) error {
	// Validate the group before updating
	if err := group.Validate(); err != nil {
		return err
	}

	// Sanitize the group name and description
	group.Sanitize()

	// Prepare the SQL statement
	query := `UPDATE groups SET name = ?, description = ? WHERE id = ?`
	result, err := r.db.ExecContext(ctx, query, group.Name, group.Description, group.ID)
	if err != nil {
		return fmt.Errorf("failed to update group: %w", err)
	}

	// Check if any rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("group with ID %d not found", group.ID)
	}

	return nil
}

// Delete removes a group by its ID
func (r *SQLiteGroupRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM groups WHERE id = ?`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete group: %w", err)
	}

	// Check if any rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("group with ID %d not found", id)
	}

	return nil
}

// List retrieves groups with pagination
func (r *SQLiteGroupRepository) List(ctx context.Context, page, pageSize int, search string) ([]models.Group, int, error) {
	// Calculate offset
	offset := (page - 1) * pageSize

	// Base query for counting total groups
	countQuery := `SELECT COUNT(*) FROM groups WHERE 1=1`
	listQuery := `SELECT id, name, description, created_at FROM groups WHERE 1=1`

	// Add search condition if provided
	args := []interface{}{}
	if search != "" {
		countQuery += ` AND (name LIKE ? OR description LIKE ?)`
		listQuery += ` AND (name LIKE ? OR description LIKE ?)`
		args = append(args, "%"+search+"%", "%"+search+"%")
	}

	// Count total groups
	var totalCount int
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&totalCount)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count groups: %w", err)
	}

	// Add pagination to list query
	listQuery += ` LIMIT ? OFFSET ?`
	args = append(args, pageSize, offset)

	// Execute list query
	rows, err := r.db.QueryContext(ctx, listQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list groups: %w", err)
	}
	defer rows.Close()

	// Scan results
	groups := []models.Group{}
	for rows.Next() {
		var group models.Group
		if err := rows.Scan(&group.ID, &group.Name, &group.Description, &group.CreatedAt); err != nil {
			return nil, 0, fmt.Errorf("failed to scan group: %w", err)
		}
		groups = append(groups, group)
	}

	// Check for any errors during iteration
	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error during group listing: %w", err)
	}

	return groups, totalCount, nil
}
