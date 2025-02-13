package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/pavittarx/lang-portal/backend/pkg/models"
	"github.com/pavittarx/lang-portal/backend/pkg/repository"
	"github.com/pavittarx/lang-portal/backend/pkg/services"
)

// GroupHandler handles HTTP requests related to groups
type GroupHandler struct {
	groupService *services.GroupService
	groupRepo    *repository.SQLiteGroupRepository
}

// NewGroupHandler creates a new instance of GroupHandler
func NewGroupHandler(service *services.GroupService, repo *repository.SQLiteGroupRepository) *GroupHandler {
	return &GroupHandler{
		groupService: service,
		groupRepo:    repo,
	}
}

// CreateGroup handles the creation of a new group
// @Summary Create a new group
// @Description Create a new group with a name and optional description
// @Tags groups
// @Accept json
// @Produce json
// @Param group body models.Group true "Group details"
// @Success 201 {object} map[string]interface{} "Group created successfully"
// @Failure 400 {object} map[string]string "Invalid request payload or validation error"
// @Router /groups [post]
func (h *GroupHandler) CreateGroup(c echo.Context) error {
	group := &models.Group{}
	if err := c.Bind(group); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	if err := h.groupService.CreateGroup(c.Request().Context(), group); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"group": group,
		"description": group.Description,
	})
}

// GetGroupByID retrieves a group by its ID
// @Summary Get a group by ID
// @Description Retrieve a group's details using its unique identifier
// @Tags groups
// @Produce json
// @Param id path int true "Group ID"
// @Success 200 {object} map[string]interface{} "Group retrieved successfully"
// @Failure 400 {object} map[string]string "Invalid group ID"
// @Failure 404 {object} map[string]string "Group not found"
// @Router /groups/{id} [get]
func (h *GroupHandler) GetGroupByID(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid group ID"})
	}

	group, err := h.groupService.GetGroupByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"group": group,
		"description": group.Description,
	})
}

// UpdateGroup updates an existing group
// @Summary Update a group
// @Description Update an existing group's details by its ID
// @Tags groups
// @Accept json
// @Produce json
// @Param id path int true "Group ID"
// @Param group body models.Group true "Updated group details"
// @Success 200 {object} map[string]interface{} "Group updated successfully"
// @Failure 400 {object} map[string]string "Invalid request payload or validation error"
// @Failure 404 {object} map[string]string "Group not found"
// @Router /groups/{id} [put]
func (h *GroupHandler) UpdateGroup(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid group ID"})
	}

	group := &models.Group{}
	if err := c.Bind(group); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	group.ID = id
	if err := h.groupService.UpdateGroup(c.Request().Context(), group); err != nil {
		if err.Error() == fmt.Sprintf("group with ID %d not found", id) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"group": group,
		"description": group.Description,
	})
}

// DeleteGroup removes a group by its ID
// @Summary Delete a group
// @Description Delete an existing group by its unique identifier
// @Tags groups
// @Produce json
// @Param id path int true "Group ID"
// @Success 200 {object} map[string]string "Group deleted successfully"
// @Failure 400 {object} map[string]string "Invalid group ID"
// @Failure 404 {object} map[string]string "Group not found"
// @Router /groups/{id} [delete]
func (h *GroupHandler) DeleteGroup(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid group ID"})
	}

	if err := h.groupService.DeleteGroup(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Group deleted successfully"})
}

// ListGroups retrieves a list of groups
// @Summary List groups
// @Description Retrieve a paginated list of groups with optional search
// @Tags groups
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param pageSize query int false "Number of items per page" default(10)
// @Param search query string false "Search term"
// @Success 200 {object} map[string]interface{} "Groups retrieved successfully"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /groups [get]
func (h *GroupHandler) ListGroups(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	pageSize, _ := strconv.Atoi(c.QueryParam("pageSize"))
	search := c.QueryParam("search")

	groups, total, err := h.groupService.ListGroups(c.Request().Context(), page, pageSize, search)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	var groupResponses []map[string]interface{}
	for _, group := range groups {
		groupResponses = append(groupResponses, map[string]interface{}{
			"group": group,
			"description": group.Description,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"groups": groupResponses,
		"total":  total,
	})
}
