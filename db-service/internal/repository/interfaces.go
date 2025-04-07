package repository

import (
	"awesomeProject22/db-service/internal/domain"
	"context"
	"github.com/google/uuid"
)

type IUserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	GetByUsername(ctx context.Context, username string) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetAll(ctx context.Context) ([]domain.User, error)
}

type IBookRepository interface {
	Create(ctx context.Context, book *domain.Book) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Book, error)
	GetAll(ctx context.Context, author, genre string) ([]*domain.Book, error)
	Update(ctx context.Context, book *domain.Book) error
	Delete(ctx context.Context, id uuid.UUID) error
}
