package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/pavittarx/lang-portal/backend/pkg/models"
)

// WordRepository defines the interface for word-related database operations
type WordRepository interface {
	// Create a new word
	Create(ctx context.Context, word *models.Word) error

	// GetByID retrieves a word by its ID
	GetByID(ctx context.Context, id int64) (*models.Word, error)

	// Update an existing word
	Update(ctx context.Context, word *models.Word) error

	// Delete a word by its ID
	Delete(ctx context.Context, id int64) error

	// List words with pagination and optional filtering
	List(ctx context.Context, params ListWordsParams) ([]models.Word, int, error)

	// GetRandomWord retrieves a random word from the database
	GetRandomWord(ctx context.Context) (*models.Word, error)
}

// ListWordsParams defines parameters for listing words
type ListWordsParams struct {
	Page     int
	PageSize int
	Search   string
	Language string // "hindi", "english", "hinglish"
}

// SQLiteWordRepository implements WordRepository for SQLite
type SQLiteWordRepository struct {
	db *sql.DB
}

// NewSQLiteWordRepository creates a new instance of SQLiteWordRepository
func NewSQLiteWordRepository(db *sql.DB) *SQLiteWordRepository {
	return &SQLiteWordRepository{db: db}
}

// Create inserts a new word into the database
func (r *SQLiteWordRepository) Create(ctx context.Context, word *models.Word) error {
	// Validate the word before insertion
	if err := word.Validate(); err != nil {
		return fmt.Errorf("invalid word: %w", err)
	}

	// Sanitize the word
	word.Sanitize()

	// Generate scrambled word if not provided
	word.GenerateScrambledWord()

	// Prepare SQL statement
	query := `
		INSERT INTO words (hindi, scrambled, hinglish, english, created_at)
		VALUES (?, ?, ?, ?, ?)
	`

	// Execute the query
	result, err := r.db.ExecContext(ctx, query,
		word.Hindi,
		word.Scrambled,
		word.Hinglish,
		word.English,
		word.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to insert word: %w", err)
	}

	// Get the last inserted ID
	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert ID: %w", err)
	}

	word.ID = id
	return nil
}

// GetByID retrieves a word by its ID
func (r *SQLiteWordRepository) GetByID(ctx context.Context, id int64) (*models.Word, error) {
	query := `
		SELECT id, hindi, scrambled, hinglish, english, created_at
		FROM words
		WHERE id = ?
	`

	word := &models.Word{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&word.ID,
		&word.Hindi,
		&word.Scrambled,
		&word.Hinglish,
		&word.English,
		&word.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("word with ID %d not found", id)
		}
		return nil, fmt.Errorf("failed to retrieve word: %w", err)
	}

	return word, nil
}

// Update modifies an existing word
func (r *SQLiteWordRepository) Update(ctx context.Context, word *models.Word) error {
	// Validate the word before update
	if err := word.Validate(); err != nil {
		return fmt.Errorf("invalid word: %w", err)
	}

	// Sanitize the word
	word.Sanitize()

	// Prepare SQL statement
	query := `
		UPDATE words
		SET hindi = ?, scrambled = ?, hinglish = ?, english = ?
		WHERE id = ?
	`

	// Execute the query
	result, err := r.db.ExecContext(ctx, query,
		word.Hindi,
		word.Scrambled,
		word.Hinglish,
		word.English,
		word.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update word: %w", err)
	}

	// Check if any rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no word found with ID %d", word.ID)
	}

	return nil
}

// Delete removes a word by its ID
func (r *SQLiteWordRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM words WHERE id = ?`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete word: %w", err)
	}

	// Check if any rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no word found with ID %d", id)
	}

	return nil
}

// List retrieves words with pagination and optional filtering
func (r *SQLiteWordRepository) List(ctx context.Context, params ListWordsParams) ([]models.Word, int, error) {
	// Log input parameters for debugging
	log.Printf("SQLiteWordRepository.List called with params: %+v", params)

	// Set default pagination
	if params.Page < 1 {
		params.Page = 1
	}
	if params.PageSize < 1 {
		params.PageSize = 10
	}
	offset := (params.Page - 1) * params.PageSize

	// Base query
	baseQuery := `FROM words WHERE 1=1`
	args := []interface{}{}

	// Add search filter if provided
	if params.Search != "" {
		searchParam := "%" + params.Search + "%"
		switch params.Language {
		case "hindi":
			baseQuery += ` AND hindi LIKE ?`
			args = append(args, searchParam)
		case "english":
			baseQuery += ` AND english LIKE ?`
			args = append(args, searchParam)
		case "hinglish":
			baseQuery += ` AND hinglish LIKE ?`
			args = append(args, searchParam)
		default:
			// Search across all fields if no specific language is specified
			baseQuery += ` AND (hindi LIKE ? OR english LIKE ? OR hinglish LIKE ?)`
			args = append(args, searchParam, searchParam, searchParam)
		}
	}

	// Log the constructed query and arguments
	log.Printf("List query: %s", baseQuery)
	log.Printf("List query arguments: %v", args)

	// Count total matching records
	countQuery := `SELECT COUNT(*) ` + baseQuery
	var totalCount int
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&totalCount); err != nil {
		log.Printf("Error counting words: %v", err)
		return nil, 0, fmt.Errorf("failed to count words: %w", err)
	}

	log.Printf("Total count of words: %d", totalCount)

	// Retrieve words with pagination
	query := `SELECT id, hindi, scrambled, hinglish, english, created_at ` +
		baseQuery + ` LIMIT ? OFFSET ?`
	args = append(args, params.PageSize, offset)

	// Log the final query with arguments
	log.Printf("Final query: %s", query)
	log.Printf("Final query arguments: %v", args)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		log.Printf("Error querying words: %v", err)
		return nil, 0, fmt.Errorf("failed to list words: %w", err)
	}
	defer rows.Close()

	// Scan results
	var words []models.Word
	for rows.Next() {
		var word models.Word
		if err := rows.Scan(
			&word.ID,
			&word.Hindi,
			&word.Scrambled,
			&word.Hinglish,
			&word.English,
			&word.CreatedAt,
		); err != nil {
			log.Printf("Error scanning word: %v", err)
			return nil, 0, fmt.Errorf("error scanning word: %w", err)
		}
		words = append(words, word)
	}

	// Check for any errors encountered during iteration
	if err := rows.Err(); err != nil {
		log.Printf("Error iterating rows: %v", err)
		return nil, 0, fmt.Errorf("error iterating rows: %w", err)
	}

	log.Printf("Retrieved %d words", len(words))

	return words, totalCount, nil
}

// GetRandomWord retrieves a random word from the database
func (r *SQLiteWordRepository) GetRandomWord(ctx context.Context) (*models.Word, error) {
	query := `
		SELECT id, hindi, scrambled, hinglish, english, created_at
		FROM words
		ORDER BY RANDOM()
		LIMIT 1
	`

	word := &models.Word{}
	err := r.db.QueryRowContext(ctx, query).Scan(
		&word.ID,
		&word.Hindi,
		&word.Scrambled,
		&word.Hinglish,
		&word.English,
		&word.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no words found in the database")
		}
		return nil, fmt.Errorf("failed to retrieve random word: %w", err)
	}

	return word, nil
}
