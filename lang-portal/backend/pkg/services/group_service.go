package services

import (
	"context"
	"fmt"

	"github.com/pavittarx/lang-portal/backend/pkg/models"
	"github.com/pavittarx/lang-portal/backend/pkg/repository"
)

// GroupService handles business logic for group operations
type GroupService struct {
	groupRepo *repository.SQLiteGroupRepository
}

// NewGroupService creates a new instance of GroupService
func NewGroupService(repo *repository.SQLiteGroupRepository) *GroupService {
	return &GroupService{groupRepo: repo}
}

// CreateGroup creates a new group with validation
func (s *GroupService) CreateGroup(ctx context.Context, group *models.Group) error {
	// Validate the group
	if err := group.Validate(); err != nil {
		return fmt.Errorf("invalid group: %w", err)
	}

	// Sanitize the group name and description
	group.Sanitize()

	// Create the group in the repository
	return s.groupRepo.Create(ctx, group)
}

// GetGroupByID retrieves a group by its ID
func (s *GroupService) GetGroupByID(ctx context.Context, id int64) (*models.Group, error) {
	// Validate ID
	if id <= 0 {
		return nil, fmt.Errorf("invalid group ID: %d", id)
	}

	// Retrieve the group from the repository
	return s.groupRepo.GetByID(ctx, id)
}

// UpdateGroup updates an existing group
func (s *GroupService) UpdateGroup(ctx context.Context, group *models.Group) error {
	// Validate the group
	if err := group.Validate(); err != nil {
		return fmt.Errorf("invalid group: %w", err)
	}

	// Validate ID
	if group.ID <= 0 {
		return fmt.Errorf("invalid group ID: %d", group.ID)
	}

	// Sanitize the group name and description
	group.Sanitize()

	// Update the group in the repository
	return s.groupRepo.Update(ctx, group)
}

// DeleteGroup removes a group by its ID
func (s *GroupService) DeleteGroup(ctx context.Context, id int64) error {
	// Validate ID
	if id <= 0 {
		return fmt.Errorf("invalid group ID: %d", id)
	}

	// Delete the group from the repository
	return s.groupRepo.Delete(ctx, id)
}

// ListGroups retrieves groups with pagination and optional search
func (s *GroupService) ListGroups(ctx context.Context, page, pageSize int, search string) ([]models.Group, int, error) {
	// Validate pagination parameters
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10 // Default page size
	}

	// List groups from the repository
	return s.groupRepo.List(ctx, page, pageSize, search)
}

// GetGroups retrieves a paginated list of groups
func (s *GroupService) GetGroups(ctx context.Context, page, pageSize int) ([]*models.Group, int64, error) {
	// Retrieve groups with pagination
	groups, totalCount, err := s.groupRepo.List(ctx, page, pageSize, "")
	if err != nil {
		return nil, 0, fmt.Errorf("failed to retrieve groups: %v", err)
	}

	// Convert to pointer slice if needed
	groupPtrs := make([]*models.Group, len(groups))
	for i := range groups {
		groupPtrs[i] = &groups[i]
	}

	return groupPtrs, int64(totalCount), nil
}
