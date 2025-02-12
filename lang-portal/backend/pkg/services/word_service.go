package services

import (
	"context"
	"log"

	"github.com/pavittarx/lang-portal/pkg/models"
	"github.com/pavittarx/lang-portal/pkg/repository"
)

// WordService provides business logic for word-related operations
type WordService struct {
	wordRepo  repository.WordRepository
	GroupRepo repository.GroupRepository
}

// NewWordService creates a new instance of WordService
func NewWordService(wordRepo repository.WordRepository, groupRepo repository.GroupRepository) *WordService {
	return &WordService{
		wordRepo:  wordRepo,
		GroupRepo: groupRepo,
	}
}

// GetRandomWord retrieves a random word with optional filtering
func (s *WordService) GetRandomWord(ctx context.Context, difficulty string, groupID *int64) (*models.Word, error) {
	// Prepare random word filter
	filter := &models.RandomWordFilter{
		Difficulty: difficulty,
		GroupID:    groupID,
	}

	// Retrieve random word
	word, err := s.wordRepo.GetRandom(ctx, filter)
	if err != nil {
		log.Printf("Error getting random word: %v", err)
		return nil, err
	}

	return word, nil
}

// GetWordDetails retrieves detailed information about a specific word
func (s *WordService) GetWordDetails(ctx context.Context, wordID int64) (*models.Word, error) {
	word, err := s.wordRepo.GetByID(ctx, wordID)
	if err != nil {
		log.Printf("Error getting word details: %v", err)
		return nil, err
	}
	
	return word, nil
}

// ListWordsByGroup retrieves words in a specific group
func (s *WordService) ListWordsByGroup(ctx context.Context, groupID int64) ([]models.Word, error) {
	// Retrieve words in the group
	words, err := s.wordRepo.ListByGroup(ctx, groupID)
	if err != nil {
		log.Printf("Error listing words by group: %v", err)
		return nil, err
	}

	return words, nil
}
