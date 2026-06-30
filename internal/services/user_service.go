package services

import (
	"context"
	"errors"

	"github.com/thetramp22/rifflog/internal/auth"
	"github.com/thetramp22/rifflog/internal/models"
	"github.com/thetramp22/rifflog/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidPassword = errors.New("Invalid password")

type UserService struct {
	Repo *repository.UserRepository
	JWT  *auth.JWTService
}

func NewUserService(repo *repository.UserRepository, jwt *auth.JWTService) *UserService {
	return &UserService{Repo: repo, JWT: jwt}
}

func (s *UserService) RegisterUser(ctx context.Context, req models.RegisterRequest) (models.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(req.Password),
		bcrypt.DefaultCost,
	)

	if err != nil {
		return models.User{}, err
	}

	user := models.User{
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
	}

	return s.Repo.CreateUser(ctx, user)
}

func (s *UserService) Login(ctx context.Context, req models.LoginRequest) (models.User, error) {
	user, err := s.Repo.GetUserByEmail(ctx, req.Email)
	if errors.Is(err, repository.ErrUserNotFound) {
		return models.User{}, ErrUserNotFound
	}
	if err != nil {
		return models.User{}, err
	}

	if !CheckPasswordHash(req.Password, user.PasswordHash) {
		return models.User{}, ErrInvalidPassword
	}

	return user, nil

}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
