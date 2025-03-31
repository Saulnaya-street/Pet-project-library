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

// BookRepositoryImpl - реализация репозитория книг
type BookRepositoryImpl struct {
	db *pgxpool.Pool
}

// NewBookRepository - конструктор репозитория книг
func NewBookRepository(db *pgxpool.Pool) IBookRepository {
	return &BookRepositoryImpl{
		db: db,
	}
}

func (r *BookRepositoryImpl) Create(book *domain.Book) error {
	if book.ID == uuid.Nil {
		book.ID = uuid.New()
	}

	query := `INSERT INTO books (id, genre, name, author, year) VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.Exec(context.Background(), query,
		book.ID, book.Genre, book.Name, book.Author, book.Year)
	if err != nil {
		return fmt.Errorf("ошибка при создании книги: %w", err)
	}
	return nil
}

func (r *BookRepositoryImpl) GetByID(id uuid.UUID) (*domain.Book, error) {
	var book domain.Book
	query := `SELECT id, genre, name, author, year FROM books WHERE id = $1`

	err := r.db.QueryRow(context.Background(), query, id).Scan(
		&book.ID, &book.Genre, &book.Name, &book.Author, &book.Year)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("книга с ID %s не найдена", id)
		}
		return nil, fmt.Errorf("ошибка при запросе книги с ID %s: %w", id, err)
	}

	return &book, nil
}

func (r *BookRepositoryImpl) GetAll(author, genre string) ([]*domain.Book, error) {
	books := []*domain.Book{}

	// Начинаем формировать базовый запрос
	query := `SELECT id, genre, name, author, year FROM books`
	params := []interface{}{}
	var conditions []string
	paramIndex := 1

	// Добавляем условия, если они заданы
	if author != "" {
		conditions = append(conditions, fmt.Sprintf("author = $%d", paramIndex))
		params = append(params, author)
		paramIndex++
	}

	if genre != "" {
		conditions = append(conditions, fmt.Sprintf("genre = $%d", paramIndex))
		params = append(params, genre)
	}

	// Добавляем WHERE, только если есть условия
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	// Выполняем запрос
	rows, err := r.db.Query(context.Background(), query, params...)
	if err != nil {
		return nil, fmt.Errorf("ошибка при выполнении запроса GetAll: %w", err)
	}
	defer rows.Close()

	// Обрабатываем результаты
	for rows.Next() {
		var book domain.Book
		err := rows.Scan(&book.ID, &book.Genre, &book.Name, &book.Author, &book.Year)
		if err != nil {
			return nil, fmt.Errorf("ошибка при сканировании результатов: %w", err)
		}
		books = append(books, &book)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при обработке результатов: %w", err)
	}

	return books, nil
}

func (r *BookRepositoryImpl) Update(book *domain.Book) error {
	query := `UPDATE books SET genre = $1, name = $2, author = $3, year = $4 WHERE id = $5`
	_, err := r.db.Exec(context.Background(), query,
		book.Genre, book.Name, book.Author, book.Year, book.ID)
	if err != nil {
		return fmt.Errorf("ошибка при обновлении книги с ID %s: %w", book.ID, err)
	}
	return nil
}

func (r *BookRepositoryImpl) Delete(id uuid.UUID) error {
	query := `DELETE FROM books WHERE id = $1`
	_, err := r.db.Exec(context.Background(), query, id)
	if err != nil {
		return fmt.Errorf("ошибка при удалении книги с ID %s: %w", id, err)
	}
	return nil
}
