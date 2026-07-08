package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thetramp22/rifflog/internal/middleware"
	"github.com/thetramp22/rifflog/internal/models"
	"github.com/thetramp22/rifflog/internal/services"
)

// PracticeSessionHandler handles requests to endpoints dealing with practice sessions.
type PracticeSessionHandler struct {
	Service *services.PracticeSessionService
}

// NewPracticeSessionHandler returns a PracticeSessionHandler.
func NewPracticeSessionHandler(service *services.PracticeSessionService) *PracticeSessionHandler {
	return &PracticeSessionHandler{Service: service}
}

// CreatePracticeSession recieves a request to create a practice session and calls
// the CreatePracticeSession service.
func (h *PracticeSessionHandler) CreatePracticeSession(c *gin.Context) {
	var req models.CreatePracticeSessionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
		})
		return
	}

	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "invalid user",
		})
		return
	}

	practiceSession, err := h.Service.CreatePracticeSession(c, userID, req)

	if errors.Is(err, services.ErrInvalidDuration) ||
		errors.Is(err, services.ErrInvalidSkillID) ||
		errors.Is(err, services.ErrInvalidUserID) ||
		errors.Is(err, services.ErrInvalidPracticedAt) ||
		errors.Is(err, services.ErrUserNotFound) ||
		errors.Is(err, services.ErrSkillNotFound) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "could not create practice session",
		})
		return
	}

	c.JSON(http.StatusCreated, practiceSession)
}

// UpdatePracticeSession recieves a request to update a practice session and calls
// the UpdatePracticeSession service.
func (h *PracticeSessionHandler) UpdatePracticeSession(c *gin.Context) {
	var req models.UpdatePracticeSessionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
		})
		return
	}

	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	var sessionID models.PracticeSessionURI
	if err := c.ShouldBindUri(&sessionID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid or missing id parameter",
		})
		return
	}

	practiceSession, err := h.Service.UpdatePracticeSession(c, userID, sessionID.ID, req)

	if errors.Is(err, services.ErrInvalidDuration) ||
		errors.Is(err, services.ErrInvalidSkillID) ||
		errors.Is(err, services.ErrInvalidUserID) ||
		errors.Is(err, services.ErrInvalidPracticedAt) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	if errors.Is(err, services.ErrPracticeSessionNotFound) ||
		errors.Is(err, services.ErrUserNotFound) ||
		errors.Is(err, services.ErrSkillNotFound) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "could not update practice session",
		})
		return
	}

	c.JSON(http.StatusOK, practiceSession)
}

// DeletePracticeSession recieves a request to delete a practice session and calls
// the DeletePracticeSession service.
func (h *PracticeSessionHandler) DeletePracticeSession(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	var sessionID models.PracticeSessionURI
	if err := c.ShouldBindUri(&sessionID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid or missing id parameter",
		})
		return
	}

	deletedSessionID, err := h.Service.DeletePracticeSession(c, userID, sessionID.ID)
	if errors.Is(err, services.ErrPracticeSessionNotFound) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "could not update practice session",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("practice session %v deleted", deletedSessionID),
	})
}

// ListPracticeSessions recieves a request to list practice sessions for the
// current user and calls the GetPracticeSessions service.
func (h *PracticeSessionHandler) ListPracticeSessions(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "invalid user",
		})
		return
	}

	var params models.FilterParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid query parameters",
		})
		return
	}
	if params.To != nil {
		to := params.To.AddDate(0, 0, 1)
		params.To = &to
	}
	if userID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request",
		})
		return
	}

	practiceSessionDetails, err := h.Service.GetPracticeSessions(c, userID, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Could not get list of practice sessions",
		})
		return
	}

	c.JSON(http.StatusOK, practiceSessionDetails)
}

// ListPracticeSessionStats recieves a request to list practice session stats for the
// current user and calls the GetPracticeSessionStats service.
func (h *PracticeSessionHandler) ListPracticeSessionStats(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "invalid user",
		})
		return
	}

	practiceSessionStats, err := h.Service.GetPracticeSessionStats(c, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Could not get practice session stats",
		})
		return
	}

	c.JSON(http.StatusOK, practiceSessionStats)
}
