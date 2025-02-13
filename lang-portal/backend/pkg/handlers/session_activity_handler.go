package handlers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/pavittarx/lang-portal/backend/pkg/services"

	// Swagger annotations and request types
	// @Summary Session Activity Handler
	// @Description Handles HTTP requests related to session activities
	// @Tags session-activities
)

// SessionActivityHandler handles HTTP requests related to session activities
type SessionActivityHandler struct {
	service *services.SessionActivityService
}

// NewSessionActivityHandler creates a new instance of SessionActivityHandler
func NewSessionActivityHandler(service *services.SessionActivityService) *SessionActivityHandler {
	return &SessionActivityHandler{service: service}
}

// AddSessionActivityRequest defines the request payload for adding a session activity
type AddSessionActivityRequest struct {
	SessionID   int64  `json:"session_id" validate:"required,min=1"`
	ActivityID  int64  `json:"activity_id" validate:"required,min=1"`
	Answer      string `json:"answer" validate:"required"`
}

// EvaluateSessionActivityRequest defines the request payload for evaluating a session activity
type EvaluateSessionActivityRequest struct {
	Result string `json:"result" validate:"required"`
	Score  int    `json:"score" validate:"min=0,max=100"`
}

// AddSessionActivity handles adding a new activity to a session
// @Summary Add session activity
// @Description Add a new activity to an existing session
// @Tags session-activities
// @Accept json
// @Produce json
// @Param session_id path int true "Session ID"
// @Param activity body AddSessionActivityRequest true "Session activity details"
// @Success 201 {object} models.SessionActivity "Session activity added successfully"
// @Failure 400 {object} map[string]string "Invalid request payload or session/activity ID"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /sessions/{session_id}/activities [post]
func (h *SessionActivityHandler) AddSessionActivity(c echo.Context) error {
	// Request body struct for adding a session activity
	var req AddSessionActivityRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	// Validate input
	if req.SessionID <= 0 || req.ActivityID <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid session or activity ID",
		})
	}

	// Add session activity
	sessionActivity, err := h.service.AddSessionActivity(
		c.Request().Context(), 
		req.SessionID, 
		req.ActivityID, 
		req.Answer,
	)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to add session activity",
		})
	}

	return c.JSON(http.StatusCreated, sessionActivity)
}

// EvaluateSessionActivity handles scoring and updating a session activity
// @Summary Evaluate session activity
// @Description Evaluate and score a session activity
// @Tags session-activities
// @Accept json
// @Produce json
// @Param id path int true "Session Activity ID"
// @Param evaluation body EvaluateSessionActivityRequest true "Evaluation details"
// @Success 204 "Session activity evaluated successfully"
// @Failure 400 {object} map[string]string "Invalid request payload or session activity ID"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /session-activities/{id} [put]
func (h *SessionActivityHandler) EvaluateSessionActivity(c echo.Context) error {
	// Parse session activity ID from URL parameter
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid session activity ID",
		})
	}

	// Request body struct for evaluating a session activity
	var req EvaluateSessionActivityRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	// Evaluate session activity
	if err := h.service.EvaluateSessionActivity(
		c.Request().Context(), 
		id, 
		req.Result, 
		req.Score,
	); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to evaluate session activity",
		})
	}

	return c.NoContent(http.StatusNoContent)
}

// GetSessionActivities retrieves all activities for a specific session
// @Summary List session activities
// @Description Retrieve all activities for a specific session
// @Tags session-activities
// @Produce json
// @Param session_id path int true "Session ID"
// @Success 200 {array} models.SessionActivity "Session activities retrieved successfully"
// @Failure 400 {object} map[string]string "Invalid session ID"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /sessions/{session_id}/activities [get]
func (h *SessionActivityHandler) GetSessionActivities(c echo.Context) error {
	// Parse session ID from URL parameter
	idStr := c.Param("session_id")
	sessionID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid session ID",
		})
	}

	// Retrieve session activities
	sessionActivities, err := h.service.GetSessionActivities(
		c.Request().Context(), 
		sessionID,
	)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to retrieve session activities",
		})
	}

	return c.JSON(http.StatusOK, sessionActivities)
}

// DeleteSessionActivity handles removing a session activity
// @Summary Delete session activity
// @Description Delete a specific session activity
// @Tags session-activities
// @Produce json
// @Param id path int true "Session Activity ID"
// @Success 204 "Session activity deleted successfully"
// @Failure 400 {object} map[string]string "Invalid session activity ID"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /session-activities/{id} [delete]
func (h *SessionActivityHandler) DeleteSessionActivity(c echo.Context) error {
	// Parse session activity ID from URL parameter
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid session activity ID",
		})
	}

	// Delete session activity
	if err := h.service.DeleteSessionActivity(
		c.Request().Context(), 
		id,
	); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to delete session activity",
		})
	}

	return c.NoContent(http.StatusNoContent)
}
