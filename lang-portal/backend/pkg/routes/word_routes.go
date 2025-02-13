package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/pavittarx/lang-portal/backend/pkg/handlers"
)

// RegisterWordRoutes sets up the routes for word-related operations
func RegisterWordRoutes(words *echo.Group, wordHandler *handlers.WordHandler) {
	// Create a new word
	words.POST("", wordHandler.CreateWord)

	// List all words
	words.GET("", wordHandler.ListWords)

	// Search words
	words.GET("/search", wordHandler.SearchWords)

	// Get a random word
	words.GET("/random", wordHandler.GetRandomWord)

	// Get a word by ID
	words.GET("/:id", wordHandler.GetWordByID)

	// Update an existing word
	words.PUT("/:id", wordHandler.UpdateWord)

	// Delete a word
	words.DELETE("/:id", wordHandler.DeleteWord)
}
