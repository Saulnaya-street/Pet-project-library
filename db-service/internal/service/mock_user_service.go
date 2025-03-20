package service

import (
	"github.com/google/uuid"

	"awesomeProject22/db-service/internal/domain"
)

// MockUserService временная заглушка для UserService
// В реальном проекте эту реализацию нужно заменить на настоящую
type MockUserService struct{}

func (s *MockUserService) Create(user *domain.User, password string) error {
	return nil
}

func (s *MockUserService) Authenticate(username, password string) (string, error) {
	return "mock-jwt-token", nil
}

func (s *MockUserService) GetByID(id uuid.UUID) (*domain.User, error) {
	// Для примера возвращаем фиктивного пользователя
	return &domain.User{
		ID:           id,
		Username:     "user",
		Email:        "user@example.com",
		PasswordHash: "hash",
		IsAdmin:      false,
	}, nil
}

func (s *MockUserService) IsAdmin(id uuid.UUID) (bool, error) {
	// Для примера возвращаем true если ID заканчивается на "admin"
	// В реальной системе здесь будет проверка в БД
	return id.String()[0:5] == "admin", nil
}
