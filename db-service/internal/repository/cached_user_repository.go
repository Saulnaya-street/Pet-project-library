package repository

import (
	"awesomeProject22/db-service/internal/cache"
	"awesomeProject22/db-service/internal/domain"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"time"
)

const (
	userKeyPrefix        = "user:"
	userByUsernamePrefix = "user:username:"
	userByEmailPrefix    = "user:email:"
	userListKey          = "users:all"
	userCacheTTL         = 30 * time.Minute
)

type CachedUserRepository struct {
	repo        IUserRepository
	redisClient *cache.RedisClient
}

func NewCachedUserRepository(repo IUserRepository, redisClient *cache.RedisClient) IUserRepository {
	return &CachedUserRepository{
		repo:        repo,
		redisClient: redisClient,
	}
}

func getUserKey(id uuid.UUID) string {
	return fmt.Sprintf("%s%s", userKeyPrefix, id.String())
}

func getUsernameKey(username string) string {
	return fmt.Sprintf("%s%s", userByUsernamePrefix, username)
}

func getEmailKey(email string) string {
	return fmt.Sprintf("%s%s", userByEmailPrefix, email)
}

func (r *CachedUserRepository) Create(user *domain.User) error {

	err := r.repo.Create(user)
	if err != nil {
		return err
	}

	userJson, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("ошибка при сериализации пользователя: %w", err)
	}

	ctx := context.Background()

	err = r.redisClient.Set(ctx, getUserKey(user.ID), string(userJson), userCacheTTL)
	if err != nil {
		fmt.Printf("Ошибка при кешировании пользователя по ID: %v\n", err)
	}

	err = r.redisClient.Set(ctx, getUsernameKey(user.Username), user.ID.String(), userCacheTTL)
	if err != nil {
		fmt.Printf("Ошибка при кешировании пользователя по имени: %v\n", err)
	}

	err = r.redisClient.Set(ctx, getEmailKey(user.Email), user.ID.String(), userCacheTTL)
	if err != nil {
		fmt.Printf("Ошибка при кешировании пользователя по email: %v\n", err)
	}

	r.redisClient.Delete(ctx, userListKey)

	return nil
}

func (r *CachedUserRepository) GetByID(id uuid.UUID) (*domain.User, error) {
	ctx := context.Background()

	userKey := getUserKey(id)
	cachedUser, err := r.redisClient.Get(ctx, userKey)
	if err == nil && cachedUser != "" {

		var user domain.User
		err = json.Unmarshal([]byte(cachedUser), &user)
		if err == nil {
			return &user, nil
		}
	}

	user, err := r.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	userJson, err := json.Marshal(user)
	if err == nil {
		r.redisClient.Set(ctx, userKey, string(userJson), userCacheTTL)
		r.redisClient.Set(ctx, getUsernameKey(user.Username), user.ID.String(), userCacheTTL)
		r.redisClient.Set(ctx, getEmailKey(user.Email), user.ID.String(), userCacheTTL)
	}

	return user, nil
}

func (r *CachedUserRepository) GetByUsername(username string) (*domain.User, error) {
	ctx := context.Background()

	usernameKey := getUsernameKey(username)
	userID, err := r.redisClient.Get(ctx, usernameKey)
	if err == nil && userID != "" {

		id, err := uuid.Parse(userID)
		if err == nil {

			return r.GetByID(id)
		}
	}

	user, err := r.repo.GetByUsername(username)
	if err != nil {
		return nil, err
	}

	userJson, err := json.Marshal(user)
	if err == nil {
		r.redisClient.Set(ctx, getUserKey(user.ID), string(userJson), userCacheTTL)
		r.redisClient.Set(ctx, usernameKey, user.ID.String(), userCacheTTL)
		r.redisClient.Set(ctx, getEmailKey(user.Email), user.ID.String(), userCacheTTL)
	}

	return user, nil
}

func (r *CachedUserRepository) GetByEmail(email string) (*domain.User, error) {
	ctx := context.Background()

	emailKey := getEmailKey(email)
	userID, err := r.redisClient.Get(ctx, emailKey)
	if err == nil && userID != "" {

		id, err := uuid.Parse(userID)
		if err == nil {
			return r.GetByID(id)
		}
	}

	user, err := r.repo.GetByEmail(email)
	if err != nil {
		return nil, err
	}

	userJson, err := json.Marshal(user)
	if err == nil {
		r.redisClient.Set(ctx, getUserKey(user.ID), string(userJson), userCacheTTL)
		r.redisClient.Set(ctx, getUsernameKey(user.Username), user.ID.String(), userCacheTTL)
		r.redisClient.Set(ctx, emailKey, user.ID.String(), userCacheTTL)
	}

	return user, nil
}

func (r *CachedUserRepository) Update(user *domain.User) error {
	ctx := context.Background()

	oldUser, err := r.repo.GetByID(user.ID)
	if err == nil {

		if oldUser.Username != user.Username {
			r.redisClient.Delete(ctx, getUsernameKey(oldUser.Username))
		}
		if oldUser.Email != user.Email {
			r.redisClient.Delete(ctx, getEmailKey(oldUser.Email))
		}
	}

	err = r.repo.Update(user)
	if err != nil {
		return err
	}

	userJson, err := json.Marshal(user)
	if err == nil {
		r.redisClient.Set(ctx, getUserKey(user.ID), string(userJson), userCacheTTL)
		r.redisClient.Set(ctx, getUsernameKey(user.Username), user.ID.String(), userCacheTTL)
		r.redisClient.Set(ctx, getEmailKey(user.Email), user.ID.String(), userCacheTTL)
	}

	r.redisClient.Delete(ctx, userListKey)

	return nil
}

func (r *CachedUserRepository) Delete(id uuid.UUID) error {
	ctx := context.Background()

	user, err := r.repo.GetByID(id)
	if err == nil {

		r.redisClient.Delete(ctx, getUsernameKey(user.Username))
		r.redisClient.Delete(ctx, getEmailKey(user.Email))
	}

	err = r.repo.Delete(id)
	if err != nil {
		return err
	}

	r.redisClient.Delete(ctx, getUserKey(id))

	r.redisClient.Delete(ctx, userListKey)

	return nil
}

func (r *CachedUserRepository) GetAll() ([]domain.User, error) {
	ctx := context.Background()

	cachedList, err := r.redisClient.Get(ctx, userListKey)
	if err == nil && cachedList != "" {

		var users []domain.User
		err = json.Unmarshal([]byte(cachedList), &users)
		if err == nil {
			return users, nil
		}
	}

	users, err := r.repo.GetAll()
	if err != nil {
		return nil, err
	}

	usersJson, err := json.Marshal(users)
	if err == nil {
		r.redisClient.Set(ctx, userListKey, string(usersJson), userCacheTTL)
	}

	return users, nil
}
