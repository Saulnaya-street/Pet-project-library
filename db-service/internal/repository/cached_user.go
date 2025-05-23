package repository

import (
	"awesomeProject22/db-service/internal/cache"
	"awesomeProject22/db-service/internal/domain"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
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
	redisClient cache.IRedisClient
}

func CachedUserRepo(repo IUserRepository, redisClient cache.IRedisClient) IUserRepository {
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

func (r *CachedUserRepository) Create(ctx context.Context, user *domain.User) error {
	err := r.repo.Create(ctx, user)
	if err != nil {
		return fmt.Errorf("error creating user in database: %w", err)
	}

	userJson, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("error serializing user: %w", err)
	}

	err = r.redisClient.Set(ctx, getUserKey(user.ID), string(userJson), userCacheTTL)
	if err != nil {
		fmt.Printf("Error caching user by ID: %v\n", err)
	}

	err = r.redisClient.Set(ctx, getUsernameKey(user.Username), user.ID.String(), userCacheTTL)
	if err != nil {
		fmt.Printf("Error caching user by name: %v\n", err)
	}

	err = r.redisClient.Set(ctx, getEmailKey(user.Email), user.ID.String(), userCacheTTL)
	if err != nil {
		fmt.Printf("Error caching user by email: %v\n", err)
	}

	r.redisClient.Delete(ctx, userListKey)

	return nil
}

func (r *CachedUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	userKey := getUserKey(id)
	cachedUser, err := r.redisClient.Get(ctx, userKey)

	if err == nil {
		var user domain.User
		if err := json.Unmarshal([]byte(cachedUser), &user); err == nil {
			return &user, nil
		} else {
			fmt.Printf("Error unmarshaling cached user: %v\n", err)
		}
	} else if err != redis.Nil {
		fmt.Printf("Error fetching user from Redis: %v\n", err)
	}

	user, err := r.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error getting user from database: %w", err)
	}

	userJson, err := json.Marshal(user)
	if err == nil {
		r.redisClient.Set(ctx, userKey, string(userJson), userCacheTTL)
		r.redisClient.Set(ctx, getUsernameKey(user.Username), user.ID.String(), userCacheTTL)
		r.redisClient.Set(ctx, getEmailKey(user.Email), user.ID.String(), userCacheTTL)
	}

	return user, nil
}

func (r *CachedUserRepository) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	usernameKey := getUsernameKey(username)
	userID, err := r.redisClient.Get(ctx, usernameKey)

	if err == nil {
		id, err := uuid.Parse(userID)
		if err == nil {
			return r.GetByID(ctx, id)
		}
	} else if err != redis.Nil {
		fmt.Printf("Error fetching user by username from Redis: %v\n", err)
	}

	user, err := r.repo.GetByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("error getting user by username from database: %w", err)
	}

	userJson, err := json.Marshal(user)
	if err == nil {
		r.redisClient.Set(ctx, getUserKey(user.ID), string(userJson), userCacheTTL)
		r.redisClient.Set(ctx, usernameKey, user.ID.String(), userCacheTTL)
		r.redisClient.Set(ctx, getEmailKey(user.Email), user.ID.String(), userCacheTTL)
	}

	return user, nil
}

func (r *CachedUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	emailKey := getEmailKey(email)
	userID, err := r.redisClient.Get(ctx, emailKey)

	if err == nil {
		id, err := uuid.Parse(userID)
		if err == nil {
			return r.GetByID(ctx, id)
		}
	} else if err != redis.Nil {
		fmt.Printf("Error fetching user by email from Redis: %v\n", err)
	}

	user, err := r.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("error getting user by email from database: %w", err)
	}

	userJson, err := json.Marshal(user)
	if err == nil {
		r.redisClient.Set(ctx, getUserKey(user.ID), string(userJson), userCacheTTL)
		r.redisClient.Set(ctx, getUsernameKey(user.Username), user.ID.String(), userCacheTTL)
		r.redisClient.Set(ctx, emailKey, user.ID.String(), userCacheTTL)
	}

	return user, nil
}

func (r *CachedUserRepository) Update(ctx context.Context, user *domain.User) error {
	oldUser, err := r.repo.GetByID(ctx, user.ID)
	if err == nil {
		if oldUser.Username != user.Username {
			r.redisClient.Delete(ctx, getUsernameKey(oldUser.Username))
		}
		if oldUser.Email != user.Email {
			r.redisClient.Delete(ctx, getEmailKey(oldUser.Email))
		}
	}

	err = r.repo.Update(ctx, user)
	if err != nil {
		return fmt.Errorf("error updating user in database: %w", err)
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

func (r *CachedUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	user, err := r.repo.GetByID(ctx, id)
	if err == nil {
		r.redisClient.Delete(ctx, getUsernameKey(user.Username))
		r.redisClient.Delete(ctx, getEmailKey(user.Email))
	}

	err = r.repo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("error deleting user from database: %w", err)
	}

	r.redisClient.Delete(ctx, getUserKey(id))
	r.redisClient.Delete(ctx, userListKey)

	return nil
}

func (r *CachedUserRepository) GetAll(ctx context.Context) ([]domain.User, error) {
	cachedList, err := r.redisClient.Get(ctx, userListKey)

	if err == nil {
		var users []domain.User
		if unmarshalErr := json.Unmarshal([]byte(cachedList), &users); unmarshalErr == nil {
			return users, nil
		} else {
			fmt.Printf("Error deserializing user list from cache: %v\n", unmarshalErr)
		}
	} else if err != redis.Nil {
		fmt.Printf("Error fetching users list from Redis: %v\n", err)
	}

	users, err := r.repo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting list of users from database: %w", err)
	}

	usersJson, err := json.Marshal(users)
	if err == nil {
		r.redisClient.Set(ctx, userListKey, string(usersJson), userCacheTTL)
	}

	return users, nil
}
