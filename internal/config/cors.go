package config

import (
	"os"
	"strings"
)

func CORSAllowedOrigins() []string {
	originsStr := os.Getenv("CORS_ALLOWED_ORIGINS")
	origins := strings.Split(originsStr, ",")
	for i, item := range origins {
		origins[i] = strings.TrimSpace(item)
	}
	return origins
}
