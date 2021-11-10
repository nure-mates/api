package models

import (
	"time"
)

type User struct {
	ID        int       `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	Archived  bool      `json:"-"`
	IsNew     bool      `json:"-" bun:"-"`
	CreatedAt time.Time `json:"-"`
}
