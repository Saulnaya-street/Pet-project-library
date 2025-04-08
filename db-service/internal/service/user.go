package service

import (
	"awesomeProject22/db-service/internal/domain"
	"awesomeProject22/db-service/internal/repository"
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserServiceImpl struct {
	repo repository.IUserRepository
}

func NewUserService(repo repository.IUserRepository) IUserService {
	return &UserServiceImpl{
		repo: repo,
	}
}

func (s *UserServiceImpl) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *UserServiceImpl) Create(ctx context.Context, user *domain.User, password string) error {
	if user.ID == uuid.Nil {
		user.ID = uuid.New()
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("ошибка при хешировании пароля: %w", err)
	}

	user.PasswordHash = string(hashedPassword)

	return s.repo.Create(ctx, user)
}

func (s *UserServiceImpl) Update(ctx context.Context, user *domain.User, passwordChanged bool, newPassword string) error {
	_, err := s.repo.GetByID(ctx, user.ID)
	if err != nil {
		return fmt.Errorf("пользователь для обновления не найден: %w", err)
	}

	if passwordChanged {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("ошибка при хешировании пароля: %w", err)
		}
		user.PasswordHash = string(hashedPassword)
	} else {
		userWithPass, err := s.getUserWithPassword(ctx, user.ID)
		if err != nil {
			return fmt.Errorf("ошибка при получении пароля пользователя: %w", err)
		}
		user.PasswordHash = userWithPass.PasswordHash
	}

	return s.repo.Update(ctx, user)
}

func (s *UserServiceImpl) getUserWithPassword(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	var user domain.User

	userWithUsername, err := s.repo.GetByUsername(ctx, id.String())
	if err == nil && userWithUsername != nil {
		return userWithUsername, nil
	}

	return &user, fmt.Errorf("не удалось получить пароль пользователя")
}

func (s *UserServiceImpl) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("пользователь для удаления не найден: %w", err)
	}

	return s.repo.Delete(ctx, id)
}

func (s *UserServiceImpl) Authenticate(ctx context.Context, username, password string) (string, error) {
	user, err := s.repo.GetByUsername(ctx, username)
	if err != nil {
		return "", errors.New("неверные учетные данные")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", errors.New("неверные учетные данные")
	}

	return user.ID.String(), nil
}

func (s *UserServiceImpl) IsAdmin(ctx context.Context, id uuid.UUID) (bool, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return false, fmt.Errorf("ошибка при проверке прав администратора: %w", err)
	}

	return user.IsAdmin, nil
}

func (s *UserServiceImpl) GetAll(ctx context.Context) ([]domain.User, error) {
	return s.repo.GetAll(ctx)
}
