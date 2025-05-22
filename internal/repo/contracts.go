package repo

import (
	"context"
	"evrone_go_hw_5_1/internal/entity"
)

type UserRepository interface {
	Save(ctx context.Context, user entity.User) (entity.User, error)
	FindByID(ctx context.Context, id string) (entity.User, error)
	FindAll(ctx context.Context) ([]entity.User, error)
	DeleteByID(ctx context.Context, id string) error
}

type ErrorUserNotFound struct{}

func (e *ErrorUserNotFound) Error() string {
	return "Not Found"
}

type UserCacheRepository interface {
	SaveUserToCache(ctx context.Context, user entity.User) error
	FetchUserFromCache(ctx context.Context, id string) (entity.User, error)
	InvalidateUserInCache(ctx context.Context, id string) error
	SaveAllUsersToCache(ctx context.Context, user []entity.User) error
	FetchAllUsersFromCache(ctx context.Context) ([]entity.User, error)
	InvalidateAllUsersCache(ctx context.Context) error
}

type MethodCalledNotifier interface {
	NotifyMethodCalled(methodName string, params map[string]string) error
}
