package services

import (
	"context"

	"github.com/thetramp22/rifflog/internal/models"
	"github.com/thetramp22/rifflog/internal/repository"
)

type SkillService struct {
	Repo *repository.SkillRepository
}

func NewSkillService(repo *repository.SkillRepository) *SkillService {
	return &SkillService{Repo: repo}
}

func (s *SkillService) GetSkills(ctx context.Context) ([]models.Skill, error) {
	skills, err := s.Repo.GetSkills(ctx)
	if err != nil {
		return nil, err
	}

	return skills, nil
}
