package services

import (
	"context"
	"errors"

	"github.com/thetramp22/rifflog/internal/models"
	"github.com/thetramp22/rifflog/internal/repository"
)

var ErrInvalidDuration = errors.New("Duration must be greater than zero")

type PracticeSessionService struct {
	Repo *repository.PracticeSessionRepository
}

func NewPracticeSessionService(repo *repository.PracticeSessionRepository) *PracticeSessionService {
	return &PracticeSessionService{Repo: repo}
}

func (s *PracticeSessionService) CreatePracticeSession(ctx context.Context, req models.CreatePracticeSessionRequest) error {
	err := validateDuration(req.DurationMinutes)
	if err != nil {
		return err
	}
	practiceSession := models.PracticeSession{
		SkillID:         req.SkillID,
		DurationMinutes: req.DurationMinutes,
		PracticedAt:     req.PracticedAt,
		Notes:           req.Notes,
		UserID:          req.UserID,
	}

	return s.Repo.CreatePracticeSession(ctx, practiceSession)
}

func validateDuration(duration int) error {
	if duration > 0 {
		return nil
	}
	return ErrInvalidDuration
}
