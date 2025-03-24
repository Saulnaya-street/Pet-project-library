package service

import (
	"awesomeProject22/db-service/internal/domain"
	"errors"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo domain.UserRepository
}

func NewUserService(repo domain.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (s *UserService) GetByID(id uuid.UUID) (*domain.User, error) {
	return s.repo.GetByID(id)
}

func (s *UserService) Create(user *domain.User, password string) error {
	if user.ID == uuid.Nil {
		user.ID = uuid.New()
	}

	// Хешируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.PasswordHash = string(hashedPassword)

	return s.repo.Create(user)
}

func (s *UserService) Update(user *domain.User, passwordChanged bool, newPassword string) error {
	// Если пароль меняется, хешируем новый пароль
	if passwordChanged {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		user.PasswordHash = string(hashedPassword)
	}

	return s.repo.Update(user)
}

func (s *UserService) Delete(id uuid.UUID) error {
	return s.repo.Delete(id)
}

func (s *UserService) Authenticate(username, password string) (string, error) {

	user, err := s.repo.GetByUsername(username)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	// Проверяем пароль
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", errors.New("invalid credentials")
	}

	// JWT токен прикрутить когда-нибудь
	// Пока возвращаем ID пользователя как токен
	return user.ID.String(), nil
}

func (s *UserService) IsAdmin(id uuid.UUID) (bool, error) {
	user, err := s.repo.GetByID(id)
	if err != nil {
		return false, err
	}

	return user.IsAdmin, nil
}
func (s *UserService) GetAll() ([]*domain.User, error) {
	return s.repo.GetAll()
}
