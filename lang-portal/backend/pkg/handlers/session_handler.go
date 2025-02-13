package handlers

import (
	"net/http"
	"strconv"

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

// CreateSession handles the creation of a new session
func (h *SessionHandler) CreateSession(c echo.Context) error {
	// Request body struct for session creation
	type CreateSessionRequest struct {
		ActivityID int64  `json:"activity_id"`
		GroupID    *int64 `json:"group_id,omitempty"`
	}

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

	// Create session
	session, err := h.service.CreateSession(c.Request().Context(), req.ActivityID, req.GroupID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to create session",
		})
	}

	return c.JSON(http.StatusCreated, session)
}

// GetSessionByID retrieves a specific session
func (h *SessionHandler) GetSessionByID(c echo.Context) error {
	// Parse session ID from URL parameter
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid session ID",
		})
	}

	// Retrieve session
	session, err := h.service.GetSessionByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "Session not found",
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
	// Request body struct for session update
	type UpdateSessionRequest struct {
		Score int `json:"score"`
	}

	// Parse session ID from URL parameter
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid session ID",
		})
	}

	var req UpdateSessionRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	// End the session
	if err := h.service.EndSession(c.Request().Context(), id, req.Score); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to update session",
		})
	}

	return c.NoContent(http.StatusNoContent)
}

// DeleteSession handles session deletion
func (h *SessionHandler) DeleteSession(c echo.Context) error {
	// Parse session ID from URL parameter
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid session ID",
		})
	}

	// Delete the session
	if err := h.service.DeleteSession(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to delete session",
		})
	}

	return c.NoContent(http.StatusNoContent)
}
