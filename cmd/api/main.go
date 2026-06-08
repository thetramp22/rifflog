package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/thetramp22/rifflog/internal/bootstrap"
	"github.com/thetramp22/rifflog/internal/database"
	"github.com/thetramp22/rifflog/internal/handlers"
	"github.com/thetramp22/rifflog/internal/repository"
	"github.com/thetramp22/rifflog/internal/services"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	conn := database.NewConnection()
	defer conn.Close(context.Background())

	log.Println("Connected to PostgreSQL")

	router := gin.Default()

	userRepo := repository.NewUserRepository(conn)
	userService := services.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)

	skillRepo := repository.NewSkillRepository(conn)
	skillService := services.NewSkillService(skillRepo)
	skillHandler := handlers.NewSkillHandler(skillService)

	practiceSessionRepo := repository.NewPracticeSessionRepository(conn)
	practiceSessionService := services.NewPracticeSessionService(practiceSessionRepo)
	practiceSessionHandler := handlers.NewPracticeSessionHandler(practiceSessionService)

	bootstrap.PopulateSkillsList(skillRepo)

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})
	router.POST("/register", userHandler.Register)
	router.GET("/skills", skillHandler.ListSkills)
	router.POST("/practice-sessions", practiceSessionHandler.CreatePracticeSession)

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	router.Run(":" + port)
}
