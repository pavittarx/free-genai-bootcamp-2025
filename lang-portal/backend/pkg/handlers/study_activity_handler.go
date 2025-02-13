package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pavittarx/lang-portal/backend/pkg/services"
)

// StudyActivityHandler handles HTTP requests for study activities
type StudyActivityHandler struct {
	service *services.StudyActivityService
}

// NewStudyActivityHandler creates a new instance of StudyActivityHandler
func NewStudyActivityHandler(service *services.StudyActivityService) *StudyActivityHandler {
	return &StudyActivityHandler{
		service: service,
	}
}

// GetStudyActivities retrieves all available study activities
func (h *StudyActivityHandler) GetStudyActivities(c echo.Context) error {
	ctx := c.Request().Context()
	activities, err := h.service.GetStudyActivities(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to retrieve study activities",
		})
	}
	return c.JSON(http.StatusOK, activities)
}
