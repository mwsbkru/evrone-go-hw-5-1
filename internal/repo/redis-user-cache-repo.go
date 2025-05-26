package repo

import (
	"context"
	"encoding/json"
	"errors"
	"evrone_go_hw_5_1/config"
	"evrone_go_hw_5_1/internal/entity"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

const USERS_CACHE_KEY = "cache-entity.Users"

type RedisUserCacheRepo struct {
	client *redis.Client
	cfg    *config.Config
}

func NewRedisUserCacheRepo(client *redis.Client, cfg *config.Config) *RedisUserCacheRepo {
	return &RedisUserCacheRepo{client: client, cfg: cfg}
}

func (r RedisUserCacheRepo) SaveUserToCache(ctx context.Context, user entity.User) error {
	if user.ID == "" {
		return errors.New("ошибка при сериализации пользователя для кеширования в Redis: у пользователя нет ID")
	}

	cacheKey := getUserCacheKey(user.ID)
	userJSON, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("ошибка при сериализации пользователя для кеширования в Redis: %w", err)
	}

	err = r.client.Set(ctx, cacheKey, userJSON, time.Duration(r.cfg.CacheLifetime)*time.Second).Err()
	if err != nil {
		return fmt.Errorf("ошибка при сохранения пользователя в кеше Redis: %w", err)
	}

	return nil
}

func (r RedisUserCacheRepo) SaveAllUsersToCache(ctx context.Context, users []entity.User) error {
	cacheKey := getUsersCacheKey()
	usersJSON, err := json.Marshal(users)
	if err != nil {
		return fmt.Errorf("ошибка при сериализации пользователей для кеширования в Redis: %w", err)
	}

	err = r.client.Set(ctx, cacheKey, usersJSON, time.Duration(r.cfg.CacheLifetime)*time.Second).Err()
	if err != nil {
		return fmt.Errorf("ошибка при сохранения пользователей в кеше Redis: %w", err)
	}

	return nil
}

func (r RedisUserCacheRepo) FetchUserFromCache(ctx context.Context, id string) (entity.User, error) {
	cacheKey := getUserCacheKey(id)
	userJson, err := r.client.Get(ctx, cacheKey).Result()
	if err != nil {
		if err == redis.Nil {
			return entity.User{}, &ErrorUserNotFound{}
		}
		return entity.User{}, fmt.Errorf("ошибка при получении пользователя из кеша Redis: %w", err)
	}

	var user entity.User
	err = json.Unmarshal([]byte(userJson), &user)
	if err != nil {
		return entity.User{}, fmt.Errorf("ошибка при десериализации пользователя из кеша Redis: %w", err)
	}

	return user, nil
}

func (r RedisUserCacheRepo) FetchAllUsersFromCache(ctx context.Context) ([]entity.User, error) {
	cacheKey := getUsersCacheKey()

	usersJson, err := r.client.Get(ctx, cacheKey).Result()
	if err != nil {
		if err == redis.Nil {
			return []entity.User{}, &ErrorUserNotFound{}
		}
		return []entity.User{}, fmt.Errorf("ошибка при получении пользователя из кеша Redis: %w", err)
	}

	var users []entity.User
	err = json.Unmarshal([]byte(usersJson), &users)
	if err != nil {
		return []entity.User{}, fmt.Errorf("ошибка при десериализации пользователя из кеша Redis: %w", err)
	}

	return users, nil
}

func (r RedisUserCacheRepo) InvalidateAllUsersCache(ctx context.Context) error {
	cacheKey := getUsersCacheKey()
	err := r.client.Del(ctx, cacheKey).Err()
	if err != nil {
		return fmt.Errorf("ошибка при удалении кеша Redis всех пользователей: %w", err)
	}

	return nil
}

func (r RedisUserCacheRepo) InvalidateUserInCache(ctx context.Context, id string) error {
	cacheKey := getUserCacheKey(id)

	err := r.client.Del(ctx, cacheKey).Err()
	if err != nil {
		return fmt.Errorf("ошибка при удалении кеша Redis пользователя с id %s: %w", id, err)
	}

	return nil
}

func getUserCacheKey(userId string) string {
	return fmt.Sprintf("cache-entity.User-%s", userId)
}

func getUsersCacheKey() string {
	return USERS_CACHE_KEY
}
