package handlers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/pavittarx/lang-portal/pkg/services"
)

// WordDetailsHandler manages HTTP handlers for word details
type WordDetailsHandler struct {
	wordService *services.WordService
}

// NewWordDetailsHandler creates a new instance of WordDetailsHandler
func NewWordDetailsHandler(wordService *services.WordService) *WordDetailsHandler {
	return &WordDetailsHandler{
		wordService: wordService,
	}
}

// GetWordDetails handles the GET request to retrieve details of a specific word
func (h *WordDetailsHandler) GetWordDetails(c echo.Context) error {
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
			"error": "Word not found: " + err.Error(),
		})
	}
	
	return c.JSON(http.StatusOK, word)
}
