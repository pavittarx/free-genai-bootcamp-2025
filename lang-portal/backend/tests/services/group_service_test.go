package services_test

import (
	"context"
	"testing"

	"github.com/pavittarx/lang-portal/backend/pkg/models"
	"github.com/pavittarx/lang-portal/backend/pkg/repository"
	"github.com/pavittarx/lang-portal/backend/pkg/services"
	"github.com/pavittarx/lang-portal/backend/tests/testutils"
	"github.com/stretchr/testify/assert"
)

func setupGroupTest(t *testing.T) (*services.GroupService, func()) {
	db, cleanup, err := testutils.CreateTestDB()
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}

	groupRepo := repository.NewSQLiteGroupRepository(db)
	groupService := services.NewGroupService(groupRepo)

	return groupService, cleanup
}

func TestGroupService_CreateGroup(t *testing.T) {
	service, cleanup := setupGroupTest(t)
	defer cleanup()

	ctx := context.Background()

	tests := []struct {
		name    string
		group   *models.Group
		wantErr bool
	}{
		{
			name: "valid group",
			group: &models.Group{
				Name: "Travel Words",
			},
			wantErr: false,
		},
		{
			name: "empty group name",
			group: &models.Group{
				Name: "",
			},
			wantErr: true,
		},
		{
			name: "invalid group name",
			group: &models.Group{
				Name: "A",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.CreateGroup(ctx, tt.group)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotZero(t, tt.group.ID)
			}
		})
	}
}

func TestGroupService_GetGroupByID(t *testing.T) {
	service, cleanup := setupGroupTest(t)
	defer cleanup()

	ctx := context.Background()

	// First create a group
	group := &models.Group{
		Name: "Travel Words",
	}
	err := service.CreateGroup(ctx, group)
	assert.NoError(t, err)

	// Retrieve the group
	retrievedGroup, err := service.GetGroupByID(ctx, group.ID)
	assert.NoError(t, err)
	assert.Equal(t, group.Name, retrievedGroup.Name)

	// Try to get non-existing group
	_, err = service.GetGroupByID(ctx, 999)
	assert.Error(t, err)

	// Try to get invalid group ID
	_, err = service.GetGroupByID(ctx, 0)
	assert.Error(t, err)
}

func TestGroupService_UpdateGroup(t *testing.T) {
	service, cleanup := setupGroupTest(t)
	defer cleanup()

	ctx := context.Background()

	// First create a group
	group := &models.Group{
		Name: "Travel Words",
	}
	err := service.CreateGroup(ctx, group)
	assert.NoError(t, err)

	// Update the group
	group.Name = "Updated Travel Words"
	err = service.UpdateGroup(ctx, group)
	assert.NoError(t, err)

	// Retrieve and verify
	retrievedGroup, err := service.GetGroupByID(ctx, group.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Updated Travel Words", retrievedGroup.Name)

	// Try to update with invalid group name
	group.Name = ""
	err = service.UpdateGroup(ctx, group)
	assert.Error(t, err)

	// Try to update with invalid group ID
	invalidGroup := &models.Group{
		ID:   0,
		Name: "Invalid Group",
	}
	err = service.UpdateGroup(ctx, invalidGroup)
	assert.Error(t, err)
}

func TestGroupService_DeleteGroup(t *testing.T) {
	service, cleanup := setupGroupTest(t)
	defer cleanup()

	ctx := context.Background()

	// First create a group
	group := &models.Group{
		Name: "Travel Words",
	}
	err := service.CreateGroup(ctx, group)
	assert.NoError(t, err)

	// Delete the group
	err = service.DeleteGroup(ctx, group.ID)
	assert.NoError(t, err)

	// Try to retrieve deleted group
	_, err = service.GetGroupByID(ctx, group.ID)
	assert.Error(t, err)

	// Try to delete non-existing group
	err = service.DeleteGroup(ctx, 999)
	assert.Error(t, err)

	// Try to delete with invalid group ID
	err = service.DeleteGroup(ctx, 0)
	assert.Error(t, err)
}

func TestGroupService_ListGroups(t *testing.T) {
	service, cleanup := setupGroupTest(t)
	defer cleanup()

	ctx := context.Background()

	// Create multiple groups
	groups := []models.Group{
		{Name: "Travel Words"},
		{Name: "Food Vocabulary"},
		{Name: "Business English"},
	}

	for i := range groups {
		err := service.CreateGroup(ctx, &groups[i])
		assert.NoError(t, err)
	}

	// Test listing all groups
	listedGroups, total, err := service.ListGroups(ctx, 1, 10, "")
	assert.NoError(t, err)
	assert.Equal(t, len(groups), total)
	assert.Equal(t, len(groups), len(listedGroups))

	// Test pagination
	listedGroups, total, err = service.ListGroups(ctx, 1, 2, "")
	assert.NoError(t, err)
	assert.Equal(t, len(groups), total)
	assert.Equal(t, 2, len(listedGroups))

	// Test search
	listedGroups, total, err = service.ListGroups(ctx, 1, 10, "Travel")
	assert.NoError(t, err)
	assert.Equal(t, 1, total)
	assert.Equal(t, 1, len(listedGroups))
	assert.Equal(t, "Travel Words", listedGroups[0].Name)

	// Test invalid pagination
	listedGroups, total, err = service.ListGroups(ctx, 0, 0, "")
	assert.NoError(t, err)
	assert.Equal(t, len(groups), total)
	assert.LessOrEqual(t, len(listedGroups), 10) // Default page size is 10
	assert.GreaterOrEqual(t, len(listedGroups), 1) // But should have at least 1 group
}
