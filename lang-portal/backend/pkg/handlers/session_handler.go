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

// CreateSessionRequest defines the request payload for creating a session
type CreateSessionRequest struct {
	ActivityID int64  `json:"activity_id" validate:"required"`
	GroupID    *int64 `json:"group_id,omitempty"`
}

// UpdateSessionRequest defines the request payload for updating a session
type UpdateSessionRequest struct {
	Score int `json:"score" validate:"min=0,max=100"`
}

// CreateSession handles the creation of a new session
// @Summary Create a new session
// @Description Create a new learning session with an activity and optional group
// @Tags sessions
// @Accept json
// @Produce json
// @Param session body CreateSessionRequest true "Session creation details"
// @Success 201 {object} models.Session "Session created successfully"
// @Failure 400 {object} map[string]string "Invalid request payload"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/sessions [post]
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
// @Summary Get a session by ID
// @Description Retrieve a specific session's details using its unique identifier
// @Tags sessions
// @Produce json
// @Param id path int true "Session ID"
// @Success 200 {object} models.Session "Session retrieved successfully"
// @Failure 400 {object} map[string]string "Invalid session ID"
// @Failure 404 {object} map[string]string "Session not found"
// @Router /api/v1/sessions/{id} [get]
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
// @Summary List sessions
// @Description Retrieve a paginated list of learning sessions
// @Tags sessions
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Success 200 {array} models.Session "Sessions retrieved successfully"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/sessions [get]
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
// @Summary Update a session
// @Description End a session and record its score
// @Tags sessions
// @Accept json
// @Produce json
// @Param id path int true "Session ID"
// @Param session body UpdateSessionRequest true "Session update details"
// @Success 204 "Session updated successfully"
// @Failure 400 {object} map[string]string "Invalid request payload or session ID"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/sessions/{id} [put]
func (h *SessionHandler) UpdateSession(c echo.Context) error {
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
// @Summary Delete a session
// @Description Delete an existing session by its unique identifier
// @Tags sessions
// @Produce json
// @Param id path int true "Session ID"
// @Success 204 "Session deleted successfully"
// @Failure 400 {object} map[string]string "Invalid session ID"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/sessions/{id} [delete]
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
