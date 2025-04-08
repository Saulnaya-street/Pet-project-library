package service

import (
	"awesomeProject22/db-service/internal/domain"
	"context"
	"github.com/google/uuid"
)

type IUserService interface {
	GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	Create(ctx context.Context, user *domain.User, password string) error
	Update(ctx context.Context, user *domain.User, passwordChanged bool, newPassword string) error
	Delete(ctx context.Context, id uuid.UUID) error
	Authenticate(ctx context.Context, username, password string) (string, error)
	IsAdmin(ctx context.Context, id uuid.UUID) (bool, error)
	GetAll(ctx context.Context) ([]domain.User, error)
}

type IBookService interface {
	GetAll(ctx context.Context, author, genre string) ([]*domain.Book, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Book, error)
	Create(ctx context.Context, book *domain.Book) error
	Update(ctx context.Context, book *domain.Book) error
	Delete(ctx context.Context, id uuid.UUID) error
}
