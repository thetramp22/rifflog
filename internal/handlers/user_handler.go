package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thetramp22/rifflog/internal/models"
	"github.com/thetramp22/rifflog/internal/services"
)

// UserHandler handles requests to endpoints dealing with user registration and login.
type UserHandler struct {
	Service *services.UserService
}

// NewUserHandler returns a UserHandler.
func NewUserHandler(service *services.UserService) *UserHandler {
	return &UserHandler{Service: service}
}

// Register recieves a request to register a new user and calls the RegisterUser service.
func (h *UserHandler) Register(c *gin.Context) {
	var req models.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
		})
		return
	}

	user, err := h.Service.RegisterUser(c, req)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "could not create user",
		})
		return
	}

	c.JSON(http.StatusCreated, user)
}

// Login recieves a request to log in and authenticate a registered user, and calls
// the Login service.
func (h *UserHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request",
		})
		return
	}

	user, err := h.Service.Login(c, req)

	if errors.Is(err, services.ErrInvalidPassword) ||
		errors.Is(err, services.ErrUserNotFound) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "could not login",
		})
		return
	}

	c.JSON(http.StatusOK, user)
}
