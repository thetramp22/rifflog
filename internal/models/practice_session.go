package models

import "time"

type PracticeSession struct {
	ID              int       `db:"id" json:"id"`
	SkillID         int       `db:"skill_id" json:"skill_id"`
	DurationMinutes int       `db:"duration_minutes" json:"duration_minutes"`
	Notes           *string   `db:"notes" json:"notes"`
	PracticedAt     time.Time `db:"practiced_at" json:"practiced_at"`
	CreatedAt       time.Time `db:"created_at" json:"created_at"`
	UserID          int       `db:"user_id" json:"user_id"`
}
