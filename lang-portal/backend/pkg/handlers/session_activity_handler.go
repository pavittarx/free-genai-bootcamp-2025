package handlers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/pavittarx/lang-portal/backend/pkg/services"
)

// SessionActivityHandler handles HTTP requests related to session activities
type SessionActivityHandler struct {
	service *services.SessionActivityService
}

// NewSessionActivityHandler creates a new instance of SessionActivityHandler
func NewSessionActivityHandler(service *services.SessionActivityService) *SessionActivityHandler {
	return &SessionActivityHandler{service: service}
}

// AddSessionActivity handles adding a new activity to a session
func (h *SessionActivityHandler) AddSessionActivity(c echo.Context) error {
	// Request body struct for adding a session activity
	type AddSessionActivityRequest struct {
		SessionID   int64  `json:"session_id"`
		ActivityID  int64  `json:"activity_id"`
		Answer      string `json:"answer"`
	}

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
func (h *SessionActivityHandler) EvaluateSessionActivity(c echo.Context) error {
	// Request body struct for evaluating a session activity
	type EvaluateSessionActivityRequest struct {
		Result string `json:"result"`
		Score  int    `json:"score"`
	}

	// Parse session activity ID from URL parameter
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid session activity ID",
		})
	}

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
