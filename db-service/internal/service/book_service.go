package service

import (
	"awesomeProject22/db-service/internal/domain"
	"github.com/google/uuid"
)

type BookService struct {
	bookRepo domain.BookRepository
}

func NewBookService(bookRepo domain.BookRepository) *BookService {
	return &BookService{
		bookRepo: bookRepo,
	}
}

// GetByID возвращает книгу по её ID
func (s *BookService) GetByID(id uuid.UUID) (domain.Book, error) {
	book, err := s.bookRepo.GetByID(id)
	if err != nil {
		return domain.Book{}, err
	}
	return *book, nil
}

func (s *BookService) GetAll(author, genre string) ([]domain.Book, error) {
	// Напрямую вызываем репозиторий с теми же параметрами
	ptrBooks, err := s.bookRepo.GetAll(author, genre)
	if err != nil {
		return nil, err
	}

	books := make([]domain.Book, len(ptrBooks))
	for i, ptr := range ptrBooks {
		books[i] = *ptr
	}

	return books, nil
}

func (s *BookService) Create(book *domain.Book) error {
	if book.ID == uuid.Nil {
		book.ID = uuid.New()
	}

	return s.bookRepo.Create(book)
}

func (s *BookService) Update(book *domain.Book) error {
	return s.bookRepo.Update(book)
}

func (s *BookService) Delete(id uuid.UUID) error {
	return s.bookRepo.Delete(id)
}
