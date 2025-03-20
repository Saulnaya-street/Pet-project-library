package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"awesomeProject22/db-service/internal/domain"
)

type BookRepository struct {
	db *sqlx.DB
}

func NewBookRepository(db *sqlx.DB) *BookRepository {
	return &BookRepository{
		db: db,
	}
}

func (r *BookRepository) Create(book *domain.Book) error {
	if book.ID == uuid.Nil {
		book.ID = uuid.New()
	}

	query := `INSERT INTO books (id, genre, name, author, year) 
              VALUES ($1, $2, $3, $4, $5)`

	_, err := r.db.Exec(query, book.ID, book.Genre, book.Name, book.Author, book.Year)
	return err
}

func (r *BookRepository) GetByID(id uuid.UUID) (*domain.Book, error) {
	var book domain.Book

	query := `SELECT id, genre, name, author, year FROM books WHERE id = $1`

	err := r.db.Get(&book, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("book not found: %w", err)
		}
		return nil, err
	}

	return &book, nil
}

func (r *BookRepository) GetAll(filters map[string]interface{}) ([]*domain.Book, error) {
	books := []*domain.Book{}

	query := `SELECT id, genre, name, author, year FROM books`

	whereConditions := []string{}
	args := []interface{}{}
	argPosition := 1

	if author, ok := filters["author"].(string); ok && author != "" {
		whereConditions = append(whereConditions, fmt.Sprintf("author = $%d", argPosition))
		args = append(args, author)
		argPosition++
	}

	if genre, ok := filters["genre"].(string); ok && genre != "" {
		whereConditions = append(whereConditions, fmt.Sprintf("genre = $%d", argPosition))
		args = append(args, genre)
		argPosition++
	}

	if year, ok := filters["year"].(int); ok && year != 0 {
		whereConditions = append(whereConditions, fmt.Sprintf("year = $%d", argPosition))
		args = append(args, year)
		argPosition++
	}

	if len(whereConditions) > 0 {
		query += " WHERE " + strings.Join(whereConditions, " AND ")
	}

	err := r.db.Select(&books, query, args...)
	if err != nil {
		return nil, err
	}

	return books, nil
}

func (r *BookRepository) Update(book *domain.Book) error {
	query := `UPDATE books SET genre = $1, name = $2, author = $3, year = $4 
              WHERE id = $5`

	result, err := r.db.Exec(query, book.Genre, book.Name, book.Author, book.Year, book.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("book with id %s not found", book.ID)
	}

	return nil
}

func (r *BookRepository) Delete(id uuid.UUID) error {
	query := `DELETE FROM books WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("book with id %s not found", id)
	}

	return nil
}
