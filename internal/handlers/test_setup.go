package handlers

import (
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"github.com/thetramp22/rifflog/internal/bootstrap"
	"github.com/thetramp22/rifflog/internal/database"
	"github.com/thetramp22/rifflog/internal/repository"
	"github.com/thetramp22/rifflog/internal/services"
)

type TestApp struct {
	Router *gin.Engine
	DB     *pgx.Conn
}

func SetupTestApp(t *testing.T) *TestApp {
	t.Log("starting setup")
	err := godotenv.Load("../../.env.test")
	if err != nil {
		t.Log("No .env file found")
	}

	t.Log(os.Getwd())

	t.Log("connecting to database")
	db := database.NewConnection()

	router := gin.Default()

	userRepo := repository.NewUserRepository(db)
	userService := services.NewUserService(userRepo)
	userHandler := NewUserHandler(userService)

	skillRepo := repository.NewSkillRepository(db)
	skillService := services.NewSkillService(skillRepo)
	skillHandler := NewSkillHandler(skillService)

	practiceSessionRepo := repository.NewPracticeSessionRepository(db)
	practiceSessionService := services.NewPracticeSessionService(practiceSessionRepo)
	practiceSessionHandler := NewPracticeSessionHandler(practiceSessionService)

	t.Log("seeding skills")
	bootstrap.PopulateSkillsList(skillRepo)

	router.POST("/register", userHandler.Register)
	router.GET("/skills", skillHandler.ListSkills)
	router.POST("/practice-sessions", practiceSessionHandler.CreatePracticeSession)

	return &TestApp{
		Router: router,
		DB:     db,
	}
}
