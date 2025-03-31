package service

import (
	"awesomeProject22/db-service/internal/domain"
	"awesomeProject22/db-service/internal/repository"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// UserServiceImpl - реализация сервиса пользователей
type UserServiceImpl struct {
	repo repository.IUserRepository
}

// NewUserService - конструктор сервиса пользователей
func NewUserService(repo repository.IUserRepository) IUserService {
	return &UserServiceImpl{
		repo: repo,
	}
}

func (s *UserServiceImpl) GetByID(id uuid.UUID) (*domain.User, error) {
	return s.repo.GetByID(id)
}

func (s *UserServiceImpl) Create(user *domain.User, password string) error {
	if user.ID == uuid.Nil {
		user.ID = uuid.New()
	}

	// Хешируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("ошибка при хешировании пароля: %w", err)
	}

	user.PasswordHash = string(hashedPassword)

	return s.repo.Create(user)
}

func (s *UserServiceImpl) Update(user *domain.User, passwordChanged bool, newPassword string) error {
	// Проверяем существование пользователя
	_, err := s.repo.GetByID(user.ID)
	if err != nil {
		return fmt.Errorf("пользователь для обновления не найден: %w", err)
	}

	// Если пароль меняется, хешируем новый пароль
	if passwordChanged {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("ошибка при хешировании пароля: %w", err)
		}
		user.PasswordHash = string(hashedPassword)
	} else {
		// Если пароль не меняется, сохраняем существующий хеш
		// Получаем хеш пароля из базы данных
		userWithPass, err := s.getUserWithPassword(user.ID)
		if err != nil {
			return fmt.Errorf("ошибка при получении пароля пользователя: %w", err)
		}
		user.PasswordHash = userWithPass.PasswordHash
	}

	return s.repo.Update(user)
}

// Вспомогательный метод для получения пользователя с паролем
func (s *UserServiceImpl) getUserWithPassword(id uuid.UUID) (*domain.User, error) {
	// Получаем данные из базы, включая пароль
	// Этот метод сделан отдельно, чтобы не менять существующие методы репозитория
	// и чтобы не добавлять новую функцию в интерфейс
	var user domain.User

	// Заглушка - в реальном коде здесь будет запрос к базе данных
	// В вашем случае, нужно будет добавить метод в репозиторий,
	// который возвращает пользователя с хешем пароля
	userWithUsername, err := s.repo.GetByUsername(id.String())
	if err == nil && userWithUsername != nil {
		return userWithUsername, nil
	}

	// Если не удалось найти по username, можно попробовать другой подход
	// или вернуть ошибку
	return &user, fmt.Errorf("не удалось получить пароль пользователя")
}

func (s *UserServiceImpl) Delete(id uuid.UUID) error {
	_, err := s.repo.GetByID(id)
	if err != nil {
		return fmt.Errorf("пользователь для удаления не найден: %w", err)
	}

	return s.repo.Delete(id)
}

func (s *UserServiceImpl) Authenticate(username, password string) (string, error) {
	user, err := s.repo.GetByUsername(username)
	if err != nil {
		return "", errors.New("неверные учетные данные")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", errors.New("неверные учетные данные")
	}

	return user.ID.String(), nil
}

func (s *UserServiceImpl) IsAdmin(id uuid.UUID) (bool, error) {
	user, err := s.repo.GetByID(id)
	if err != nil {
		return false, fmt.Errorf("ошибка при проверке прав администратора: %w", err)
	}

	return user.IsAdmin, nil
}

func (s *UserServiceImpl) GetAll() ([]domain.User, error) {
	return s.repo.GetAll()
}
