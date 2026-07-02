package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/thetramp22/rifflog/internal/auth"
)

const ContextUserID = "userID"

type AuthMiddleware struct {
	JWT *auth.JWTService
}

func NewAuthMiddleware(jwt *auth.JWTService) *AuthMiddleware {
	return &AuthMiddleware{JWT: jwt}
}

func (m *AuthMiddleware) Authenticate(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authorization header is missing"})
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	if tokenString == authHeader {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authorization header is malformed"})
		return
	}

	claims, err := m.JWT.ValidateToken(tokenString)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	c.Set(ContextUserID, claims.UserID)

	c.Next()
}

func GetUserID(c *gin.Context) (int, error) {
	id, ok := c.Get(ContextUserID)
	if !ok {
		return 0, errors.New("user id not found")
	}

	userID, ok := id.(int)
	if !ok {
		return 0, errors.New("invalid user id type")
	}

	return userID, nil
}
