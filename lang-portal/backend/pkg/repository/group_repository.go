package repository

import (
	"context"
	"database/sql"
	"errors"
	"log"

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

// ListPaginated retrieves paginated groups
func (r *SQLiteGroupRepository) ListPaginated(ctx context.Context, req models.PaginationRequest) ([]models.Group, int64, error) {
	// Base query to count total groups
	countQuery := `SELECT COUNT(*) FROM groups`
	var total int64
	err := r.DB.QueryRowContext(ctx, countQuery).Scan(&total)
	if err != nil {
		log.Printf("Error counting groups: %v", err)
		return nil, 0, err
	}
	
	// Paginated query with sorting
	query := `SELECT id, name, description FROM groups 
			  ORDER BY ` + req.Sort + ` ` + req.Order + ` 
			  LIMIT ? OFFSET ?`
	
	rows, err := r.DB.QueryContext(ctx, query, req.Limit, req.Offset)
	if err != nil {
		log.Printf("Error listing paginated groups: %v", err)
		return nil, 0, err
	}
	defer rows.Close()
	
	var groups []models.Group
	for rows.Next() {
		var group models.Group
		err := rows.Scan(&group.ID, &group.Name, &group.Description)
		if err != nil {
			log.Printf("Error scanning paginated group: %v", err)
			return nil, 0, err
		}
		groups = append(groups, group)
	}
	
	if err = rows.Err(); err != nil {
		log.Printf("Error in paginated group rows: %v", err)
		return nil, 0, err
	}
	
	return groups, total, nil
}
