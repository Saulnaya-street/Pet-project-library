package domain

import (
	"github.com/google/uuid"
)

// BookRepository определяет методы доступа к хранилищу книг
type BookRepository interface {
	Create(book *Book) error
	GetByID(id uuid.UUID) (*Book, error)
	GetAll(filters map[string]interface{}) ([]*Book, error)
	Update(book *Book) error
	Delete(id uuid.UUID) error
}

// BookService определяет бизнес-логику работы с книгами
type BookService interface {
	Create(book *Book) error
	GetByID(id uuid.UUID) (*Book, error)
	GetAll(filters map[string]interface{}) ([]*Book, error)
	Update(book *Book) error
	Delete(id uuid.UUID) error
}
