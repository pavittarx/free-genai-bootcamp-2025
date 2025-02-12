package handlers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/pavittarx/lang-portal/pkg/models"
	"github.com/pavittarx/lang-portal/pkg/services"
)

// GroupHandler manages HTTP handlers for group-related operations
type GroupHandler struct {
	groupService *services.GroupService
}

// NewGroupHandler creates a new instance of GroupHandler
func NewGroupHandler(groupService *services.GroupService) *GroupHandler {
	return &GroupHandler{
		groupService: groupService,
	}
}

// ListGroups handles the GET request to retrieve all groups
func (h *GroupHandler) ListGroups(c echo.Context) error {
	ctx := c.Request().Context()

	// Parse pagination parameters
	page, _ := strconv.Atoi(c.QueryParam("page"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	sort := c.QueryParam("sort")
	order := c.QueryParam("order")

	// Create pagination request
	req := models.PaginationRequest{
		Page:   page,
		Limit:  limit,
		Sort:   sort,
		Order:  order,
	}

	// Get paginated groups
	result, err := h.groupService.ListGroups(ctx, req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to retrieve groups: " + err.Error(),
		})
	}

	return c.JSON(http.StatusOK, result)
}

// GetGroupDetails handles the GET request to retrieve details of a specific group
func (h *GroupHandler) GetGroupDetails(c echo.Context) error {
	ctx := c.Request().Context()

	// Get group ID from URL parameter
	groupIDStr := c.Param("id")
	groupID, err := strconv.ParseInt(groupIDStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid group ID",
		})
	}

	// Get group details
	group, err := h.groupService.GetGroupDetails(ctx, groupID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to retrieve group details: " + err.Error(),
		})
	}

	return c.JSON(http.StatusOK, group)
}

// ListGroupWords handles the GET request to retrieve words in a group
func (h *GroupHandler) ListGroupWords(c echo.Context) error {
	ctx := c.Request().Context()

	// Get group ID from URL parameter
	groupIDStr := c.Param("id")
	groupID, err := strconv.ParseInt(groupIDStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid group ID",
		})
	}

	// Parse pagination parameters
	page, _ := strconv.Atoi(c.QueryParam("page"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	sort := c.QueryParam("sort")
	order := c.QueryParam("order")

	// Create pagination request
	req := models.PaginationRequest{
		Page:   page,
		Limit:  limit,
		Sort:   sort,
		Order:  order,
	}

	// Get paginated words in the group
	result, err := h.groupService.ListGroupWords(ctx, groupID, req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to retrieve words for group: " + err.Error(),
		})
	}

	return c.JSON(http.StatusOK, result)
}
