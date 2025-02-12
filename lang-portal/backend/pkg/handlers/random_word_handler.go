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
			"error": "Failed to retrieve random word: " + err.Error(),
		})
	}
	
	return c.JSON(http.StatusOK, word)
}
