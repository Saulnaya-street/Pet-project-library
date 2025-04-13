package repository

import (
	"awesomeProject22/db-service/internal/cache"
	"awesomeProject22/db-service/internal/domain"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
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

func CachedBookRepo(repo IBookRepository, redisClient cache.IRedisClient) IBookRepository {
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

func (r *CachedBookRepository) Create(ctx context.Context, book *domain.Book) error {
	err := r.repo.Create(ctx, book)
	if err != nil {
		return fmt.Errorf("error creating book in database: %w", err)
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

func (r *CachedBookRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Book, error) {
	bookKey := getBookKey(id)
	cachedBook, err := r.redisClient.Get(ctx, bookKey)

	if err == nil {
		var book domain.Book
		if err := json.Unmarshal([]byte(cachedBook), &book); err == nil {
			return &book, nil
		} else {
			fmt.Printf("Error unmarshaling cached book: %v\n", err)
		}
	} else if err != redis.Nil {
		fmt.Printf("Error fetching book from Redis: %v\n", err)
	}

	book, err := r.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error getting book from database: %w", err)
	}

	bookJson, err := json.Marshal(book)
	if err != nil {
		fmt.Printf("Error serializing book for cache: %v\n", err)
		return book, nil
	}

	redisErr := r.redisClient.Set(ctx, bookKey, string(bookJson), cacheTTL)
	if redisErr != nil {
		fmt.Printf("Error caching book data: %v\n", redisErr)
	}

	return book, nil
}

func (r *CachedBookRepository) GetAll(ctx context.Context, author, genre string) ([]*domain.Book, error) {
	listKey := getBookListKey(author, genre)

	cachedList, err := r.redisClient.Get(ctx, listKey)

	if err == nil {
		var books []*domain.Book
		if unmarshalErr := json.Unmarshal([]byte(cachedList), &books); unmarshalErr == nil {
			return books, nil
		} else {
			fmt.Printf("Error deserializing book list from cache: %v\n", unmarshalErr)
		}
	} else if err != redis.Nil {
		fmt.Printf("Error getting list of books from Redis: %v\n", err)
	}

	books, err := r.repo.GetAll(ctx, author, genre)
	if err != nil {
		return nil, fmt.Errorf("error getting list of books from database: %w", err)
	}

	booksJson, err := json.Marshal(books)
	if err != nil {
		fmt.Printf("Error serializing book list for cache: %v\n", err)
		return books, nil
	}

	redisErr := r.redisClient.Set(ctx, listKey, string(booksJson), cacheTTL)
	if redisErr != nil {
		fmt.Printf("Book list caching error: %v\n", redisErr)
	}

	return books, nil
}

func (r *CachedBookRepository) Update(ctx context.Context, book *domain.Book) error {
	err := r.repo.Update(ctx, book)
	if err != nil {
		return fmt.Errorf("error updating book in database: %w", err)
	}

	bookJson, err := json.Marshal(book)
	if err != nil {
		fmt.Printf("Error serializing book to update cache: %v\n", err)
		return nil
	}

	redisErr := r.redisClient.Set(ctx, getBookKey(book.ID), string(bookJson), cacheTTL)
	if redisErr != nil {
		fmt.Printf("Error updating book in cache: %v\n", redisErr)
	}

	return nil
}

func (r *CachedBookRepository) Delete(ctx context.Context, id uuid.UUID) error {
	err := r.repo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("error deleting book from database: %w", err)
	}

	redisErr := r.redisClient.Delete(ctx, getBookKey(id))
	if redisErr != nil {
		fmt.Printf("Error deleting book from cache: %v\n", redisErr)
	}

	return nil
}
