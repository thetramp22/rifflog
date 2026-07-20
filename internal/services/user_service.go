package services

import (
	"context"
	"errors"

	"github.com/thetramp22/rifflog/internal/auth"
	"github.com/thetramp22/rifflog/internal/models"
	repository "github.com/thetramp22/rifflog/internal/repositories"
	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidPassword = errors.New("Invalid password")

// UserService provides methods dealing with users and authentication.
type UserService struct {
	Repo *repository.UserRepository
	JWT  *auth.JWTService
}

func NewUserService(repo *repository.UserRepository, jwt *auth.JWTService) *UserService {
	return &UserService{Repo: repo, JWT: jwt}
}

// RegisterUser takes a request from a handler, encrypts the users password,
// and calls the repository method to create the user.
// Returns the created user on success.
func (s *UserService) RegisterUser(ctx context.Context, req models.RegisterRequest) (models.UserResponse, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(req.Password),
		bcrypt.DefaultCost,
	)

	if err != nil {
		return models.UserResponse{}, err
	}

	user := models.User{
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
	}

	newUser, err := s.Repo.CreateUser(ctx, user)
	if err != nil {
		return models.UserResponse{}, err
	}

	return models.UserResponse{
		ID:        newUser.ID,
		Email:     newUser.Email,
		CreatedAt: newUser.CreatedAt,
	}, nil
}

// Login takes a request from a handler, retrieves the user information from the repository,
// and authenticates the user.
// Returns a response containing the user information and authentication token.
func (s *UserService) Login(ctx context.Context, req models.LoginRequest) (models.LoginResponse, error) {
	user, err := s.Repo.GetUserByEmail(ctx, req.Email)
	if errors.Is(err, repository.ErrUserNotFound) {
		return models.LoginResponse{}, ErrUserNotFound
	}
	if err != nil {
		return models.LoginResponse{}, err
	}

	if !checkPasswordHash(req.Password, user.PasswordHash) {
		return models.LoginResponse{}, ErrInvalidPassword
	}

	token, err := s.JWT.GenerateToken(user.ID)
	if err != nil {
		return models.LoginResponse{}, err
	}

	loginResponse := models.LoginResponse{
		Token: token,
		User: models.UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
		},
	}

	return loginResponse, nil

}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
