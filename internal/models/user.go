package models

import "time"

type User struct {
	ID           int32
	Name         string
	Surname      string
	Nickname     string
	Email        string
	PasswordHash string
	CreatedAt    time.Time
}
