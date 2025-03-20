package service

import (
	"awesomeProject22/db-service/internal/domain"
	"errors"
	"fmt"
	"github.com/google/uuid"
)

var (
	ErrNotAuthorized = errors.New("user not authorized to perform this action")
	ErrNotFound      = errors.New("entity not found")
)

type BookService struct {
	bookRepo domain.BookRepository
	userSrv  domain.UserService
}

func NewBookService(bookRepo domain.BookRepository, userSrv domain.UserService) *BookService {
	return &BookService{
		bookRepo: bookRepo,
		userSrv:  userSrv,
	}
}

func (s *BookService) Create(book *domain.Book, userID uuid.UUID) error {
	// Проверяем, что пользователь является админом
	isAdmin, err := s.userSrv.IsAdmin(userID)
	if err != nil {
		return fmt.Errorf("checking admin permissions: %w", err)
	}

	if !isAdmin {
		return ErrNotAuthorized
	}

	return s.bookRepo.Create(book)
}

func (s *BookService) GetByID(id uuid.UUID) (*domain.Book, error) {
	return s.bookRepo.GetByID(id)
}

func (s *BookService) GetAll(filters map[string]interface{}) ([]*domain.Book, error) {
	return s.bookRepo.GetAll(filters)
}

func (s *BookService) Update(book *domain.Book, userID uuid.UUID) error {
	// Проверяем, что пользователь является админом
	isAdmin, err := s.userSrv.IsAdmin(userID)
	if err != nil {
		return fmt.Errorf("checking admin permissions: %w", err)
	}

	if !isAdmin {
		return ErrNotAuthorized
	}

	return s.bookRepo.Update(book)
}

func (s *BookService) Delete(id uuid.UUID, userID uuid.UUID) error {
	// Проверяем, что пользователь является админом
	isAdmin, err := s.userSrv.IsAdmin(userID)
	if err != nil {
		return fmt.Errorf("checking admin permissions: %w", err)
	}

	if !isAdmin {
		return ErrNotAuthorized
	}

	return s.bookRepo.Delete(id)
}
