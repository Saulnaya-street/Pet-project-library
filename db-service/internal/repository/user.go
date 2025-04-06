package repository

import (
	"awesomeProject22/db-service/internal/domain"
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type UserRepositoryImpl struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) IUserRepository {
	return &UserRepositoryImpl{
		db: db,
	}
}

func (r *UserRepositoryImpl) Create(user *domain.User) error {
	if user.ID == uuid.Nil {
		user.ID = uuid.New()
	}

	query := `INSERT INTO users (id, username, email, password_hash, is_admin) 
              VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.Exec(context.Background(), query,
		user.ID, user.Username, user.Email, user.PasswordHash, user.IsAdmin)
	if err != nil {
		return fmt.Errorf("ошибка при создании пользователя: %w", err)
	}
	return nil
}

func (r *UserRepositoryImpl) GetByID(id uuid.UUID) (*domain.User, error) {
	var user domain.User
	query := `SELECT id, username, email, is_admin FROM users WHERE id = $1`

	err := r.db.QueryRow(context.Background(), query, id).Scan(
		&user.ID, &user.Username, &user.Email, &user.IsAdmin)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("пользователь с ID %s не найден", id)
		}
		return nil, fmt.Errorf("ошибка при запросе пользователя с ID %s: %w", id, err)
	}

	return &user, nil
}

func (r *UserRepositoryImpl) GetByUsername(username string) (*domain.User, error) {
	var user domain.User
	query := `SELECT id, username, email, password_hash, is_admin FROM users WHERE username = $1`

	err := r.db.QueryRow(context.Background(), query, username).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.IsAdmin)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("пользователь с именем %s не найден", username)
		}
		return nil, fmt.Errorf("ошибка при запросе пользователя с именем %s: %w", username, err)
	}

	return &user, nil
}

func (r *UserRepositoryImpl) GetByEmail(email string) (*domain.User, error) {
	var user domain.User

	query := `SELECT id, username, email, is_admin FROM users WHERE email = $1`

	err := r.db.QueryRow(context.Background(), query, email).Scan(
		&user.ID, &user.Username, &user.Email, &user.IsAdmin)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("пользователь с email %s не найден", email)
		}
		return nil, fmt.Errorf("ошибка при запросе пользователя с email %s: %w", email, err)
	}

	return &user, nil
}

func (r *UserRepositoryImpl) Update(user *domain.User) error {
	query := `UPDATE users SET username = $1, email = $2, password_hash = $3, is_admin = $4 WHERE id = $5`
	_, err := r.db.Exec(context.Background(), query,
		user.Username, user.Email, user.PasswordHash, user.IsAdmin, user.ID)
	if err != nil {
		return fmt.Errorf("ошибка при обновлении пользователя с ID %s: %w", user.ID, err)
	}
	return nil
}

func (r *UserRepositoryImpl) Delete(id uuid.UUID) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.Exec(context.Background(), query, id)
	if err != nil {
		return fmt.Errorf("ошибка при удалении пользователя с ID %s: %w", id, err)
	}
	return nil
}

func (r *UserRepositoryImpl) GetAll() ([]domain.User, error) {

	query := `SELECT id, username, email, is_admin FROM users`

	rows, err := r.db.Query(context.Background(), query)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении списка пользователей: %w", err)
	}
	defer rows.Close()

	users := []domain.User{}

	for rows.Next() {
		var user domain.User
		err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.IsAdmin)
		if err != nil {
			return nil, fmt.Errorf("ошибка при сканировании данных пользователя: %w", err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка после обработки результатов: %w", err)
	}

	return users, nil
}
