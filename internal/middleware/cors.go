package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/thetramp22/rifflog/internal/config"
)

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		allowedOrigin := config.CORSAllowedOrigin()

		origin := c.GetHeader("Origin")

		if origin == allowedOrigin {
			c.Writer.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
		}

		c.Next()
	}
}
