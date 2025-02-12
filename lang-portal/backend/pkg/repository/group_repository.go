package repository

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"strings"

	"github.com/pavittarx/lang-portal/pkg/models"
)

var (
	// ErrGroupNotFound is returned when a group cannot be found
	ErrGroupNotFound = errors.New("group not found")
)

// GroupRepository defines methods for group-related database operations
type GroupRepository interface {
	// GetByID retrieves a single group by its unique identifier
	GetByID(ctx context.Context, id int64) (*models.Group, error)
	
	// List retrieves all available groups
	List(ctx context.Context) ([]models.Group, error)
	
	// Create adds a new group to the database
	Create(ctx context.Context, group *models.Group) (int64, error)
	
	// ListPaginated retrieves paginated groups
	ListPaginated(ctx context.Context, req models.PaginationRequest) ([]models.Group, int64, error)
}

// SQLiteGroupRepository provides SQLite-specific implementation for group-related operations
type SQLiteGroupRepository struct {
	DB *sql.DB
}

// NewSQLiteGroupRepository creates a new instance of SQLiteGroupRepository
func NewSQLiteGroupRepository(db *sql.DB) *SQLiteGroupRepository {
	return &SQLiteGroupRepository{DB: db}
}

// GetByID retrieves a group by its ID
func (r *SQLiteGroupRepository) GetByID(ctx context.Context, id int64) (*models.Group, error) {
	query := `SELECT id, name, description FROM groups WHERE id = ?`
	
	group := &models.Group{}
	err := r.DB.QueryRowContext(ctx, query, id).Scan(
		&group.ID, &group.Name, &group.Description,
	)
	
	if err == sql.ErrNoRows {
		return nil, ErrGroupNotFound
	}
	
	if err != nil {
		log.Printf("Error retrieving group: %v", err)
		return nil, err
	}
	
	return group, nil
}

// List retrieves all available groups
func (r *SQLiteGroupRepository) List(ctx context.Context) ([]models.Group, error) {
	query := `SELECT id, name, description FROM groups`
	
	rows, err := r.DB.QueryContext(ctx, query)
	if err != nil {
		log.Printf("Error listing groups: %v", err)
		return nil, err
	}
	defer rows.Close()
	
	var groups []models.Group
	for rows.Next() {
		var group models.Group
		err := rows.Scan(&group.ID, &group.Name, &group.Description)
		if err != nil {
			log.Printf("Error scanning group: %v", err)
			return nil, err
		}
		groups = append(groups, group)
	}
	
	if err = rows.Err(); err != nil {
		log.Printf("Error in group rows: %v", err)
		return nil, err
	}
	
	return groups, nil
}

// Create adds a new group to the database
func (r *SQLiteGroupRepository) Create(ctx context.Context, group *models.Group) (int64, error) {
	query := `INSERT INTO groups (name, description) VALUES (?, ?)`
	
	result, err := r.DB.ExecContext(ctx, query, group.Name, group.Description)
	if err != nil {
		log.Printf("Error creating group: %v", err)
		return 0, err
	}
	
	id, err := result.LastInsertId()
	if err != nil {
		log.Printf("Error getting last insert ID: %v", err)
		return 0, err
	}
	
	return id, nil
}

// ListPaginated retrieves a paginated list of groups
func (r *SQLiteGroupRepository) ListPaginated(ctx context.Context, req models.PaginationRequest) ([]models.Group, int64, error) {
	// Prepare base query
	query := `SELECT id, name, description FROM groups`
	countQuery := `SELECT COUNT(*) FROM groups`

	// Prepare optional filters
	var whereClauses []string
	var args []interface{}

	// Add name filter if provided
	if req.Filter != nil {
		if nameFilter, ok := req.Filter.(models.GroupFilter); ok && nameFilter.Name != "" {
			whereClauses = append(whereClauses, "name LIKE ?")
			args = append(args, "%"+nameFilter.Name+"%")
		}
	}

	// Construct where clause if filters exist
	whereClause := ""
	if len(whereClauses) > 0 {
		whereClause = " WHERE " + strings.Join(whereClauses, " AND ")
	}

	// Count total records
	var total int64
	countQuery += whereClause
	err := r.DB.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		log.Printf("Error counting groups: %v", err)
		return nil, 0, err
	}

	// Add pagination to query
	query += whereClause
	query += ` LIMIT ? OFFSET ?`
	args = append(args, req.Limit, (req.Page-1)*req.Limit)

	// Execute query
	rows, err := r.DB.QueryContext(ctx, query, args...)
	if err != nil {
		log.Printf("Error listing groups: %v", err)
		return nil, 0, err
	}
	defer rows.Close()

	// Scan results
	var groups []models.Group
	for rows.Next() {
		var g models.Group
		if err := rows.Scan(&g.ID, &g.Name, &g.Description); err != nil {
			log.Printf("Error scanning group: %v", err)
			return nil, 0, err
		}
		groups = append(groups, g)
	}

	// Check for any errors during iteration
	if err = rows.Err(); err != nil {
		log.Printf("Error in group rows: %v", err)
		return nil, 0, err
	}

	return groups, total, nil
}
