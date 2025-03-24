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
	GetAll() ([]*User, error)
}

type UserService interface {
	GetByID(id uuid.UUID) (*User, error)
	Create(user *User, password string) error
	Update(user *User, passwordChanged bool, newPassword string) error
	Delete(id uuid.UUID) error
	Authenticate(username, password string) (string, error)
	IsAdmin(id uuid.UUID) (bool, error)
	GetAll() ([]*User, error)
}
