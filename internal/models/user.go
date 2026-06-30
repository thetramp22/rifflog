package models

import "time"

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type User struct {
	ID           int       `db:"id"`
	Email        string    `db:"email"`
	PasswordHash string    `db:"password_hash"`
	CreatedAt    time.Time `db:"created_at"`
}
