package domain

import (
	"github.com/google/uuid"
)

type UserRepository interface {
	Create(user *User) error
	GetByID(id uuid.UUID) (*User, error)
	GetByUsername(username string) (*User, error)
	GetByEmail(email string) (*User, error)
	Update(user *User) error
	Delete(id uuid.UUID) error
}

type UserService interface {
	Create(user *User, password string) error
	Authenticate(username, password string) (string, error)
	GetByID(id uuid.UUID) (*User, error)
	IsAdmin(id uuid.UUID) (bool, error)
}
