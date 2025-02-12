package services

import (
	"context"
	"errors"
	"log"

	"github.com/pavittarx/lang-portal/pkg/models"
	"github.com/pavittarx/lang-portal/pkg/repository"
)

// RandomWordService provides business logic for retrieving random words
type RandomWordService struct {
	wordRepo  repository.WordRepository
	groupRepo repository.GroupRepository
}

// NewRandomWordService creates a new instance of RandomWordService
func NewRandomWordService(wordRepo repository.WordRepository, groupRepo repository.GroupRepository) *RandomWordService {
	return &RandomWordService{
		wordRepo:  wordRepo,
		groupRepo: groupRepo,
	}
}

// GetRandomWord retrieves a random word with optional filtering
func (s *RandomWordService) GetRandomWord(ctx context.Context, difficulty string, groupName *string) (*models.Word, error) {
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
