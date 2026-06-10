package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
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

	practiceSession, err := h.Service.CreatePracticeSession(c, req)

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
