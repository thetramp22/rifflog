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

func (r *UserRepository) CreateUser(ctx context.Context, user models.User) (models.User, error) {
	var createdUser models.User

	query := `
		INSERT INTO users (email, password_hash)
		VALUES ($1, $2)
		RETURNING
			id,
			email,
			password_hash,
			created_at
	`

	err := r.DB.QueryRow(
		ctx,
		query,
		user.Email,
		user.PasswordHash,
	).Scan(
		&createdUser.ID,
		&createdUser.Email,
		&createdUser.PasswordHash,
		&createdUser.CreatedAt,
	)
	if err != nil {
		return models.User{}, err
	}

	return createdUser, err
}

func (r *UserRepository) GetUserByEmail()
