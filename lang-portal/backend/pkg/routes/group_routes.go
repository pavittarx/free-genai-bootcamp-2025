package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/pavittarx/lang-portal/backend/pkg/handlers"
)

// RegisterGroupRoutes sets up the routes for group-related operations
func RegisterGroupRoutes(groups *echo.Group, groupHandler *handlers.GroupHandler) {
	// Create a new group
	groups.POST("", groupHandler.CreateGroup)

	// Get a group by ID
	groups.GET("/:id", groupHandler.GetGroupByID)

	// Update an existing group
	groups.PUT("/:id", groupHandler.UpdateGroup)

	// Delete a group
	groups.DELETE("/:id", groupHandler.DeleteGroup)

	// List groups with optional pagination and search
	groups.GET("", groupHandler.ListGroups)
}
