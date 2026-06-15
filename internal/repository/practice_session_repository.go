package repository

import (
	"context"
	"errors"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/thetramp22/rifflog/internal/models"
)

var ErrSkillNotFound = errors.New("Skill not found")
var ErrUserNotFound = errors.New("User not found")

type PracticeSessionRepository struct {
	DB *pgxpool.Pool
}

func NewPracticeSessionRepository(db *pgxpool.Pool) *PracticeSessionRepository {
	return &PracticeSessionRepository{DB: db}
}

func (r *PracticeSessionRepository) CreatePracticeSession(ctx context.Context, practiceSession models.PracticeSession) (models.PracticeSession, error) {
	var session models.PracticeSession

	query := `
		INSERT INTO practice_sessions (skill_id, duration_minutes, notes, practiced_at, user_id)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING 
			skill_id, 
			duration_minutes, 
			notes, 
			practiced_at, 
			user_id
	`

	err := r.DB.QueryRow(
		ctx,
		query,
		practiceSession.SkillID,
		practiceSession.DurationMinutes,
		practiceSession.Notes,
		practiceSession.PracticedAt,
		practiceSession.UserID,
	).Scan(
		&session.SkillID,
		&session.DurationMinutes,
		&session.Notes,
		&session.PracticedAt,
		&session.UserID,
	)
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == "23503" {
			if strings.Contains(pgErr.Detail, "Key (skill_id)=") {
				return models.PracticeSession{}, ErrSkillNotFound
			}
			if strings.Contains(pgErr.Detail, "Key (user_id)=") {
				return models.PracticeSession{}, ErrUserNotFound
			}
		}
	}
	if err != nil {
		return models.PracticeSession{}, err
	}

	return session, nil
}

func (r *PracticeSessionRepository) GetPracticeSessions(ctx context.Context, userID int) ([]models.PracticeSessionDetails, error) {
	query := `
		SELECT
			practice_sessions.id AS session_id,
			practice_sessions.skill_id AS skill_id,
			skills.name AS skill_name,
			skills.description AS skill_description,
			practice_sessions.duration_minutes AS duration_minutes,
			practice_sessions.notes AS notes,
			practice_sessions.practiced_at AS practiced_at,
			practice_sessions.created_at AS created_at,
			practice_sessions.user_id AS user_id
		FROM practice_sessions
		INNER JOIN skills ON practice_sessions.skill_id = skills.id
		WHERE practice_sessions.user_id = $1
		ORDER BY practice_sessions.practiced_at DESC;
	`

	rows, err := r.DB.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	practiceSessions, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.PracticeSessionDetails])
	if err != nil {
		return nil, err
	}

	return practiceSessions, nil
}
