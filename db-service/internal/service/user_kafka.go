package service

import (
	"awesomeProject22/db-service/internal/domain"
	"awesomeProject22/db-service/internal/kafka"
	"awesomeProject22/db-service/internal/repository"
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"log"
)

type UserServiceImpl struct {
	repo          repository.IUserRepository
	eventProducer kafka.IEventProducer
}

func UserService(repo repository.IUserRepository, eventProducer kafka.IEventProducer) IUserService {
	return &UserServiceImpl{
		repo:          repo,
		eventProducer: eventProducer,
	}
}

func (s *UserServiceImpl) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error getting user by ID: %w", err)
	}
	return user, nil
}

func (s *UserServiceImpl) Create(ctx context.Context, user *domain.User, password string) error {
	if user.ID == uuid.Nil {
		user.ID = uuid.New()
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("error hashing password: %w", err)
	}

	user.PasswordHash = string(hashedPassword)

	if err := s.repo.Create(ctx, user); err != nil {
		return fmt.Errorf("error creating user: %w", err)
	}

	if err := s.eventProducer.PublishUserCreated(ctx, user); err != nil {
		log.Printf("Error publishing user creation event: %v", err)
	} else {
		log.Printf("User creation event published: %s (%s)", user.Username, user.ID)
	}

	return nil
}

func (s *UserServiceImpl) Update(ctx context.Context, user *domain.User, passwordChanged bool, newPassword string) error {

	currentUser, err := s.repo.GetByID(ctx, user.ID)
	if err != nil {
		return fmt.Errorf("user not found for update: %w", err)
	}

	if passwordChanged {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("error hashing password: %w", err)
		}
		user.PasswordHash = string(hashedPassword)
	} else {
		user.PasswordHash = currentUser.PasswordHash
	}

	if err := s.repo.Update(ctx, user); err != nil {
		return fmt.Errorf("error updating user: %w", err)
	}

	if err := s.eventProducer.PublishUserUpdated(ctx, user); err != nil {
		log.Printf("Error publishing user update event: %v", err)
	} else {
		log.Printf("User update event published: %s (%s)", user.Username, user.ID)
	}

	return nil
}

func (s *UserServiceImpl) Delete(ctx context.Context, id uuid.UUID) error {

	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("user to delete not found: %w", err)
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("error deleting user: %w", err)
	}

	if err := s.eventProducer.PublishUserDeleted(ctx, id); err != nil {
		log.Printf("Error publishing user deletion event: %v", err)
	} else {
		log.Printf("User deletion event published: %s (%s)", user.Username, id)
	}

	return nil
}

func (s *UserServiceImpl) Authenticate(ctx context.Context, username, password string) (string, error) {
	user, err := s.repo.GetByUsername(ctx, username)
	if err != nil {
		return "", errors.New("Invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", errors.New("Invalid credentials")
	}

	if err := s.eventProducer.PublishUserLoggedIn(ctx, user.ID, user.Username); err != nil {
		log.Printf("Error publishing user login event: %v", err)
	} else {
		log.Printf("User login event published: %s (%s)", user.Username, user.ID)
	}

	return user.ID.String(), nil
}

func (s *UserServiceImpl) IsAdmin(ctx context.Context, id uuid.UUID) (bool, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return false, fmt.Errorf("error checking administrator rights: %w", err)
	}

	return user.IsAdmin, nil
}

func (s *UserServiceImpl) GetAll(ctx context.Context) ([]domain.User, error) {
	users, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting list of users: %w", err)
	}
	return users, nil
}
