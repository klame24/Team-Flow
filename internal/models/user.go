package models

import "time"

type User struct {
	ID           int
	Name         string
	Surname      string
	Nickname     string
	Email        string
	PasswordHash string
	CreatedAt    time.Time
}
