package service

import (
	"awesomeProject22/db-service/internal/domain"
	"github.com/google/uuid"
)

type IUserService interface {
	GetByID(id uuid.UUID) (*domain.User, error)
	Create(user *domain.User, password string) error
	Update(user *domain.User, passwordChanged bool, newPassword string) error
	Delete(id uuid.UUID) error
	Authenticate(username, password string) (string, error)
	IsAdmin(id uuid.UUID) (bool, error)
	GetAll() ([]domain.User, error)
}

type IBookService interface {
	GetAll(author, genre string) ([]*domain.Book, error)
	GetByID(id uuid.UUID) (*domain.Book, error)
	Create(book *domain.Book) error
	Update(book *domain.Book) error
	Delete(id uuid.UUID) error
}
