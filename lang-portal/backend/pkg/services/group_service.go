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
func NewGroupService(groupRepo repository.GroupRepository, wordRepo repository.WordRepository) *GroupService {
	return &GroupService{
		groupRepo: groupRepo,
		wordRepo:  wordRepo,
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

// ListGroups retrieves paginated list of groups
func (s *GroupService) ListGroups(ctx context.Context, req models.PaginationRequest) (*models.PaginationResponse, error) {
	// Validate pagination request
	req.Validate()

	// Retrieve paginated groups
	groups, total, err := s.groupRepo.ListPaginated(ctx, req)
	if err != nil {
		log.Printf("Error listing groups: %v", err)
		return nil, err
	}

	// Create pagination response
	return &models.PaginationResponse{
		Data:       groups,
		Total:      total,
		Page:       req.Page,
		Limit:      req.Limit,
		TotalPages: int((total + int64(req.Limit) - 1) / int64(req.Limit)),
	}, nil
}

// GetGroupDetails retrieves details of a specific group
func (s *GroupService) GetGroupDetails(ctx context.Context, groupID int64) (*models.Group, error) {
	// Retrieve group details
	group, err := s.groupRepo.GetByID(ctx, groupID)
	if err != nil {
		log.Printf("Error getting group details: %v", err)
		return nil, err
	}

	return group, nil
}

// ListGroupWords retrieves paginated words in a specific group
func (s *GroupService) ListGroupWords(ctx context.Context, groupID int64, req models.PaginationRequest) (*models.PaginationResponse, error) {
	// Validate pagination request
	req.Validate()

	// Retrieve words in the group
	words, err := s.wordRepo.ListByGroup(ctx, groupID)
	if err != nil {
		log.Printf("Error listing words by group: %v", err)
		return nil, err
	}

	// Calculate total and pagination
	total := int64(len(words))
	startIndex := int64((req.Page - 1) * req.Limit)
	endIndex := startIndex + int64(req.Limit)

	// Slice words based on pagination
	if startIndex > total {
		startIndex = 0
	}
	if endIndex > total {
		endIndex = total
	}

	var paginatedWords []models.Word
	if total > 0 {
		paginatedWords = words[startIndex:endIndex]
	}

	// Create pagination response
	return &models.PaginationResponse{
		Data:       paginatedWords,
		Total:      total,
		Page:       req.Page,
		Limit:      req.Limit,
		TotalPages: int((total + int64(req.Limit) - 1) / int64(req.Limit)),
	}, nil
}
