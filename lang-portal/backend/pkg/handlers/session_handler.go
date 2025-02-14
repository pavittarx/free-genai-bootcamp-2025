package handlers

import (
	"net/http"
	"strconv"
	"log"
	"strings"
	"database/sql"
	"errors"
	"context"

	"github.com/labstack/echo/v4"
	"github.com/pavittarx/lang-portal/backend/pkg/services"
)

// SessionHandler handles HTTP requests related to sessions
type SessionHandler struct {
	service *services.SessionService
}

// NewSessionHandler creates a new instance of SessionHandler
func NewSessionHandler(service *services.SessionService) *SessionHandler {
	return &SessionHandler{service: service}
}

// CreateSessionRequest defines the request payload for creating a session
type CreateSessionRequest struct {
	ActivityID int64 `json:"activity_id" validate:"required"`
}

// UpdateSessionRequest defines the request payload for updating a session
type UpdateSessionRequest struct {
	SessionID int64 `json:"session_id" validate:"required,min=1"`
	Score     int   `json:"score" validate:"min=0,max=100"`
}

// CreateSession handles the creation of a new session

func (h *SessionHandler) CreateSession(c echo.Context) error {
	var req CreateSessionRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	// Validate activity ID
	if req.ActivityID <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid activity ID",
		})
	}

	// Create session with automatic start_time
	session, err := h.service.CreateSession(c.Request().Context(), req.ActivityID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to create session",
		})
	}

	return c.JSON(http.StatusCreated, session)
}

// GetSessionByID retrieves a specific session with its activities

func (h *SessionHandler) GetSessionByID(c echo.Context) error {
	// Validate input parameters
	idStr := c.Param("id")
	if idStr == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Session ID is required",
		})
	}

	// Parse session ID from URL parameter
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid session ID: must be a positive integer",
		})
	}

	// Retrieve session with activities
	session, err := h.service.GetSessionByIDWithActivities(c.Request().Context(), id)
	if err != nil {
		// Log the error for server-side tracking
		log.Printf("Error retrieving session %d: %v", id, err)

		// Check for specific error types
		if strings.Contains(err.Error(), "no session available") {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": err.Error(),
			})
		}

		// Generic server error for unexpected issues
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Unable to retrieve session",
		})
	}

	return c.JSON(http.StatusOK, session)
}

// GetSessions retrieves a list of sessions with pagination

func (h *SessionHandler) GetSessions(c echo.Context) error {
	// Parse pagination parameters
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(c.QueryParam("page_size"))
	if err != nil || pageSize < 1 {
		pageSize = 10 // Default page size
	}

	// Retrieve sessions
	sessions, err := h.service.ListSessions(c.Request().Context(), page, pageSize)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to retrieve sessions",
		})
	}

	return c.JSON(http.StatusOK, sessions)
}

// UpdateSession handles updating a session (e.g., ending a session)

func (h *SessionHandler) UpdateSession(c echo.Context) error {
	var req UpdateSessionRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	// End the session
	if err := h.service.EndSession(c.Request().Context(), req.SessionID, req.Score); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to update session",
		})
	}

	return c.NoContent(http.StatusNoContent)
}

// DeleteAllSessions handles deletion of all sessions
func (h *SessionHandler) DeleteAllSessions(c echo.Context) (int, error) {
	// Delete all sessions and their associated session activities
	sessionsDeleted, err := h.service.DeleteAllSessions(c.Request().Context())
	if err != nil {
		// Log the error for server-side tracking
		log.Printf("Error deleting all sessions: %v", err)

		// Provide more specific error responses
		switch {
		case errors.Is(err, sql.ErrConnDone):
			return http.StatusServiceUnavailable, c.JSON(http.StatusServiceUnavailable, map[string]string{
				"error": "Database connection is closed",
			})
		case errors.Is(err, context.DeadlineExceeded):
			return http.StatusRequestTimeout, c.JSON(http.StatusRequestTimeout, map[string]string{
				"error": "Operation timed out",
			})
		default:
			return http.StatusInternalServerError, c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "Failed to delete all sessions",
			})
		}
	}

	// If no sessions were deleted, return a specific response
	if sessionsDeleted == 0 {
		return http.StatusOK, c.JSON(http.StatusOK, map[string]string{
			"message": "No sessions found to delete",
		})
	}

	return http.StatusNoContent, c.NoContent(http.StatusNoContent)
}
