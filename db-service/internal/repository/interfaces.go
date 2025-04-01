package repository

import (
	"awesomeProject22/db-service/internal/domain"
	"github.com/google/uuid"
)

type IUserRepository interface {
	Create(user *domain.User) error
	GetByID(id uuid.UUID) (*domain.User, error)
	GetByUsername(username string) (*domain.User, error)
	GetByEmail(email string) (*domain.User, error)
	Update(user *domain.User) error
	Delete(id uuid.UUID) error
	GetAll() ([]domain.User, error)
}

type IBookRepository interface {
	Create(book *domain.Book) error
	GetByID(id uuid.UUID) (*domain.Book, error)
	GetAll(author, genre string) ([]*domain.Book, error)
	Update(book *domain.Book) error
	Delete(id uuid.UUID) error
}
