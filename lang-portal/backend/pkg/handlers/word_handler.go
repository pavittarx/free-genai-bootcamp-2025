package handlers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/pavittarx/lang-portal/pkg/services"
)

// WordHandler manages HTTP handlers for word-related endpoints
type WordHandler struct {
	wordService *services.WordService
}

// NewWordHandler creates a new instance of WordHandler
func NewWordHandler(wordService *services.WordService) *WordHandler {
	return &WordHandler{
		wordService: wordService,
	}
}

// GetRandomWord handles the GET request to retrieve a random word
func (h *WordHandler) GetRandomWord(c echo.Context) error {
	ctx := c.Request().Context()
	
	// Get optional query parameters
	difficulty := c.QueryParam("difficulty")
	groupName := c.QueryParam("group")
	
	var groupNamePtr *string
	if groupName != "" {
		groupNamePtr = &groupName
	}
	
	// Get random word
	word, err := h.wordService.GetRandomWord(ctx, difficulty, groupNamePtr)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to retrieve random word",
		})
	}
	
	return c.JSON(http.StatusOK, word)
}

// GetWordDetails handles the GET request to retrieve details of a specific word
func (h *WordHandler) GetWordDetails(c echo.Context) error {
	ctx := c.Request().Context()
	
	// Get word ID from URL parameter
	wordIDStr := c.Param("id")
	wordID, err := strconv.ParseInt(wordIDStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid word ID",
		})
	}
	
	// Get word details
	word, err := h.wordService.GetWordDetails(ctx, wordID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "Word not found",
		})
	}
	
	return c.JSON(http.StatusOK, word)
}

// ListWordsByGroup handles the GET request to retrieve words in a specific group
func (h *WordHandler) ListWordsByGroup(c echo.Context) error {
	ctx := c.Request().Context()
	
	// Get group name from URL parameter
	groupName := c.Param("group")
	
	// Get words in the group
	words, err := h.wordService.ListWordsByGroup(ctx, groupName)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to retrieve words for group",
		})
	}
	
	return c.JSON(http.StatusOK, words)
}
