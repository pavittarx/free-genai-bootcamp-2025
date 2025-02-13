package repository

import (
	"context"
	"fmt"

	"github.com/pavittarx/lang-portal/backend/pkg/models"
)

// GetWordsByGroupID retrieves all words associated with a specific group
func (r *SQLiteWordRepository) GetWordsByGroupID(ctx context.Context, groupID int64) ([]models.Word, error) {
	query := `
		SELECT w.id, w.hindi, w.english, w.hinglish, w.created_at
		FROM words w
		INNER JOIN word_groups wg ON w.id = wg.word_id
		WHERE wg.group_id = ?
		ORDER BY w.created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, groupID)
	if err != nil {
		return nil, fmt.Errorf("failed to query words by group ID: %w", err)
	}
	defer rows.Close()

	var words []models.Word
	for rows.Next() {
		var word models.Word
		err := rows.Scan(
			&word.ID,
			&word.Hindi,
			&word.English,
			&word.Hinglish,
			&word.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan word: %w", err)
		}
		words = append(words, word)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return words, nil
}
