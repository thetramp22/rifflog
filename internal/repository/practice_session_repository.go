package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/thetramp22/rifflog/internal/models"
)

var ErrSkillNotFound = errors.New("Skill not found")

type PracticeSessionRepository struct {
	DB *pgx.Conn
}

func NewPracticeSessionRepository(db *pgx.Conn) *PracticeSessionRepository {
	return &PracticeSessionRepository{DB: db}
}

func (r *PracticeSessionRepository) CreatePracticeSession(ctx context.Context, practiceSession models.PracticeSession) error {
	query := `
		INSERT INTO practice_sessions (skill_id, duration_minutes, notes, practiced_at, user_id)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.DB.Exec(
		ctx,
		query,
		practiceSession.SkillID,
		practiceSession.DurationMinutes,
		practiceSession.Notes,
		practiceSession.PracticedAt,
		practiceSession.UserID,
	)

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == "23503" {
			return ErrSkillNotFound
		}
	}

	return err
}
