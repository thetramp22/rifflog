package config

import (
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

// DatabaseURL return the URL to the database built from env variables.
func DatabaseURL() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)
}

// ParseConfig returns a configuration for a pgx connection pool.
func ParseConfig() (*pgxpool.Config, error) {
	return pgxpool.ParseConfig(DatabaseURL())
}
