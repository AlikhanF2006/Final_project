package model

import "time"

type User struct {
	ID           int
	Username     string
	Email        string
	PasswordHash string
	Role         string
	CreatedAt    time.Time
}
