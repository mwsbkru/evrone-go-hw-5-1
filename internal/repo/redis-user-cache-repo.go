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

func NewRedisUserCacheRepo(client *redis.Client, cfg *config.Config) RedisUserCacheRepo {
	return RedisUserCacheRepo{client: client, cfg: cfg}
}

func (r RedisUserCacheRepo) SaveUserToCache(user entity.User) error {
	if user.ID == "" {
		return errors.New("ошибка при сериализации пользователя для кеширования: у пользователя нет ID")
	}

	cacheKey := getUserCacheKey(user.ID)
	userJSON, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("ошибка при сериализации пользователя для кеширования: %w", err)
	}

	err = r.client.Set(context.Background(), cacheKey, userJSON, time.Duration(r.cfg.CacheLifetime)*time.Second).Err()
	if err != nil {
		return fmt.Errorf("ошибка при сохранения пользователя в кеше: %w", err)
	}

	return nil
}

func (r RedisUserCacheRepo) SaveAllUsersToCache(users []entity.User) error {
	cacheKey := getUsersCacheKey()
	usersJSON, err := json.Marshal(users)
	if err != nil {
		return fmt.Errorf("ошибка при сериализации пользователей для кеширования: %w", err)
	}

	err = r.client.Set(context.Background(), cacheKey, usersJSON, time.Duration(r.cfg.CacheLifetime)*time.Second).Err()
	if err != nil {
		return fmt.Errorf("ошибка при сохранения пользователей в кеше: %w", err)
	}

	return nil
}

func (r RedisUserCacheRepo) FetchUserFromCache(id string) (entity.User, error) {
	cacheKey := getUserCacheKey(id)
	userJson, err := r.client.Get(context.Background(), cacheKey).Result()
	if err != nil {
		if err == redis.Nil {
			return entity.User{}, &ErrorUserNotFound{}
		}
		return entity.User{}, fmt.Errorf("ошибка при получении пользователя из кеша: %w", err)
	}

	var user entity.User
	err = json.Unmarshal([]byte(userJson), &user)
	if err != nil {
		return entity.User{}, fmt.Errorf("ошибка при десериализации пользователя из кеша: %w", err)
	}

	return user, nil
}

func (r RedisUserCacheRepo) FetchAllUsersFromCache() ([]entity.User, error) {
	cacheKey := getUsersCacheKey()

	usersJson, err := r.client.Get(context.Background(), cacheKey).Result()
	if err != nil {
		if err == redis.Nil {
			return []entity.User{}, &ErrorUserNotFound{}
		}
		return []entity.User{}, fmt.Errorf("ошибка при получении пользователя из кеша: %w", err)
	}

	var users []entity.User
	err = json.Unmarshal([]byte(usersJson), &users)
	if err != nil {
		return []entity.User{}, fmt.Errorf("ошибка при десериализации пользователя из кеша: %w", err)
	}

	return users, nil
}

func (r RedisUserCacheRepo) InvalidateAllUsersCache() error {
	cacheKey := getUsersCacheKey()
	err := r.client.Del(context.Background(), cacheKey).Err()
	if err != nil {
		return fmt.Errorf("ошибка при удалении кеша всех пользователей: %w", err)
	}

	return nil
}

func getUserCacheKey(userId string) string {
	return fmt.Sprintf("cache-entity.User-%s", userId)
}

func getUsersCacheKey() string {
	return USERS_CACHE_KEY
}
