package service

import (
	"awesomeProject22/db-service/internal/domain"
	"awesomeProject22/db-service/internal/repository"
	"fmt"
	"github.com/google/uuid"
)

// BookServiceImpl - реализация сервиса книг
type BookServiceImpl struct {
	bookRepo repository.IBookRepository
}

// NewBookService - конструктор сервиса книг
func NewBookService(bookRepo repository.IBookRepository) IBookService {
	return &BookServiceImpl{
		bookRepo: bookRepo,
	}
}

func (s *BookServiceImpl) GetByID(id uuid.UUID) (*domain.Book, error) {
	return s.bookRepo.GetByID(id)
}

func (s *BookServiceImpl) GetAll(author, genre string) ([]*domain.Book, error) {
	return s.bookRepo.GetAll(author, genre)
}

func (s *BookServiceImpl) Create(book *domain.Book) error {
	if book.ID == uuid.Nil {
		book.ID = uuid.New()
	}

	return s.bookRepo.Create(book)
}

func (s *BookServiceImpl) Update(book *domain.Book) error {
	// Проверяем существование книги перед обновлением
	_, err := s.bookRepo.GetByID(book.ID)
	if err != nil {
		return fmt.Errorf("книга для обновления не найдена: %w", err)
	}

	return s.bookRepo.Update(book)
}

func (s *BookServiceImpl) Delete(id uuid.UUID) error {
	// Проверяем существование книги перед удалением
	_, err := s.bookRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("книга для удаления не найдена: %w", err)
	}

	return s.bookRepo.Delete(id)
}
