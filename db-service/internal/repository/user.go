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

func UserRepo(db *pgxpool.Pool) IUserRepository {
	return &UserRepositoryImpl{
		db: db,
	}
}

func (r *UserRepositoryImpl) Create(ctx context.Context, user *domain.User) error {
	if user.ID == uuid.Nil {
		user.ID = uuid.New()
	}

	query := `INSERT INTO users (id, username, email, password_hash, is_admin) 
              VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.Exec(ctx, query,
		user.ID, user.Username, user.Email, user.PasswordHash, user.IsAdmin)
	if err != nil {
		return fmt.Errorf("error creating user: %w", err)
	}
	return nil
}

func (r *UserRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	var user domain.User
	query := `SELECT id, username, email, is_admin FROM users WHERE id = $1`

	err := r.db.QueryRow(ctx, query, id).Scan(
		&user.ID, &user.Username, &user.Email, &user.IsAdmin)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("user with ID %s not found", id)
		}
		return nil, fmt.Errorf("error requesting user with ID %s: %w", id, err)
	}

	return &user, nil
}

func (r *UserRepositoryImpl) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	var user domain.User
	query := `SELECT id, username, email, password_hash, is_admin FROM users WHERE username = $1`

	err := r.db.QueryRow(ctx, query, username).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.IsAdmin)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("username %s not found", username)
		}
		return nil, fmt.Errorf("error requesting user with name %s: %w", username, err)
	}

	return &user, nil
}

func (r *UserRepositoryImpl) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User

	query := `SELECT id, username, email, is_admin FROM users WHERE email = $1`

	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID, &user.Username, &user.Email, &user.IsAdmin)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("user with email %s not found", email)
		}
		return nil, fmt.Errorf("error when requesting user with email %s: %w", email, err)
	}

	return &user, nil
}

func (r *UserRepositoryImpl) Update(ctx context.Context, user *domain.User) error {
	query := `UPDATE users SET username = $1, email = $2, password_hash = $3, is_admin = $4 WHERE id = $5`
	_, err := r.db.Exec(ctx, query,
		user.Username, user.Email, user.PasswordHash, user.IsAdmin, user.ID)
	if err != nil {
		return fmt.Errorf("error updating user with ID %s: %w", user.ID, err)
	}
	return nil
}

func (r *UserRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("error deleting user with ID %s: %w", id, err)
	}
	return nil
}

func (r *UserRepositoryImpl) GetAll(ctx context.Context) ([]domain.User, error) {
	query := `SELECT id, username, email, is_admin FROM users`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error getting list of users: %w", err)
	}
	defer rows.Close()

	users := []domain.User{}

	for rows.Next() {
		var user domain.User
		err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.IsAdmin)
		if err != nil {
			return nil, fmt.Errorf("error scanning user data: %w", err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error after processing results: %w", err)
	}

	return users, nil
}
