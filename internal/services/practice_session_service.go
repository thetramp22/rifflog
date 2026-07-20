package services

import (
	"context"
	"errors"

	"github.com/thetramp22/rifflog/internal/models"
	repository "github.com/thetramp22/rifflog/internal/repositories"
)

var ErrInvalidDuration = errors.New("Duration must be greater than zero")
var ErrInvalidSkillID = errors.New("Invalid skill id")
var ErrInvalidUserID = errors.New("Invalid user id")
var ErrInvalidPracticedAt = errors.New("Invalid practiced at time")
var ErrSkillNotFound = errors.New("Skill not found")
var ErrUserNotFound = errors.New("User not found")
var ErrPracticeSessionNotFound = errors.New("practice session not found")

// PracticeSessionService provides methods dealing with practice sessions.
type PracticeSessionService struct {
	Repo *repository.PracticeSessionRepository
}

// NewPracticeSessionService returns a PracticeSessionService.
func NewPracticeSessionService(repo *repository.PracticeSessionRepository) *PracticeSessionService {
	return &PracticeSessionService{Repo: repo}
}

// CreatePracticeSession takes a request from a handler, validates the request, and calls the repository method to create the session.
// Returns the created practice session on success.
func (s *PracticeSessionService) CreatePracticeSession(ctx context.Context, userID int, req models.CreatePracticeSessionRequest) (models.PracticeSession, error) {
	err := validateCreateSessionRequest(userID, req)
	if err != nil {
		return models.PracticeSession{}, err
	}

	practiceSession := models.PracticeSession{
		SkillID:         req.SkillID,
		DurationMinutes: req.DurationMinutes,
		PracticedAt:     req.PracticedAt,
		Notes:           req.Notes,
		UserID:          userID,
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

// PracticeSessionService takes a request from a handler, validates the request, and calls the repository method to update the session.
// Returns the updated practice session on success.
func (s *PracticeSessionService) UpdatePracticeSession(ctx context.Context, userID int, sessionID int, req models.UpdatePracticeSessionRequest) (models.PracticeSession, error) {
	err := validateUpdateSessionRequest(userID, req)
	if err != nil {
		return models.PracticeSession{}, err
	}

	practiceSession := models.PracticeSession{
		SkillID:         req.SkillID,
		DurationMinutes: req.DurationMinutes,
		PracticedAt:     req.PracticedAt,
		Notes:           req.Notes,
		UserID:          userID,
	}

	returnedSession, err := s.Repo.UpdatePracticeSession(ctx, userID, sessionID, practiceSession)
	if err != nil {
		if errors.Is(err, repository.ErrSkillNotFound) {
			return models.PracticeSession{}, ErrSkillNotFound
		}
		if errors.Is(err, repository.ErrUserNotFound) {
			return models.PracticeSession{}, ErrUserNotFound
		}
		if errors.Is(err, repository.ErrPracticeSessionNotFound) {
			return models.PracticeSession{}, ErrPracticeSessionNotFound
		}
		return models.PracticeSession{}, err
	}
	return returnedSession, nil
}

// DeletePracticeSession takes a request from a handler and calls the repository method to delete the given session.
// Returns the deleted practice session id on success.
func (s *PracticeSessionService) DeletePracticeSession(ctx context.Context, userID int, sessionID int) (int, error) {
	deletedSessionID, err := s.Repo.DeletePracticeSession(ctx, userID, sessionID)
	if err != nil {
		if errors.Is(err, repository.ErrPracticeSessionNotFound) {
			return 0, ErrPracticeSessionNotFound
		}
		return 0, err
	}
	return deletedSessionID, nil
}

func validateCreateSessionRequest(userID int, req models.CreatePracticeSessionRequest) error {
	if req.SkillID <= 0 {
		return ErrInvalidSkillID
	}
	if userID <= 0 {
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

func validateUpdateSessionRequest(userID int, req models.UpdatePracticeSessionRequest) error {
	if req.SkillID <= 0 {
		return ErrInvalidSkillID
	}
	if userID <= 0 {
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

// GetPracticeSessions takes a request from a handler and calls the repository method to retrieve the sessions based on the given parameters.
// Returns a list of practice sessions.
func (s *PracticeSessionService) GetPracticeSessions(ctx context.Context, userID int, params models.FilterParams) ([]models.PracticeSessionDetails, error) {
	practiceSessionDetails, err := s.Repo.GetPracticeSessions(ctx, userID, params)
	if err != nil {
		return nil, err
	}

	return practiceSessionDetails, nil
}

// GetPracticeSessionStats takes a request from a handler and calls the repository method to retrieve practice session stats.
// Returns a list of practice sessions.
func (s *PracticeSessionService) GetPracticeSessionStats(ctx context.Context, userID int) (models.PracticeSessionStats, error) {
	return s.Repo.GetPracticeSessionStats(ctx, userID)
}
