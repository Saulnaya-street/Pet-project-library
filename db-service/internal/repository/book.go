package repository

import (
	"awesomeProject22/db-service/internal/domain"
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"strings"
)

type BookRepositoryImpl struct {
	db *pgxpool.Pool
}

func NewBookRepository(db *pgxpool.Pool) IBookRepository {
	return &BookRepositoryImpl{
		db: db,
	}
}

func (r *BookRepositoryImpl) Create(ctx context.Context, book *domain.Book) error {
	if book.ID == uuid.Nil {
		book.ID = uuid.New()
	}

	query := `INSERT INTO books (id, genre, name, author, year) VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.Exec(ctx, query,
		book.ID, book.Genre, book.Name, book.Author, book.Year)
	if err != nil {
		return fmt.Errorf("error creating book: %w", err)
	}
	return nil
}

func (r *BookRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*domain.Book, error) {
	var book domain.Book
	query := `SELECT id, genre, name, author, year FROM books WHERE id = $1`

	err := r.db.QueryRow(ctx, query, id).Scan(
		&book.ID, &book.Genre, &book.Name, &book.Author, &book.Year)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("book with ID %s not found", id)
		}
		return nil, fmt.Errorf("error requesting book with ID %s: %w", id, err)
	}

	return &book, nil
}

func (r *BookRepositoryImpl) GetAll(ctx context.Context, author, genre string) ([]*domain.Book, error) {
	books := []*domain.Book{}

	query := `SELECT id, genre, name, author, year FROM books`
	params := []interface{}{}
	var conditions []string
	paramIndex := 1

	if author != "" {
		conditions = append(conditions, fmt.Sprintf("author = $%d", paramIndex))
		params = append(params, author)
		paramIndex++
	}

	if genre != "" {
		conditions = append(conditions, fmt.Sprintf("genre = $%d", paramIndex))
		params = append(params, genre)
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	rows, err := r.db.Query(ctx, query, params...)
	if err != nil {
		return nil, fmt.Errorf("error while executing GetAll request: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var book domain.Book
		err := rows.Scan(&book.ID, &book.Genre, &book.Name, &book.Author, &book.Year)
		if err != nil {
			return nil, fmt.Errorf("error scanning results: %w", err)
		}
		books = append(books, &book)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error when processing results: %w", err)
	}

	return books, nil
}

func (r *BookRepositoryImpl) Update(ctx context.Context, book *domain.Book) error {
	query := `UPDATE books SET genre = $1, name = $2, author = $3, year = $4 WHERE id = $5`
	_, err := r.db.Exec(ctx, query,
		book.Genre, book.Name, book.Author, book.Year, book.ID)
	if err != nil {
		return fmt.Errorf("error updating book with ID %s: %w", book.ID, err)
	}
	return nil
}

func (r *BookRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM books WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("error deleting book with ID %s: %w", id, err)
	}
	return nil
}
