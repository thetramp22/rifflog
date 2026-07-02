package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thetramp22/rifflog/internal/middleware"
	"github.com/thetramp22/rifflog/internal/models"
	"github.com/thetramp22/rifflog/internal/repository"
	"github.com/thetramp22/rifflog/internal/services"
)

type PracticeSessionHandler struct {
	Service *services.PracticeSessionService
}

func NewPracticeSessionHandler(service *services.PracticeSessionService) *PracticeSessionHandler {
	return &PracticeSessionHandler{Service: service}
}

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
		errors.Is(err, repository.ErrSkillNotFound) {
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
