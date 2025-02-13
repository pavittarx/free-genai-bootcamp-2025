package services

import (
	"context"
	"fmt"
	"log"

	"github.com/pavittarx/lang-portal/backend/pkg/models"
	"github.com/pavittarx/lang-portal/backend/pkg/repository"
)

// WordService provides business logic for word-related operations
type WordService struct {
	repo repository.WordRepository
}

// NewWordService creates a new instance of WordService
func NewWordService(repo repository.WordRepository) *WordService {
	return &WordService{repo: repo}
}

// CreateWord handles the creation of a new word
func (s *WordService) CreateWord(ctx context.Context, word *models.Word) error {
	// Additional business logic can be added here
	// For example, checking for duplicates, logging, etc.

	// Validate the word
	if err := word.Validate(); err != nil {
		return fmt.Errorf("word validation failed: %w", err)
	}

	// Generate scrambled word if not provided
	word.GenerateScrambledWord()

	// Persist the word
	return s.repo.Create(ctx, word)
}

// GetWordByID retrieves a word by its ID
func (s *WordService) GetWordByID(ctx context.Context, id int64) (*models.Word, error) {
	// Additional business logic can be added here
	// For example, caching, logging, etc.

	return s.repo.GetByID(ctx, id)
}

// UpdateWord updates an existing word
func (s *WordService) UpdateWord(ctx context.Context, word *models.Word) error {
	// Additional business logic can be added here
	// For example, checking permissions, logging changes, etc.

	// Validate the word
	if err := word.Validate(); err != nil {
		return fmt.Errorf("word validation failed: %w", err)
	}

	// Ensure the word exists before updating
	existingWord, err := s.repo.GetByID(ctx, word.ID)
	if err != nil {
		return fmt.Errorf("word not found: %w", err)
	}

	// Merge existing and new word data
	if word.Hindi == "" {
		word.Hindi = existingWord.Hindi
	}
	if word.English == "" {
		word.English = existingWord.English
	}
	if word.Hinglish == "" {
		word.Hinglish = existingWord.Hinglish
	}
	if word.Scrambled == "" {
		word.GenerateScrambledWord()
	}

	// Update the word
	return s.repo.Update(ctx, word)
}

// DeleteWord removes a word by its ID
func (s *WordService) DeleteWord(ctx context.Context, id int64) error {
	// Additional business logic can be added here
	// For example, checking permissions, logging, etc.

	// Verify the word exists before deleting
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("word not found: %w", err)
	}

	return s.repo.Delete(ctx, id)
}

// ListWords retrieves a list of words with pagination and filtering
func (s *WordService) ListWords(ctx context.Context, params repository.ListWordsParams) ([]models.Word, int, error) {
	// Log the input parameters for debugging
	log.Printf("ListWords called with params: %+v", params)

	// Call repository method
	words, totalCount, err := s.repo.List(ctx, params)
	if err != nil {
		log.Printf("Error in ListWords: %v", err)
		return nil, 0, fmt.Errorf("failed to list words: %w", err)
	}

	// Log the results
	log.Printf("ListWords found %d words, total count: %d", len(words), totalCount)

	return words, totalCount, nil
}

// SearchWords provides a convenient method for searching words
func (s *WordService) SearchWords(ctx context.Context, query string, language string) ([]models.Word, int, error) {
	params := repository.ListWordsParams{
		Search:   query,
		Language: language,
		Page:     1,
		PageSize: 50, // Allow a larger default page size for search results
	}

	return s.repo.List(ctx, params)
}

// GetRandomWord retrieves a random word from the repository
func (s *WordService) GetRandomWord(ctx context.Context) (*models.Word, error) {
	return s.repo.GetRandomWord(ctx)
}

// GetWordsByGroupID retrieves all words associated with a specific group
func (s *WordService) GetWordsByGroupID(ctx context.Context, groupID int64) ([]models.Word, error) {
	return s.repo.GetWordsByGroupID(ctx, groupID)
}
