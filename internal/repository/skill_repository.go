package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/thetramp22/rifflog/internal/models"
)

type SkillRepository struct {
	DB *pgxpool.Pool
}

func NewSkillRepository(db *pgxpool.Pool) *SkillRepository {
	return &SkillRepository{DB: db}
}

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
