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
var ErrSkillNotFound = errors.New("Skill not found")
var ErrUserNotFound = errors.New("User not found")

type PracticeSessionService struct {
	Repo *repository.PracticeSessionRepository
}

func NewPracticeSessionService(repo *repository.PracticeSessionRepository) *PracticeSessionService {
	return &PracticeSessionService{Repo: repo}
}

func (s *PracticeSessionService) CreatePracticeSession(ctx context.Context, req models.CreatePracticeSessionRequest) (models.PracticeSession, error) {
	err := validateRequest(req)
	if err != nil {
		return models.PracticeSession{}, err
	}

	practiceSession := models.PracticeSession{
		SkillID:         req.SkillID,
		DurationMinutes: req.DurationMinutes,
		PracticedAt:     req.PracticedAt,
		Notes:           req.Notes,
		UserID:          req.UserID,
	}

	returnedSession, err := s.Repo.CreatePracticeSession(ctx, practiceSession)
	if err != nil {
		if errors.Is(err, repository.ErrSkillNotFound) {
			return models.PracticeSession{}, ErrSkillNotFound
		}
		if errors.Is(err, repository.ErrUserNotFound) {
			return models.PracticeSession{}, ErrUserNotFound
		}
		return models.PracticeSession{}, err
	}
	return returnedSession, nil
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

func (s *PracticeSessionService) GetPracticeSessions(ctx context.Context, userID int, params models.FilterParams) ([]models.PracticeSessionDetails, error) {
	practiceSessionDetails, err := s.Repo.GetPracticeSessions(ctx, userID, params)
	if err != nil {
		return nil, err
	}

	return practiceSessionDetails, nil
}

func (s *PracticeSessionService) GetPracticeSessionStats(ctx context.Context, userID int) (models.PracticeSessionStats, error) {
	return s.Repo.GetPracticeSessionStats(ctx, userID)
}
