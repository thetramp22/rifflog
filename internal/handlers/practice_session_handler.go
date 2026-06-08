package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thetramp22/rifflog/internal/models"
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

	err := h.Service.CreatePracticeSession(c, req)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "could not create practice session",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "practice session created",
	})
}
