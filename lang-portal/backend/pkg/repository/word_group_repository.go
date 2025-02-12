package repository

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/pavittarx/lang-portal/pkg/models"
)

var (
	// ErrWordGroupNotFound is returned when a word group relationship cannot be found
	ErrWordGroupNotFound = errors.New("word group relationship not found")
)

// WordGroupRepository defines methods for word group relationship operations
type WordGroupRepository interface {
	// Create adds a new word-group relationship
	Create(ctx context.Context, wordGroup *models.WordGroup) (int64, error)
	
	// GetByWordID retrieves all groups for a specific word
	GetByWordID(ctx context.Context, wordID int64) ([]models.Group, error)
	
	// GetByGroupID retrieves all words in a specific group
	GetByGroupID(ctx context.Context, groupID int64) ([]models.Word, error)
	
	// Remove deletes a word-group relationship
	Remove(ctx context.Context, wordID, groupID int64) error
}

// SQLiteWordGroupRepository provides SQLite-specific implementation for word group operations
type SQLiteWordGroupRepository struct {
	db *sql.DB
}

// NewWordGroupRepository creates a new instance of SQLiteWordGroupRepository
func NewWordGroupRepository(db *sql.DB) *SQLiteWordGroupRepository {
	return &SQLiteWordGroupRepository{db: db}
}

// Create adds a new word-group relationship
func (r *SQLiteWordGroupRepository) Create(ctx context.Context, wordGroup *models.WordGroup) (int64, error) {
	// Validate the word group relationship
	if err := wordGroup.Validate(); err != nil {
		return 0, err
	}

	query := `INSERT INTO word_groups (word_id, group_id) VALUES (?, ?)`
	
	result, err := r.db.ExecContext(ctx, query, wordGroup.WordID, wordGroup.GroupID)
	if err != nil {
		log.Printf("Error creating word group relationship: %v", err)
		return 0, err
	}
	
	id, err := result.LastInsertId()
	if err != nil {
		log.Printf("Error getting last insert ID: %v", err)
		return 0, err
	}
	
	return id, nil
}

// GetByWordID retrieves all groups for a specific word
func (r *SQLiteWordGroupRepository) GetByWordID(ctx context.Context, wordID int64) ([]models.Group, error) {
	query := `SELECT g.id, g.name, g.description 
			  FROM groups g
			  JOIN word_groups wg ON g.id = wg.group_id
			  WHERE wg.word_id = ?`
	
	rows, err := r.db.QueryContext(ctx, query, wordID)
	if err != nil {
		log.Printf("Error retrieving groups for word: %v", err)
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
	
	if len(groups) == 0 {
		return nil, ErrWordGroupNotFound
	}
	
	return groups, nil
}

// GetByGroupID retrieves all words in a specific group
func (r *SQLiteWordGroupRepository) GetByGroupID(ctx context.Context, groupID int64) ([]models.Word, error) {
	query := `SELECT w.id, w.hindi, w.scrambled, w.hinglish, w.english, w.difficulty 
			  FROM words w
			  JOIN word_groups wg ON w.id = wg.word_id
			  WHERE wg.group_id = ?`
	
	rows, err := r.db.QueryContext(ctx, query, groupID)
	if err != nil {
		log.Printf("Error retrieving words for group: %v", err)
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
	
	if len(words) == 0 {
		return nil, ErrWordGroupNotFound
	}
	
	return words, nil
}

// Remove deletes a word-group relationship
func (r *SQLiteWordGroupRepository) Remove(ctx context.Context, wordID, groupID int64) error {
	query := `DELETE FROM word_groups WHERE word_id = ? AND group_id = ?`
	
	result, err := r.db.ExecContext(ctx, query, wordID, groupID)
	if err != nil {
		log.Printf("Error removing word group relationship: %v", err)
		return err
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error checking rows affected: %v", err)
		return err
	}
	
	if rowsAffected == 0 {
		return ErrWordGroupNotFound
	}
	
	return nil
}
