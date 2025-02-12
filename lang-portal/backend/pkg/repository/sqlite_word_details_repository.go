package repository

import (
	"context"
	"database/sql"
	"log"

	"github.com/pavittarx/lang-portal/pkg/models"
)

// SQLiteWordDetailsRepository provides SQLite-specific implementation for word details retrieval
type SQLiteWordDetailsRepository struct {
	db *sql.DB
}

// NewSQLiteWordDetailsRepository creates a new instance of SQLiteWordDetailsRepository
func NewSQLiteWordDetailsRepository(db *sql.DB) *SQLiteWordDetailsRepository {
	return &SQLiteWordDetailsRepository{db: db}
}

// GetByID retrieves a word by its ID with full details
func (r *SQLiteWordDetailsRepository) GetByID(ctx context.Context, id int64) (*models.Word, error) {
	query := `SELECT id, hindi, scrambled, hinglish, english, difficulty, description 
			  FROM words WHERE id = ?`
	
	word := &models.Word{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&word.ID, &word.Hindi, &word.Scrambled, 
		&word.Hinglish, &word.English, &word.Difficulty, &word.Description,
	)
	
	if err == sql.ErrNoRows {
		return nil, ErrWordNotFound
	}
	
	if err != nil {
		log.Printf("Error retrieving word details: %v", err)
		return nil, err
	}
	
	return word, nil
}
