package domain

import (
	"github.com/google/uuid"
)

type BookRepository interface {
	Create(book *Book) error
	GetByID(id uuid.UUID) (*Book, error)
	GetAll(author, genre string) ([]*Book, error)
	Update(book *Book) error
	Delete(id uuid.UUID) error
}

type BookService interface {
	GetAll(author, genre string) ([]Book, error)
	GetByID(id uuid.UUID) (Book, error)
	Create(book *Book) error
	Update(book *Book) error
	Delete(id uuid.UUID) error
}

//да ты щас скажешь что должны быть там где используются я переделаю но позже мне нужно понять что с другим кодом у меня
