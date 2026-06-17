package handlers

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/thetramp22/rifflog/internal/bootstrap"
	"github.com/thetramp22/rifflog/internal/database"
	"github.com/thetramp22/rifflog/internal/models"
	"github.com/thetramp22/rifflog/internal/repository"
	"github.com/thetramp22/rifflog/internal/services"
	"golang.org/x/crypto/bcrypt"
)

type TestApp struct {
	Router   *gin.Engine
	DB       *pgxpool.Pool
	UserRepo *repository.UserRepository
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
	ctx := context.Background()
	bootstrap.PopulateSkillsList(ctx, skillRepo)

	router.POST("/register", userHandler.Register)
	router.GET("/skills", skillHandler.ListSkills)
	router.POST("/practice-sessions", practiceSessionHandler.CreatePracticeSession)
	router.GET("/practice-sessions", practiceSessionHandler.ListPracticeSessions)

	return &TestApp{
		Router:   router,
		DB:       db,
		UserRepo: userRepo,
	}
}

func CreateTestUser(r *repository.UserRepository, email string, password string) (models.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return models.User{}, err
	}

	var user models.User
	query := `
		INSERT INTO users (email, password_hash)
		VALUES ($1, $2)
		RETURNING email, password_hash, id, created_at
	`

	err = r.DB.QueryRow(
		context.Background(),
		query,
		email,
		hashedPassword,
	).Scan(
		&user.Email,
		&user.PasswordHash,
		&user.ID,
		&user.CreatedAt,
	)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func SetupTestUser(t *testing.T) (*TestApp, models.User) {
	t.Helper()

	app := SetupTestApp(t)

	email := fmt.Sprintf("test-%d@test.com", time.Now().UnixNano())
	password := "1234"

	user, err := CreateTestUser(app.UserRepo, email, password)
	if err != nil {
		t.Fatalf("Failed to register user: %v", err)
	}
	t.Logf("registered user id=%d", user.ID)

	return app, user
}
