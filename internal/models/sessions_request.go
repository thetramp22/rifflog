package models

import "time"

type SessionsRequest struct {
	SkillID         string    `json:"skill_id"`
	DurationMinutes int       `json:"duration_minutes"`
	PracticedAt     time.Time `json:"practiced_at"`
	Notes           string    `json:"notes"`
}
