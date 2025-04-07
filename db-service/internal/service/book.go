package service

import (
	"awesomeProject22/db-service/internal/domain"
	"awesomeProject22/db-service/internal/repository"
	"context"
	"fmt"
	"github.com/google/uuid"
)

type BookServiceImpl struct {
	bookRepo repository.IBookRepository
}

func NewBookService(bookRepo repository.IBookRepository) IBookService {
	return &BookServiceImpl{
		bookRepo: bookRepo,
	}
}

func (s *BookServiceImpl) GetByID(ctx context.Context, id uuid.UUID) (*domain.Book, error) {
	return s.bookRepo.GetByID(ctx, id)
}

func (s *BookServiceImpl) GetAll(ctx context.Context, author, genre string) ([]*domain.Book, error) {
	return s.bookRepo.GetAll(ctx, author, genre)
}

func (s *BookServiceImpl) Create(ctx context.Context, book *domain.Book) error {
	if book.ID == uuid.Nil {
		book.ID = uuid.New()
	}

	return s.bookRepo.Create(ctx, book)
}

func (s *BookServiceImpl) Update(ctx context.Context, book *domain.Book) error {
	_, err := s.bookRepo.GetByID(ctx, book.ID)
	if err != nil {
		return fmt.Errorf("книга для обновления не найдена: %w", err)
	}

	return s.bookRepo.Update(ctx, book)
}

func (s *BookServiceImpl) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := s.bookRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("книга для удаления не найдена: %w", err)
	}

	return s.bookRepo.Delete(ctx, id)
}
