package services

import (
	"context"
	"errors"
	"log"

	"github.com/pavittarx/lang-portal/pkg/models"
	"github.com/pavittarx/lang-portal/pkg/repository"
)

// GroupService provides business logic for group-related operations
type GroupService struct {
	wordRepo  repository.WordRepository
	groupRepo repository.GroupRepository
}

// NewGroupService creates a new instance of GroupService
func NewGroupService(wordRepo repository.WordRepository, groupRepo repository.GroupRepository) *GroupService {
	return &GroupService{
		wordRepo:  wordRepo,
		groupRepo: groupRepo,
	}
}

// ListWordsByGroup retrieves all words in a specific group
func (s *GroupService) ListWordsByGroup(ctx context.Context, groupName string) ([]models.Word, error) {
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
