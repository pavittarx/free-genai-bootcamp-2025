package handlers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/pavittarx/lang-portal/pkg/models"
	"github.com/pavittarx/lang-portal/pkg/services"
)

// WordPaginationHandler manages HTTP handlers for paginated word retrieval
type WordPaginationHandler struct {
	wordService *services.WordPaginationService
}

// NewWordPaginationHandler creates a new instance of WordPaginationHandler
func NewWordPaginationHandler(wordService *services.WordPaginationService) *WordPaginationHandler {
	return &WordPaginationHandler{
		wordService: wordService,
	}
}

// ListWords handles the GET request to retrieve paginated words
func (h *WordPaginationHandler) ListWords(c echo.Context) error {
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
	
	// Get paginated words
	result, err := h.wordService.ListWords(ctx, req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to retrieve words: " + err.Error(),
		})
	}
	
	return c.JSON(http.StatusOK, result)
}

// ListWordsByGroup handles the GET request to retrieve paginated words in a group
func (h *WordPaginationHandler) ListWordsByGroup(c echo.Context) error {
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
	
	// Get paginated words for group
	result, err := h.wordService.ListWordsByGroup(ctx, groupID, req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to retrieve words for group: " + err.Error(),
		})
	}
	
	return c.JSON(http.StatusOK, result)
}

// GroupPaginationHandler manages HTTP handlers for paginated group retrieval
type GroupPaginationHandler struct {
	groupService *services.GroupPaginationService
}

// NewGroupPaginationHandler creates a new instance of GroupPaginationHandler
func NewGroupPaginationHandler(groupService *services.GroupPaginationService) *GroupPaginationHandler {
	return &GroupPaginationHandler{
		groupService: groupService,
	}
}

// ListGroups handles the GET request to retrieve paginated groups
func (h *GroupPaginationHandler) ListGroups(c echo.Context) error {
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
