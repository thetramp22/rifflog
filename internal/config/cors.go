package config

import "os"

func CORSAllowedOrigin() string {
	return os.Getenv("CORS_ALLOWED_ORIGIN")
}
