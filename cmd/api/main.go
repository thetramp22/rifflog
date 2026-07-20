package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/thetramp22/rifflog/internal/auth"
	"github.com/thetramp22/rifflog/internal/bootstrap"
	"github.com/thetramp22/rifflog/internal/config"
	"github.com/thetramp22/rifflog/internal/database"
	"github.com/thetramp22/rifflog/internal/handlers"
	"github.com/thetramp22/rifflog/internal/middleware"
	repository "github.com/thetramp22/rifflog/internal/repositories"
	"github.com/thetramp22/rifflog/internal/services"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	dbPool := database.NewConnection()
	defer dbPool.Close()

	log.Println("Connected to PostgreSQL")

	router := gin.Default()

	jwtService := auth.NewJWTService(config.JWTSecret())

	authMiddleware := middleware.NewAuthMiddleware(jwtService)

	userRepo := repository.NewUserRepository(dbPool)
	userService := services.NewUserService(userRepo, jwtService)
	userHandler := handlers.NewUserHandler(userService)

	skillRepo := repository.NewSkillRepository(dbPool)
	skillService := services.NewSkillService(skillRepo)
	skillHandler := handlers.NewSkillHandler(skillService)

	practiceSessionRepo := repository.NewPracticeSessionRepository(dbPool)
	practiceSessionService := services.NewPracticeSessionService(practiceSessionRepo)
	practiceSessionHandler := handlers.NewPracticeSessionHandler(practiceSessionService)

	ctx := context.Background()
	bootstrap.PopulateSkillsList(ctx, skillRepo)

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})
	router.POST("/register", userHandler.Register)
	router.POST("/login", userHandler.Login)
	router.GET("/skills", skillHandler.ListSkills)

	protected := router.Group("/api")
	protected.Use(authMiddleware.Authenticate)
	{
		protected.POST("/practice-sessions", practiceSessionHandler.CreatePracticeSession)
		protected.GET("/practice-sessions", practiceSessionHandler.ListPracticeSessions)
		protected.GET("/practice-sessions/stats", practiceSessionHandler.ListPracticeSessionStats)
		protected.PUT("/practice-sessions/:id", practiceSessionHandler.UpdatePracticeSession)
		protected.DELETE("/practice-sessions/:id", practiceSessionHandler.DeletePracticeSession)
	}

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	router.Run(":" + port)
}
