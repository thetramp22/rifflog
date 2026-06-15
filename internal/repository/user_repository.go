package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/thetramp22/rifflog/internal/models"
)

type UserRepository struct {
	DB *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) CreateUser(user models.User) error {
	query := `
		INSERT INTO users (email, password_hash)
		VALUES ($1, $2)
	`

	_, err := r.DB.Exec(
		context.Background(),
		query,
		user.Email,
		user.PasswordHash,
	)

	return err
}
