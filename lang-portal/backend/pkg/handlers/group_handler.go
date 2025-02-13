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
