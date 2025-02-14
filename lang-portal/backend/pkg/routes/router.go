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
	sessionHandler *handlers.SessionHandler,
	sessionActivityHandler *handlers.SessionActivityHandler) {
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

	// Serve Swagger documentation
	e.Static("/docs", "swagger-ui")
	e.File("/swagger.json", "swagger.json")

	// Words routes
	e.GET("/api/words", wordHandler.GetWords)
	e.GET("/api/words/random", wordHandler.GetRandomWordFiltered)
	e.GET("/api/words/search", wordHandler.SearchWordsTerm)
	e.GET("/api/words/groups/:group-id", wordHandler.GetWordsByGroup)

	// Groups routes
	e.GET("/api/groups", groupHandler.GetGroups)

	// Study Activities routes
	e.GET("/api/study-activities", studyActivityHandler.GetStudyActivities)

	// Session routes
	e.POST("/api/sessions", sessionHandler.CreateSession)
	e.PUT("/api/sessions", sessionHandler.UpdateSession)
	e.GET("/api/sessions", sessionHandler.GetSessions)
	e.GET("/api/sessions/:id", sessionHandler.GetSessionByID)
	e.DELETE("/api/sessions", func(c echo.Context) error {
		status, err := sessionHandler.DeleteAllSessions(c)
		if err != nil {
			return err
		}
		return c.NoContent(status)
	})

	// Session Activity routes
	e.POST("/api/session-activity", sessionActivityHandler.AddSessionActivity)
}

// SetupSessionRoutes sets up routes for session-related endpoints
func SetupSessionRoutes(e *echo.Echo, sessionHandler *handlers.SessionHandler) {
    // Sessions routes
    e.POST("/api/sessions", sessionHandler.CreateSession)
    e.GET("/api/sessions", sessionHandler.GetSessions)
    e.PUT("/api/sessions", sessionHandler.UpdateSession)
    e.GET("/api/sessions/:id", sessionHandler.GetSessionByID)
    e.DELETE("/api/sessions", func(c echo.Context) error {
        status, err := sessionHandler.DeleteAllSessions(c)
        if err != nil {
            return err
        }
        return c.NoContent(status)
    })
}
