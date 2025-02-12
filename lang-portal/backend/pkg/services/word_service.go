package services

import (
	"context"
	"errors"
	"log"

	"github.com/pavittarx/lang-portal/pkg/models"
	"github.com/pavittarx/lang-portal/pkg/repository"
)

// WordService provides business logic for word-related operations
type WordService struct {
	wordRepo  repository.WordRepository
	groupRepo repository.GroupRepository
}

// NewWordService creates a new instance of WordService
func NewWordService(wordRepo repository.WordRepository, groupRepo repository.GroupRepository) *WordService {
	return &WordService{
		wordRepo:  wordRepo,
		groupRepo: groupRepo,
	}
}

// GetRandomWord retrieves a random word with optional filtering
func (s *WordService) GetRandomWord(ctx context.Context, difficulty string, groupName *string) (*models.Word, error) {
	// Validate difficulty
	validDifficulties := map[string]bool{
		"easy":   true,
		"medium": true,
		"hard":   true,
	}
	
	if difficulty != "" && !validDifficulties[difficulty] {
		return nil, errors.New("invalid difficulty level")
	}
	
	// If group name is provided, first find the group ID
	var groupID *int64
	if groupName != nil {
		groups, err := s.groupRepo.List(ctx)
		if err != nil {
			log.Printf("Error listing groups: %v", err)
			return nil, err
		}
		
		for _, g := range groups {
			if g.Name == *groupName {
				groupIDVal := g.ID
				groupID = &groupIDVal
				break
			}
		}
		
		if groupID == nil {
			return nil, errors.New("group not found")
		}
	}
	
	// Prepare filter
	filter := &models.RandomWordFilter{
		Difficulty: difficulty,
		GroupID:    groupID,
	}
	
	// Get random word
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

// ListWordsByGroup retrieves all words in a specific group
func (s *WordService) ListWordsByGroup(ctx context.Context, groupName string) ([]models.Word, error) {
	// First, find the group ID
	groups, err := s.groupRepo.List(ctx)
	if err != nil {
		log.Printf("Error listing groups: %v", err)
		return nil, err
	}
	
	var groupID int64
	found := false
	for _, g := range groups {
		if g.Name == groupName {
			groupID = g.ID
			found = true
			break
		}
	}
	
	if !found {
		return nil, errors.New("group not found")
	}
	
	// Get words for the group
	words, err := s.wordRepo.ListByGroup(ctx, groupID)
	if err != nil {
		log.Printf("Error getting words by group: %v", err)
		return nil, err
	}
	
	return words, nil
}
