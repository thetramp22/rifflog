package models

import "time"

type SessionsRequest struct {
	SkillID         int       `json:"skill_id"`
	DurationMinutes int       `json:"duration_minutes"`
	PracticedAt     time.Time `json:"practiced_at"`
	Notes           string    `json:"notes"`
	UserID          int       `json:"user_id"`
	// TODO: Remove UserID from request body once authentication is implemented.
	// Session ownership should be derived from the authenticated user context.
}
