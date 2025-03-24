package repository

import (
	"awesomeProject22/db-service/internal/domain"
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) Create(user *domain.User) error {
	// Генерируем UUID, если он не был установлен
	if user.ID == uuid.Nil {
		user.ID = uuid.New()
	}

	query := `INSERT INTO users (id, username, email, password_hash, is_admin) 
              VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.Exec(context.Background(), query,
		user.ID, user.Username, user.Email, user.PasswordHash, user.IsAdmin)
	return err
}

func (r *UserRepository) GetByID(id uuid.UUID) (*domain.User, error) {
	var user domain.User
	query := `SELECT id, username, email, password_hash, is_admin FROM users WHERE id = $1`

	err := r.db.QueryRow(context.Background(), query, id).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.IsAdmin)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetByUsername(username string) (*domain.User, error) {
	var user domain.User
	query := `SELECT id, username, email, password_hash, is_admin FROM users WHERE username = $1`

	err := r.db.QueryRow(context.Background(), query, username).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.IsAdmin)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetByEmail(email string) (*domain.User, error) {
	var user domain.User
	query := `SELECT id, username, email, password_hash, is_admin FROM users WHERE email = $1`

	err := r.db.QueryRow(context.Background(), query, email).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.IsAdmin)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) Update(user *domain.User) error {
	query := `UPDATE users SET username = $1, email = $2, password_hash = $3, is_admin = $4 WHERE id = $5`
	_, err := r.db.Exec(context.Background(), query,
		user.Username, user.Email, user.PasswordHash, user.IsAdmin, user.ID)
	return err
}

func (r *UserRepository) Delete(id uuid.UUID) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.Exec(context.Background(), query, id)
	return err
}

func (r *UserRepository) GetAll() ([]*domain.User, error) {
	query := `SELECT id, username, email, password_hash, is_admin FROM users`

	rows, err := r.db.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []*domain.User{}

	for rows.Next() {
		var user domain.User
		err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.IsAdmin)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	return users, nil
}
