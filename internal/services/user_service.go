package services

import (
	"context"

	"github.com/thetramp22/rifflog/internal/models"
	"github.com/thetramp22/rifflog/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	Repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{Repo: repo}
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
