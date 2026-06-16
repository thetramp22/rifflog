package models

type PracticeSessionStats struct {
	TotalMinutes       int                `json:"total_minutes"`
	TotalSessions      int                `json:"total_sessions"`
	MostPracticedSkill MostPracticedSkill `json:"most_practiced_skill"`
	LongestSession     int                `json:"longest_session"`
}

type MostPracticedSkill struct {
	Name         string `json:"name"`
	TotalMinutes int    `json:"total_minutes"`
}
