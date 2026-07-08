package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thetramp22/rifflog/internal/services"
)

// SkillHandler handles requests to endpoints dealing with skills.
type SkillHandler struct {
	Service *services.SkillService
}

// NewSkillHandler returns a SkillHandler
func NewSkillHandler(service *services.SkillService) *SkillHandler {
	return &SkillHandler{Service: service}
}

// ListSkills recieves a request to list skills and calls the GetSkills service.
func (h *SkillHandler) ListSkills(c *gin.Context) {
	skills, err := h.Service.GetSkills(c)
	if err != nil {
		log.Printf("Error getting skills: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "could not get list of skills",
		})
		return
	}

	c.JSON(http.StatusOK, skills)
}
