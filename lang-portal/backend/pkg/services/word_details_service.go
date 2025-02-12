package services

import (
	"context"
	"log"

	"github.com/pavittarx/lang-portal/pkg/models"
	"github.com/pavittarx/lang-portal/pkg/repository"
)

// WordDetailsService provides business logic for retrieving word details
type WordDetailsService struct {
	wordRepo repository.WordRepository
}

// NewWordDetailsService creates a new instance of WordDetailsService
func NewWordDetailsService(wordRepo repository.WordRepository) *WordDetailsService {
	return &WordDetailsService{
		wordRepo: wordRepo,
	}
}

// GetWordDetails retrieves detailed information about a specific word
func (s *WordDetailsService) GetWordDetails(ctx context.Context, wordID int64) (*models.Word, error) {
	word, err := s.wordRepo.GetByID(ctx, wordID)
	if err != nil {
		log.Printf("Error getting word details: %v", err)
		return nil, err
	}
	
	return word, nil
}
