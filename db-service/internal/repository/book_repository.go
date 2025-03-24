package repository

import (
	"awesomeProject22/db-service/internal/domain"
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

type BookRepository struct {
	db *pgxpool.Pool
}

func NewBookRepository(db *pgxpool.Pool) *BookRepository {
	return &BookRepository{
		db: db,
	}
}

func (r *BookRepository) Create(book *domain.Book) error {
	if book.ID == uuid.Nil {
		book.ID = uuid.New()
	}

	query := `INSERT INTO books (id, genre, name, author, year) VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.Exec(context.Background(), query,
		book.ID, book.Genre, book.Name, book.Author, book.Year)
	return err
}

func (r *BookRepository) GetByID(id uuid.UUID) (*domain.Book, error) {
	var book domain.Book
	query := `SELECT id, genre, name, author, year FROM books WHERE id = $1`

	err := r.db.QueryRow(context.Background(), query, id).Scan(
		&book.ID, &book.Genre, &book.Name, &book.Author, &book.Year)
	if err != nil {
		return nil, err
	}

	return &book, nil
}

func (r *BookRepository) GetAll(author, genre string) ([]*domain.Book, error) {
	var books []*domain.Book

	query := `SELECT id, genre, name, author, year FROM books WHERE 1=1`
	params := make([]interface{}, 0)
	paramCount := 1

	if author != "" {
		query += fmt.Sprintf(" AND author = $%d", paramCount)
		params = append(params, author)
		paramCount++
	}

	if genre != "" {
		query += fmt.Sprintf(" AND genre = $%d", paramCount)
		params = append(params, genre)
	}

	rows, err := r.db.Query(context.Background(), query, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var book domain.Book
		err := rows.Scan(&book.ID, &book.Genre, &book.Name, &book.Author, &book.Year)
		if err != nil {
			return nil, err
		}
		books = append(books, &book)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return books, nil
}

func (r *BookRepository) Update(book *domain.Book) error {
	query := `UPDATE books SET genre = $1, name = $2, author = $3, year = $4 WHERE id = $5`
	_, err := r.db.Exec(context.Background(), query,
		book.Genre, book.Name, book.Author, book.Year, book.ID)
	return err
}

func (r *BookRepository) Delete(id uuid.UUID) error {
	query := `DELETE FROM books WHERE id = $1`
	_, err := r.db.Exec(context.Background(), query, id)
	return err
}
