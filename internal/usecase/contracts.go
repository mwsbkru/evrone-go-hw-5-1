package usecase

import (
	"context"
	"evrone_go_hw_5_1/internal/entity"
)

// UserRepository describes functionality for operating users in DB
type UserRepository interface {
	Save(ctx context.Context, user entity.User) (entity.User, error)
	FindByID(ctx context.Context, id string) (entity.User, error)
	FindAll(ctx context.Context) ([]entity.User, error)
	DeleteByID(ctx context.Context, id string) error
}

// UserCacheRepository describes functionality for operating users in cache
type UserCacheRepository interface {
	SaveUserToCache(ctx context.Context, user entity.User) error
	FetchUserFromCache(ctx context.Context, id string) (entity.User, error)
	InvalidateUserInCache(ctx context.Context, id string) error
	SaveAllUsersToCache(ctx context.Context, user []entity.User) error
	FetchAllUsersFromCache(ctx context.Context) ([]entity.User, error)
	InvalidateAllUsersCache(ctx context.Context) error
}

// MethodCalledNotifier describes functionality for notifications about service actions called
type MethodCalledNotifier interface {
	NotifyMethodCalled(methodName string, params map[string]string) error
}

// ErrUserNotFound specific error for cases when user not found in repo
type ErrUserNotFound struct{}

func (e *ErrUserNotFound) Error() string {
	return "Not Found"
}
