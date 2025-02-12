package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pavittarx/lang-portal/pkg/services"
)

// RandomWordHandler manages HTTP handlers for retrieving random words
type RandomWordHandler struct {
	wordService *services.WordService
}

// NewRandomWordHandler creates a new instance of RandomWordHandler
func NewRandomWordHandler(wordService *services.WordService) *RandomWordHandler {
	return &RandomWordHandler{
		wordService: wordService,
	}
}

// GetRandomWord handles the GET request to retrieve a random word
func (h *RandomWordHandler) GetRandomWord(c echo.Context) error {
	// Extract query parameters
	difficulty := c.QueryParam("difficulty")
	groupNamePtr := c.QueryParam("group")

	var groupID *int64
	if groupNamePtr != "" {
		// Find group ID by name
		groups, err := h.wordService.GroupRepo.List(c.Request().Context())
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve groups")
		}

		for _, g := range groups {
			if g.Name == groupNamePtr {
				groupIDVal := g.ID
				groupID = &groupIDVal
				break
			}
		}

		if groupID == nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Group not found")
		}
	}

	// Get random word
	word, err := h.wordService.GetRandomWord(c.Request().Context(), difficulty, groupID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve random word")
	}

	// Return the word
	return c.JSON(http.StatusOK, word)
}
