package models

import "time"

type PracticeSessionDetailsRequest struct {
	UserID int `json:"user_id" form:"user_id"`
}

type FilterParams struct {
	Skill int        `form:"skill"`
	From  *time.Time `form:"from" time_format:"2006-01-02"`
	To    *time.Time `form:"to" time_format:"2006-01-02"`
}
