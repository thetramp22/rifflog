package models

import "time"

type CreatePracticeSessionRequest struct {
	SkillID         int       `json:"skill_id"`
	DurationMinutes int       `json:"duration_minutes"`
	PracticedAt     time.Time `json:"practiced_at"`
	Notes           string    `json:"notes"`
}
type UpdatePracticeSessionRequest struct {
	SkillID         int       `json:"skill_id"`
	DurationMinutes int       `json:"duration_minutes"`
	PracticedAt     time.Time `json:"practiced_at"`
	Notes           string    `json:"notes"`
}

type FilterParams struct {
	Skill int        `form:"skill"`
	From  *time.Time `form:"from" time_format:"2006-01-02"`
	To    *time.Time `form:"to" time_format:"2006-01-02"`
}

type PracticeSessionURI struct {
	ID int `uri:"id" binding:"required,min=1"`
}

type PracticeSession struct {
	ID              int       `db:"id" json:"id"`
	SkillID         int       `db:"skill_id" json:"skill_id"`
	DurationMinutes int       `db:"duration_minutes" json:"duration_minutes"`
	Notes           string    `db:"notes" json:"notes"`
	PracticedAt     time.Time `db:"practiced_at" json:"practiced_at"`
	CreatedAt       time.Time `db:"created_at" json:"created_at"`
	UserID          int       `db:"user_id" json:"user_id"`
}

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

type PracticeSessionStats struct {
	TotalMinutes       int                 `json:"total_minutes"`
	TotalSessions      int                 `json:"total_sessions"`
	MostPracticedSkill *MostPracticedSkill `json:"most_practiced_skill"`
	LongestSession     int                 `json:"longest_session"`
}

type MostPracticedSkill struct {
	Name         string `json:"name"`
	TotalMinutes int    `json:"total_minutes"`
}
