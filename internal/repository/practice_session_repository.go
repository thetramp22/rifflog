package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/thetramp22/rifflog/internal/models"
)

var ErrSkillNotFound = errors.New("skill not found")
var ErrUserNotFound = errors.New("user not found")
var ErrPracticeSessionNotFound = errors.New("practice session not found")

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
			id,
			skill_id, 
			duration_minutes, 
			notes, 
			practiced_at, 
			created_at,
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
		&session.ID,
		&session.SkillID,
		&session.DurationMinutes,
		&session.Notes,
		&session.PracticedAt,
		&session.CreatedAt,
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

func (r *PracticeSessionRepository) UpdatePracticeSession(ctx context.Context, userID int, practiceSessionID int, practiceSession models.PracticeSession) (models.PracticeSession, error) {
	var session models.PracticeSession

	query := `
		UPDATE practice_sessions
		SET 
			skill_id = $1,
			duration_minutes = $2,
			practiced_at = $3,
			notes = $4
		WHERE id = $5
		AND user_id = $6
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
		practiceSession.PracticedAt,
		practiceSession.Notes,
		practiceSessionID,
		userID,
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
	if errors.Is(err, pgx.ErrNoRows) {
		return models.PracticeSession{}, ErrPracticeSessionNotFound
	}
	if err != nil {
		return models.PracticeSession{}, err
	}

	return session, nil
}

func (r *PracticeSessionRepository) DeletePracticeSession(ctx context.Context, userID int, practiceSessionID int) (int, error) {
	var sessionID int

	query := `
		DELETE
		FROM practice_sessions
		WHERE
			id = $1
		AND
			user_id = $2
		RETURNING id
	`

	err := r.DB.QueryRow(
		ctx,
		query,
		practiceSessionID,
		userID,
	).Scan(
		&sessionID,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return 0, ErrPracticeSessionNotFound
	}

	return sessionID, nil
}

func (r *PracticeSessionRepository) GetPracticeSessions(ctx context.Context, userID int, params models.FilterParams) ([]models.PracticeSessionDetails, error) {
	baseQuery := `
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
	`
	var conditions []string
	var args []interface{}
	argCount := 2

	if params.Skill != 0 {
		conditions = append(conditions, fmt.Sprintf("practice_sessions.skill_id = $%d", argCount))
		args = append(args, params.Skill)
		argCount++
	}
	if params.From != nil {
		conditions = append(conditions, fmt.Sprintf("practice_sessions.practiced_at >= $%d", argCount))
		args = append(args, params.From)
		argCount++
	}
	if params.To != nil {
		conditions = append(conditions, fmt.Sprintf("practice_sessions.practiced_at < $%d", argCount))
		args = append(args, params.To)
		argCount++
	}

	finalQuery := baseQuery
	if len(conditions) > 0 {
		finalQuery += " AND " + strings.Join(conditions, " AND ")
	}
	finalQuery += " ORDER BY practice_sessions.practiced_at DESC;"

	queryArgs := append([]interface{}{userID}, args...)
	rows, err := r.DB.Query(ctx, finalQuery, queryArgs...)
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

func (r *PracticeSessionRepository) GetPracticeSessionStats(ctx context.Context, userID int) (models.PracticeSessionStats, error) {
	var totalMinutes int
	totalMinutesQuery := `
		SELECT COALESCE(SUM(duration_minutes), 0)
		FROM practice_sessions
		WHERE user_id = $1
	`
	err := r.DB.QueryRow(ctx, totalMinutesQuery, userID).Scan(&totalMinutes)
	if err != nil {
		return models.PracticeSessionStats{}, err
	}

	var totalSessions int
	totalSessionsQuery := `
		SELECT COUNT(*)
		FROM practice_sessions
		WHERE user_id = $1
	`
	err = r.DB.QueryRow(ctx, totalSessionsQuery, userID).Scan(&totalSessions)
	if err != nil {
		return models.PracticeSessionStats{}, err
	}

	var mostPracticedSkill models.MostPracticedSkill
	mostPracticedSkillQuery := `
		SELECT skills.name AS skill_name, COALESCE(SUM(practice_sessions.duration_minutes), 0) AS total_minutes
		FROM practice_sessions
		INNER JOIN skills ON practice_sessions.skill_id = skills.id		
		WHERE user_id = $1
		GROUP BY skills.id, skills.name
		ORDER BY total_minutes DESC
		LIMIT 1
	`
	err = r.DB.QueryRow(ctx, mostPracticedSkillQuery, userID).Scan(&mostPracticedSkill.Name, &mostPracticedSkill.TotalMinutes)
	if errors.Is(err, pgx.ErrNoRows) {
		return models.PracticeSessionStats{
			TotalMinutes:       0,
			TotalSessions:      0,
			MostPracticedSkill: nil,
			LongestSession:     0,
		}, nil
	}
	if err != nil {
		return models.PracticeSessionStats{}, err
	}

	var longestSession int
	longestSessionQuery := `
		SELECT COALESCE(MAX(duration_minutes), 0)
		FROM practice_sessions
		WHERE user_id = $1
	`
	err = r.DB.QueryRow(ctx, longestSessionQuery, userID).Scan(&longestSession)
	if err != nil {
		return models.PracticeSessionStats{}, err
	}

	stats := models.PracticeSessionStats{
		TotalMinutes:       totalMinutes,
		TotalSessions:      totalSessions,
		MostPracticedSkill: &mostPracticedSkill,
		LongestSession:     longestSession,
	}

	return stats, nil
}
