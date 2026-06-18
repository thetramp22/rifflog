package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
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

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	var user models.User

	query := `
		SELECT
			id,
			email,
			password_hash,
			created_at
		FROM users
		WHERE email = $1
	`
	err := r.DB.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return models.User{}, ErrUserNotFound
	}
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}
