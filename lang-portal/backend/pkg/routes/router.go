package routes

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pavittarx/lang-portal/backend/pkg/handlers"
)

// RegisterRoutes sets up all routes for the application
func RegisterRoutes(e *echo.Echo, 
	wordHandler *handlers.WordHandler, 
	groupHandler *handlers.GroupHandler, 
	studyActivityHandler *handlers.StudyActivityHandler,
	sessionHandler *handlers.SessionHandler) {
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
	e.GET("/api/words", wordHandler.ListWords)
	e.GET("/api/words/:id", wordHandler.GetWordByID)
	e.GET("/api/words/random", wordHandler.GetRandomWord)
	e.POST("/api/words", wordHandler.CreateWord)
	e.PUT("/api/words/:id", wordHandler.UpdateWord)
	e.DELETE("/api/words/:id", wordHandler.DeleteWord)

	// Register group routes
	e.GET("/api/groups", groupHandler.ListGroups)
	e.GET("/api/groups/:id", groupHandler.GetGroupByID)
	e.GET("/api/words/groups/:group-id", wordHandler.GetWordsByGroupID)

	// Register study activities routes
	e.GET("/api/study-activities", studyActivityHandler.GetStudyActivities)

	// Session routes
	e.GET("/api/sessions", sessionHandler.GetSessions)
	e.GET("/api/sessions/:id", sessionHandler.GetSessionByID)
	e.POST("/api/sessions", sessionHandler.CreateSession)
	e.PUT("/api/sessions/:id", sessionHandler.UpdateSession)
	e.DELETE("/api/sessions/:id", sessionHandler.DeleteSession)
}
