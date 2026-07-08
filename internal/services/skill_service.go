package services

import (
	"context"

	"github.com/thetramp22/rifflog/internal/models"
	"github.com/thetramp22/rifflog/internal/repository"
)

// SkillService provides methods dealing with skill.
type SkillService struct {
	Repo *repository.SkillRepository
}

// NewSkillService returns a SkillService.
func NewSkillService(repo *repository.SkillRepository) *SkillService {
	return &SkillService{Repo: repo}
}

// GetSkills is called from the handler and calls the repository method.
// Returns a list of skills.
func (s *SkillService) GetSkills(ctx context.Context) ([]models.Skill, error) {
	skills, err := s.Repo.GetSkills(ctx)
	if err != nil {
		return nil, err
	}

	return skills, nil
}
