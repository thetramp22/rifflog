package services

import (
	"context"
	"errors"

	"github.com/thetramp22/rifflog/internal/models"
	"github.com/thetramp22/rifflog/internal/repository"
)

var ErrInvalidDuration = errors.New("Duration must be greater than zero")
var ErrInvalidSkillID = errors.New("Invalid skill id")
var ErrInvalidUserID = errors.New("Invalid user id")
var ErrInvalidPracticedAt = errors.New("Invalid practiced at time")

type PracticeSessionService struct {
	Repo *repository.PracticeSessionRepository
}

func NewPracticeSessionService(repo *repository.PracticeSessionRepository) *PracticeSessionService {
	return &PracticeSessionService{Repo: repo}
}

func (s *PracticeSessionService) CreatePracticeSession(ctx context.Context, req models.CreatePracticeSessionRequest) error {
	err := validateRequest(req)
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

func validateRequest(req models.CreatePracticeSessionRequest) error {
	if req.SkillID <= 0 {
		return ErrInvalidSkillID
	}
	if req.UserID <= 0 {
		return ErrInvalidUserID
	}
	if req.DurationMinutes <= 0 {
		return ErrInvalidDuration
	}
	if req.PracticedAt.IsZero() {
		return ErrInvalidPracticedAt
	}
	return nil
}
