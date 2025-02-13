package routes

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pavittarx/lang-portal/backend/pkg/handlers"
)

// RegisterRoutes sets up all routes for the application
func RegisterRoutes(e *echo.Echo, wordHandler *handlers.WordHandler, groupHandler *handlers.GroupHandler) {
	// Health check endpoints
	e.GET("/api", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status":  "healthy",
			"message": "Language Portal Backend is running successfully!",
		})
	})

	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status": "healthy",
		})
	})

	// Register word routes
	words := e.Group("/api/words")
	RegisterWordRoutes(words, wordHandler)

	// Register group routes
	groups := e.Group("/api/groups")
	RegisterGroupRoutes(groups, groupHandler)
}
