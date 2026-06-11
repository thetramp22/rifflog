package models

import "time"

type PracticeSessionDetails struct {
	ID               int       `db:"session_id" json:"session_id"`
	SkillID          int       `db:"skill_id" json:"skill_id"`
	SkillName        string    `db:"skill_name" json:"skill_name"`
	SkillDescription string    `db:"skill_description" json:"skill_description"`
	DurationMinutes  int       `db:"duration_minutes" json:"duration_minutes"`
	Notes            string    `db:"notes" json:"notes"`
	PracticedAt      time.Time `db:"practiced_at" json:"practiced_at"`
	CreatedAt        time.Time `db:"created_at" json:"created_at"`
	UserID           int       `db:"user_id" json:"user_id"`
}
