package domain

import (
	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `json:"id" db:"id"`
	Username     string    `json:"username" db:"username"`
	Email        string    `json:"email" db:"email"`
	PasswordHash string    `json:"password_hash" db:"password_hash"`
	IsAdmin      bool      `json:"is_admin" db:"is_admin"`
}

type Book struct {
	ID     uuid.UUID `json:"id" db:"id"`
	Genre  string    `json:"genre" db:"genre"`
	Name   string    `json:"name" db:"name"`
	Author string    `json:"author" db:"author"`
	Year   int       `json:"year" db:"year"`
}

type UserBook struct {
	UserID uuid.UUID `json:"user_id" db:"user_id"`
	BookID uuid.UUID `json:"book_id" db:"book_id"`
}
