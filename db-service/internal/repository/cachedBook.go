package repository

import (
	"awesomeProject22/db-service/internal/cache"
	"awesomeProject22/db-service/internal/domain"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"time"
)

const (
	bookKeyPrefix     = "book:"
	bookListKeyPrefix = "books:"
	cacheTTL          = 30 * time.Minute
)

type CachedBookRepository struct {
	repo        IBookRepository
	redisClient cache.IRedisClient
}

func NewCachedBookRepository(repo IBookRepository, redisClient cache.IRedisClient) IBookRepository {
	return &CachedBookRepository{
		repo:        repo,
		redisClient: redisClient,
	}
}

func getBookKey(id uuid.UUID) string {
	return fmt.Sprintf("%s%s", bookKeyPrefix, id.String())
}

func getBookListKey(author, genre string) string {
	return fmt.Sprintf("%s%s:%s", bookListKeyPrefix, author, genre)
}

func (r *CachedBookRepository) Create(book *domain.Book) error {
	ctx := context.Background()
	return r.CreateWithContext(ctx, book)
}

func (r *CachedBookRepository) CreateWithContext(ctx context.Context, book *domain.Book) error {
	err := r.repo.Create(book)
	if err != nil {
		return err
	}

	bookJson, err := json.Marshal(book)
	if err != nil {
		return fmt.Errorf("error serializing book: %w", err)
	}

	err = r.redisClient.Set(ctx, getBookKey(book.ID), string(bookJson), cacheTTL)
	if err != nil {
		fmt.Printf("Error caching book: %v\n", err)
	}

	return nil
}

func (r *CachedBookRepository) GetByID(id uuid.UUID) (*domain.Book, error) {
	ctx := context.Background()
	return r.GetByIDWithContext(ctx, id)
}

func (r *CachedBookRepository) GetByIDWithContext(ctx context.Context, id uuid.UUID) (*domain.Book, error) {
	bookKey := getBookKey(id)
	cachedBook, err := r.redisClient.Get(ctx, bookKey)

	if err == nil {
		var book domain.Book
		if err := json.Unmarshal([]byte(cachedBook), &book); err == nil {
			return &book, nil
		}
	}

	book, err := r.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	bookJson, err := json.Marshal(book)
	if err == nil {
		r.redisClient.Set(ctx, bookKey, string(bookJson), cacheTTL)
	}

	return book, nil
}

func (r *CachedBookRepository) GetAll(author, genre string) ([]*domain.Book, error) {
	ctx := context.Background()
	return r.GetAllWithContext(ctx, author, genre)
}

func (r *CachedBookRepository) GetAllWithContext(ctx context.Context, author, genre string) ([]*domain.Book, error) {
	listKey := getBookListKey(author, genre)

	cachedList, err := r.redisClient.Get(ctx, listKey)

	if err == nil {
		var books []*domain.Book
		if err := json.Unmarshal([]byte(cachedList), &books); err == nil {
			return books, nil
		}
	}

	books, err := r.repo.GetAll(author, genre)
	if err != nil {
		return nil, err
	}

	booksJson, err := json.Marshal(books)
	if err == nil {
		r.redisClient.Set(ctx, listKey, string(booksJson), cacheTTL)
	}

	return books, nil
}

func (r *CachedBookRepository) Update(book *domain.Book) error {
	ctx := context.Background()
	return r.UpdateWithContext(ctx, book)
}

func (r *CachedBookRepository) UpdateWithContext(ctx context.Context, book *domain.Book) error {
	err := r.repo.Update(book)
	if err != nil {
		return err
	}

	bookJson, err := json.Marshal(book)
	if err == nil {
		r.redisClient.Set(ctx, getBookKey(book.ID), string(bookJson), cacheTTL)
	}
	return nil
}

func (r *CachedBookRepository) Delete(id uuid.UUID) error {
	ctx := context.Background()
	return r.DeleteWithContext(ctx, id)
}

func (r *CachedBookRepository) DeleteWithContext(ctx context.Context, id uuid.UUID) error {
	err := r.repo.Delete(id)
	if err != nil {
		return err
	}

	r.redisClient.Delete(ctx, getBookKey(id))

	return nil
}
