package handlers

import (
	"context"
	"fmt"
	"io"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/thetramp22/rifflog/internal/auth"
	"github.com/thetramp22/rifflog/internal/bootstrap"
	"github.com/thetramp22/rifflog/internal/config"
	"github.com/thetramp22/rifflog/internal/database"
	"github.com/thetramp22/rifflog/internal/middleware"
	"github.com/thetramp22/rifflog/internal/models"
	repository "github.com/thetramp22/rifflog/internal/repositories"
	"github.com/thetramp22/rifflog/internal/services"
	"golang.org/x/crypto/bcrypt"
)

type TestApp struct {
	Router     *gin.Engine
	DB         *pgxpool.Pool
	UserRepo   *repository.UserRepository
	JWTService *auth.JWTService
}

type TestUser struct {
	User  models.User
	Token string
}

func SetupTestApp(t *testing.T) *TestApp {
	err := godotenv.Load("../../.env.test")
	if err != nil {
		t.Log("No .env file found")
	}
	db := database.NewConnection()
	router := gin.Default()
	jwtService := auth.NewJWTService(config.JWTSecret())
	authMiddleware := middleware.NewAuthMiddleware(jwtService)

	userRepo := repository.NewUserRepository(db)
	userService := services.NewUserService(userRepo, jwtService)
	userHandler := NewUserHandler(userService)

	skillRepo := repository.NewSkillRepository(db)
	skillService := services.NewSkillService(skillRepo)
	skillHandler := NewSkillHandler(skillService)

	practiceSessionRepo := repository.NewPracticeSessionRepository(db)
	practiceSessionService := services.NewPracticeSessionService(practiceSessionRepo)
	practiceSessionHandler := NewPracticeSessionHandler(practiceSessionService)

	gin.SetMode(gin.TestMode)
	gin.DefaultWriter = io.Discard
	gin.DefaultErrorWriter = io.Discard

	ctx := context.Background()
	bootstrap.PopulateSkillsList(ctx, skillRepo)

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

	return &TestApp{
		Router:     router,
		DB:         db,
		UserRepo:   userRepo,
		JWTService: jwtService,
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

func SetupTestUser(t *testing.T, password string) (*TestApp, TestUser) {
	t.Helper()

	app := SetupTestApp(t)

	email := fmt.Sprintf("test-%d@test.com", time.Now().UnixNano())

	user, err := CreateTestUser(app.UserRepo, email, password)
	if err != nil {
		t.Fatalf("failed to register user: %v", err)
	}
	t.Logf("registered user id=%d", user.ID)

	token, err := app.JWTService.GenerateToken(user.ID)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	return app, TestUser{
		User:  user,
		Token: token,
	}
}

func SetupExtraTestUser(t *testing.T, app *TestApp, password string) TestUser {
	t.Helper()

	email := fmt.Sprintf("test-%d@test.com", time.Now().UnixNano())

	user, err := CreateTestUser(app.UserRepo, email, password)
	if err != nil {
		t.Fatalf("failed to register user: %v", err)
	}
	t.Logf("registered user id=%d", user.ID)

	token, err := app.JWTService.GenerateToken(user.ID)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	return TestUser{
		User:  user,
		Token: token,
	}
}
