package usecase

import (
	"context"
	"errors"
	"evrone_go_hw_5_1/internal/entity"
	"evrone_go_hw_5_1/internal/repo"
	"fmt"
	"log/slog"
)

type UserService struct {
	repo                 repo.UserRepository
	cacheRepo            repo.UserCacheRepository
	methodCallerNotifier repo.MethodCalledNotifier
}

func NewUserService(repo repo.UserRepository, cacheRepo repo.UserCacheRepository, methodCallerNotifier repo.MethodCalledNotifier) *UserService {
	return &UserService{repo: repo, cacheRepo: cacheRepo, methodCallerNotifier: methodCallerNotifier}
}

func (u UserService) CreateUser(ctx context.Context, name string, email string, role entity.UserRole) (entity.User, error) {
	u.methodCallerNotifier.NotifyMethodCalled("CreateUser", map[string]string{
		"name":  name,
		"email": email,
		"role":  string(role),
	})

	user := entity.User{Name: name, Email: email, Role: role}
	savedUser, err := u.repo.Save(ctx, user)
	if err != nil {
		return entity.User{}, fmt.Errorf("ошибка записи пользователя в БД: %w", err)
	}

	u.cacheRepo.InvalidateAllUsersCache(ctx)
	return savedUser, nil
}

func (u UserService) GetUser(ctx context.Context, id string) (entity.User, error) {
	u.methodCallerNotifier.NotifyMethodCalled("GetUser", map[string]string{
		"id": id,
	})

	user, err := u.cacheRepo.FetchUserFromCache(ctx, id)
	if err != nil {
		if !errors.Is(err, &repo.ErrorUserNotFound{}) {
			slog.Error("не удалось получить пользователя из кеша", slog.String("error", err.Error()))
		}

		user, err = u.repo.FindByID(ctx, id)
		if err != nil {
			if errors.Is(err, &repo.ErrorUserNotFound{}) {
				return entity.User{}, err
			}

			return entity.User{}, fmt.Errorf("не удалосьполучить пользователя из БД: %w", err)
		}

		u.cacheRepo.SaveUserToCache(ctx, user)
	}

	return user, err
}

func (u UserService) ListUsers(ctx context.Context) ([]entity.User, error) {
	u.methodCallerNotifier.NotifyMethodCalled("ListUsers", map[string]string{})

	users, err := u.cacheRepo.FetchAllUsersFromCache(ctx)
	if err != nil {
		if !errors.Is(err, &repo.ErrorUserNotFound{}) {
			slog.Error("Не удалось получить пользователей из кеша", slog.String("error", err.Error()))
		}

		users, err = u.repo.FindAll(ctx)
		if err != nil {
			return []entity.User{}, fmt.Errorf("не удалось получить пользоватей из БД: %w", err)
		}

		u.cacheRepo.SaveAllUsersToCache(ctx, users)
	}

	return users, err
}

func (u UserService) RemoveUser(ctx context.Context, id string) error {
	u.methodCallerNotifier.NotifyMethodCalled("RemoveUser", map[string]string{
		"id": id,
	})

	err := u.repo.DeleteByID(ctx, id)
	if err != nil {
		return fmt.Errorf("не удалось удалить пользоватя из БД: %w", err)
	}

	u.cacheRepo.InvalidateAllUsersCache(ctx)
	u.cacheRepo.InvalidateUserInCache(ctx, id)
	return nil
}
