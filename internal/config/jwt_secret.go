package config

import "os"

// JWTSecret returns the jwt secret env variable.
func JWTSecret() string {
	return os.Getenv("JWT_SECRET")
}
