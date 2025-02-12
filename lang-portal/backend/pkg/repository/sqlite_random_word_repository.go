package repository

import (
	"context"
	"database/sql"
	"log"
	"strings"

	"github.com/pavittarx/lang-portal/pkg/models"
)

// SQLiteRandomWordRepository provides SQLite-specific implementation for random word retrieval
type SQLiteRandomWordRepository struct {
	db *sql.DB
}

// NewSQLiteRandomWordRepository creates a new instance of SQLiteRandomWordRepository
func NewSQLiteRandomWordRepository(db *sql.DB) *SQLiteRandomWordRepository {
	return &SQLiteRandomWordRepository{db: db}
}

// GetRandom retrieves a random word with optional filtering
func (r *SQLiteRandomWordRepository) GetRandom(ctx context.Context, filter *models.RandomWordFilter) (*models.Word, error) {
	query := `SELECT id, hindi, scrambled, hinglish, english, difficulty 
			  FROM words`
	
	var args []interface{}
	var conditions []string
	
	// Apply difficulty filter
	if filter != nil && filter.Difficulty != "" {
		conditions = append(conditions, "difficulty = ?")
		args = append(args, filter.Difficulty)
	}
	
	// Apply group filter
	if filter != nil && filter.GroupID != nil {
		conditions = append(conditions, 
			`id IN (SELECT word_id FROM word_groups WHERE group_id = ?)`)
		args = append(args, *filter.GroupID)
	}
	
	// Combine conditions
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}
	
	query += ` ORDER BY RANDOM() LIMIT 1`
	
	word := &models.Word{}
	err := r.db.QueryRowContext(ctx, query, args...).Scan(
		&word.ID, &word.Hindi, &word.Scrambled, 
		&word.Hinglish, &word.English, &word.Difficulty,
	)
	
	if err == sql.ErrNoRows {
		return nil, ErrWordNotFound
	}
	
	if err != nil {
		log.Printf("Error getting random word: %v", err)
		return nil, err
	}
	
	return word, nil
}
