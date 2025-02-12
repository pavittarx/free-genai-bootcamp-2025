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
	// ErrWordNotFound is returned when a word cannot be found
	ErrWordNotFound = errors.New("word not found")
)

// WordRepository defines methods for word-related database operations
type WordRepository interface {
	// GetByID retrieves a single word by its unique identifier
	GetByID(ctx context.Context, id int64) (*models.Word, error)
	
	// GetRandom retrieves a random word with optional filtering
	GetRandom(ctx context.Context, filter *models.RandomWordFilter) (*models.Word, error)
	
	// ListByGroup retrieves words belonging to a specific group
	ListByGroup(ctx context.Context, groupID int64) ([]models.Word, error)
	
	// Create adds a new word to the database
	Create(ctx context.Context, word *models.Word) (int64, error)
	
	// Update modifies an existing word
	Update(ctx context.Context, word *models.Word) error
	
	// Delete removes a word from the database
	Delete(ctx context.Context, id int64) error
}

// SQLiteWordRepository provides SQLite-specific implementation for word-related operations
type SQLiteWordRepository struct {
	db *sql.DB
}

// NewWordRepository creates a new instance of SQLiteWordRepository
func NewWordRepository(db *sql.DB) *SQLiteWordRepository {
	return &SQLiteWordRepository{db: db}
}

// GetByID retrieves a word by its ID
func (r *SQLiteWordRepository) GetByID(ctx context.Context, id int64) (*models.Word, error) {
	query := `SELECT id, hindi, scrambled, hinglish, english, difficulty 
			  FROM words WHERE id = ?`
	
	word := &models.Word{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&word.ID, &word.Hindi, &word.Scrambled, 
		&word.Hinglish, &word.English, &word.Difficulty,
	)
	
	if err == sql.ErrNoRows {
		return nil, ErrWordNotFound
	}
	
	if err != nil {
		log.Printf("Error retrieving word: %v", err)
		return nil, err
	}
	
	return word, nil
}

// GetRandom retrieves a random word with optional filtering
func (r *SQLiteWordRepository) GetRandom(ctx context.Context, filter *models.RandomWordFilter) (*models.Word, error) {
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

// ListByGroup retrieves words in a specific group
func (r *SQLiteWordRepository) ListByGroup(ctx context.Context, groupID int64) ([]models.Word, error) {
	query := `SELECT w.id, w.hindi, w.scrambled, w.hinglish, w.english, w.difficulty 
			  FROM words w
			  JOIN word_groups wg ON w.id = wg.word_id
			  WHERE wg.group_id = ?`
	
	rows, err := r.db.QueryContext(ctx, query, groupID)
	if err != nil {
		log.Printf("Error listing words by group: %v", err)
		return nil, err
	}
	defer rows.Close()
	
	var words []models.Word
	for rows.Next() {
		var word models.Word
		err := rows.Scan(
			&word.ID, &word.Hindi, &word.Scrambled, 
			&word.Hinglish, &word.English, &word.Difficulty,
		)
		if err != nil {
			log.Printf("Error scanning word: %v", err)
			return nil, err
		}
		words = append(words, word)
	}
	
	if err = rows.Err(); err != nil {
		log.Printf("Error in word rows: %v", err)
		return nil, err
	}
	
	return words, nil
}

// Create adds a new word to the database
func (r *SQLiteWordRepository) Create(ctx context.Context, word *models.Word) (int64, error) {
	query := `INSERT INTO words (hindi, scrambled, hinglish, english, difficulty) 
			  VALUES (?, ?, ?, ?, ?)`
	
	result, err := r.db.ExecContext(ctx, query, 
		word.Hindi, word.Scrambled, word.Hinglish, word.English, word.Difficulty,
	)
	
	if err != nil {
		log.Printf("Error creating word: %v", err)
		return 0, err
	}
	
	id, err := result.LastInsertId()
	if err != nil {
		log.Printf("Error getting last insert ID: %v", err)
		return 0, err
	}
	
	return id, nil
}

// Update modifies an existing word
func (r *SQLiteWordRepository) Update(ctx context.Context, word *models.Word) error {
	query := `UPDATE words SET hindi = ?, scrambled = ?, hinglish = ?, english = ?, difficulty = ?
			  WHERE id = ?`
	
	_, err := r.db.ExecContext(ctx, query, 
		word.Hindi, word.Scrambled, word.Hinglish, word.English, word.Difficulty, word.ID,
	)
	
	if err != nil {
		log.Printf("Error updating word: %v", err)
		return err
	}
	
	return nil
}

// Delete removes a word from the database
func (r *SQLiteWordRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM words WHERE id = ?`
	
	_, err := r.db.ExecContext(ctx, query, id)
	
	if err != nil {
		log.Printf("Error deleting word: %v", err)
		return err
	}
	
	return nil
}
