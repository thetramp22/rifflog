package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/thetramp22/rifflog/internal/models"
)

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

	return err
}
