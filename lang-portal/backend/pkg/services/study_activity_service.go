package services

import (
	"context"

	"github.com/pavittarx/lang-portal/backend/pkg/models"
	"github.com/pavittarx/lang-portal/backend/pkg/repository"
)

// StudyActivityService provides business logic for study activity operations
type StudyActivityService struct {
	repo *repository.StudyActivityRepository
}

// NewStudyActivityService creates a new instance of StudyActivityService
func NewStudyActivityService(repo *repository.StudyActivityRepository) *StudyActivityService {
	return &StudyActivityService{
		repo: repo,
	}
}

// GetStudyActivities retrieves all available study activities
func (s *StudyActivityService) GetStudyActivities(ctx context.Context) ([]models.StudyActivity, error) {
	return s.repo.GetAll(ctx)
}
