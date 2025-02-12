package repository_test

import (
	"context"
	"testing"

	"github.com/pavittarx/lang-portal/backend/pkg/models"
	"github.com/pavittarx/lang-portal/backend/pkg/repository"
	"github.com/pavittarx/lang-portal/backend/tests/testutils"
	"github.com/stretchr/testify/assert"
)

func setupGroupTest(t *testing.T) (*repository.SQLiteGroupRepository, func()) {
	db, cleanup, err := testutils.CreateTestDB()
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}

	repo := repository.NewSQLiteGroupRepository(db)
	return repo, cleanup
}

func TestGroupRepository_Create(t *testing.T) {
	repo, cleanup := setupGroupTest(t)
	defer cleanup()

	ctx := context.Background()

	group := &models.Group{
		Name: "Travel Words",
	}

	err := repo.Create(ctx, group)
	assert.NoError(t, err)
	assert.NotZero(t, group.ID)
}

func TestGroupRepository_CreateValidation(t *testing.T) {
	repo, cleanup := setupGroupTest(t)
	defer cleanup()

	ctx := context.Background()

	tests := []struct {
		name    string
		group   *models.Group
		wantErr bool
	}{
		{
			name:    "empty group name",
			group:   &models.Group{Name: ""},
			wantErr: true,
		},
		{
			name:    "invalid group name",
			group:   &models.Group{Name: "A"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Create(ctx, tt.group)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGroupRepository_GetByID(t *testing.T) {
	repo, cleanup := setupGroupTest(t)
	defer cleanup()

	ctx := context.Background()

	// First create a group
	group := &models.Group{
		Name: "Travel Words",
	}
	err := repo.Create(ctx, group)
	assert.NoError(t, err)

	// Then retrieve it
	retrievedGroup, err := repo.GetByID(ctx, group.ID)
	assert.NoError(t, err)
	assert.Equal(t, group.Name, retrievedGroup.Name)
	assert.NotZero(t, retrievedGroup.CreatedAt)

	// Try to get non-existing group
	_, err = repo.GetByID(ctx, 999)
	assert.Error(t, err)
}

func TestGroupRepository_Update(t *testing.T) {
	repo, cleanup := setupGroupTest(t)
	defer cleanup()

	ctx := context.Background()

	// First create a group
	group := &models.Group{
		Name: "Travel Words",
	}
	err := repo.Create(ctx, group)
	assert.NoError(t, err)

	// Update the group
	group.Name = "Updated Travel Words"
	err = repo.Update(ctx, group)
	assert.NoError(t, err)

	// Retrieve and verify
	updatedGroup, err := repo.GetByID(ctx, group.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Updated Travel Words", updatedGroup.Name)

	// Try to update non-existing group
	nonExistingGroup := &models.Group{
		ID:   999,
		Name: "Non-existing Group",
	}
	err = repo.Update(ctx, nonExistingGroup)
	assert.Error(t, err)
}

func TestGroupRepository_Delete(t *testing.T) {
	repo, cleanup := setupGroupTest(t)
	defer cleanup()

	ctx := context.Background()

	// First create a group
	group := &models.Group{
		Name: "Travel Words",
	}
	err := repo.Create(ctx, group)
	assert.NoError(t, err)

	// Delete the group
	err = repo.Delete(ctx, group.ID)
	assert.NoError(t, err)

	// Try to retrieve deleted group
	_, err = repo.GetByID(ctx, group.ID)
	assert.Error(t, err)

	// Try to delete non-existing group
	err = repo.Delete(ctx, 999)
	assert.Error(t, err)
}

func TestGroupRepository_List(t *testing.T) {
	repo, cleanup := setupGroupTest(t)
	defer cleanup()

	ctx := context.Background()

	// Create multiple groups
	groups := []models.Group{
		{Name: "Travel Words"},
		{Name: "Food Vocabulary"},
		{Name: "Business English"},
	}

	for i := range groups {
		err := repo.Create(ctx, &groups[i])
		assert.NoError(t, err)
	}

	// Test listing all groups
	listedGroups, total, err := repo.List(ctx, 1, 10, "")
	assert.NoError(t, err)
	assert.Equal(t, len(groups), total)
	assert.Equal(t, len(groups), len(listedGroups))

	// Test pagination
	listedGroups, total, err = repo.List(ctx, 1, 2, "")
	assert.NoError(t, err)
	assert.Equal(t, len(groups), total)
	assert.Equal(t, 2, len(listedGroups))

	// Test search
	listedGroups, total, err = repo.List(ctx, 1, 10, "Travel")
	assert.NoError(t, err)
	assert.Equal(t, 1, total)
	assert.Equal(t, 1, len(listedGroups))
	assert.Equal(t, "Travel Words", listedGroups[0].Name)
}
