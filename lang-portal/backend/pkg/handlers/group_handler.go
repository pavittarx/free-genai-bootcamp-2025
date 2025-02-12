package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pavittarx/lang-portal/pkg/services"
)

// GroupHandler manages HTTP handlers for group-related operations
type GroupHandler struct {
	wordService *services.WordService
}

// NewGroupHandler creates a new instance of GroupHandler
func NewGroupHandler(wordService *services.WordService) *GroupHandler {
	return &GroupHandler{
		wordService: wordService,
	}
}

// ListWordsByGroup handles the GET request to retrieve words in a specific group
func (h *GroupHandler) ListWordsByGroup(c echo.Context) error {
	ctx := c.Request().Context()
	
	// Get group name from query parameter
	groupName := c.QueryParam("name")
	if groupName == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Group name is required",
		})
	}
	
	// Get words for the group
	words, err := h.wordService.ListWordsByGroup(ctx, groupName)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "Failed to retrieve words for group: " + err.Error(),
		})
	}
	
	return c.JSON(http.StatusOK, words)
}
