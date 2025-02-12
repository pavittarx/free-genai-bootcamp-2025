package handlers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/pavittarx/lang-portal/pkg/services"
)

// WordHandler manages HTTP handlers for word-related operations
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

	// Parse optional difficulty and group parameters
	difficulty := c.QueryParam("difficulty")
	groupIDStr := c.QueryParam("group")

	var groupID *int64
	if groupIDStr != "" {
		parsedGroupID, err := strconv.ParseInt(groupIDStr, 10, 64)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "Invalid group ID",
			})
		}
		groupID = &parsedGroupID
	}

	// Get random word
	word, err := h.wordService.GetRandomWord(ctx, difficulty, groupID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to retrieve random word: " + err.Error(),
		})
	}

	return c.JSON(http.StatusOK, word)
}

// GetWordDetails handles the GET request to retrieve word details
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
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to retrieve word details: " + err.Error(),
		})
	}

	return c.JSON(http.StatusOK, word)
}

// ListWordsByGroup handles the GET request to retrieve words in a group
func (h *WordHandler) ListWordsByGroup(c echo.Context) error {
	ctx := c.Request().Context()

	// Get group ID from URL parameter
	groupIDStr := c.Param("group")
	groupID, err := strconv.ParseInt(groupIDStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid group ID",
		})
	}

	// Get words in the group
	words, err := h.wordService.ListWordsByGroup(ctx, groupID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to retrieve words for group: " + err.Error(),
		})
	}

	return c.JSON(http.StatusOK, words)
}
