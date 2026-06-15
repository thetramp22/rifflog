package database

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/thetramp22/rifflog/internal/config"
)

func NewConnection() *pgxpool.Pool {
	dbConfig, err := config.ParseConfig()
	if err != nil {
		log.Fatal(err)
	}

	var dbPool *pgxpool.Pool

	for i := 0; i < 10; i++ {
		dbPool, err = pgxpool.NewWithConfig(context.Background(), dbConfig)

		if err == nil {
			log.Println("Connected to PostgreSQL")
			return dbPool
		}

		log.Printf("Database not ready... retrying (%d/10)\n", i+1)

		time.Sleep(2 * time.Second)
	}

	log.Fatal("Unable to connect to database:", err)

	return nil
}
