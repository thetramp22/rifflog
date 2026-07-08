package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/thetramp22/rifflog/internal/models"
)

// SkillRepository provides methods to access and manipulate the application database.
type SkillRepository struct {
	DB *pgxpool.Pool
}

// NewSkillRepository returns a SkillRepository.
func NewSkillRepository(db *pgxpool.Pool) *SkillRepository {
	return &SkillRepository{DB: db}
}

// CreateSkill stores a skill in the database.
func (r *SkillRepository) CreateSkill(ctx context.Context, skill models.Skill) error {
	query := `
		INSERT INTO skills (name, description)
		VALUES ($1, $2)
	`

	_, err := r.DB.Exec(
		ctx,
		query,
		skill.Name,
		skill.Description,
	)

	return err
}

// SeedSkill stores a skill in the database. Used by the bootstrap package to populate
// skills on application startup.
func (r *SkillRepository) SeedSkill(ctx context.Context, skill models.Skill) error {
	query := `
		INSERT INTO skills (name, description)
		VALUES ($1, $2)
		ON CONFLICT (name) DO NOTHING
	`

	_, err := r.DB.Exec(
		ctx,
		query,
		skill.Name,
		skill.Description,
	)

	return err
}

// GetSkills retrieves all skills stored in the database and returns them in a list.
func (r *SkillRepository) GetSkills(ctx context.Context) ([]models.Skill, error) {
	query := `
		SELECT id, name, description, created_at
		FROM skills
		ORDER BY name
	`

	rows, err := r.DB.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	skills, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.Skill])
	if err != nil {
		return nil, err
	}

	return skills, nil
}
