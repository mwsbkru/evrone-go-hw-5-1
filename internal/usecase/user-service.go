package usecase

import (
	"errors"
	"evrone_go_hw_5_1/internal/entity"
	"evrone_go_hw_5_1/internal/repo"
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

func (u UserService) CreateUser(name string, email string, role entity.UserRole) (entity.User, error) {
	u.methodCallerNotifier.NotifyMethodCalled("CreateUser", map[string]string{
		"name":  name,
		"email": email,
		"role":  string(role),
	})

	user := entity.User{Name: name, Email: email, Role: role}
	savedUser, err := u.repo.Save(user)
	if err != nil {
		return entity.User{}, err
	}

	u.cacheRepo.InvalidateAllUsersCache()
	return savedUser, nil
}

func (u UserService) GetUser(id string) (entity.User, error) {
	u.methodCallerNotifier.NotifyMethodCalled("GetUser", map[string]string{
		"id": id,
	})

	user, err := u.cacheRepo.FetchUserFromCache(id)
	if err != nil {
		if !errors.Is(err, &repo.ErrorUserNotFound{}) {
			slog.Error("Не удалось получить пользователя из кеша", slog.String("error", err.Error()))
		}

		user, err = u.repo.FindByID(id)
		if err != nil {
			return entity.User{}, err
		}

		u.cacheRepo.SaveUserToCache(user)
	}

	return user, err
}

func (u UserService) ListUsers() ([]entity.User, error) {
	u.methodCallerNotifier.NotifyMethodCalled("ListUsers", map[string]string{})

	users, err := u.cacheRepo.FetchAllUsersFromCache()
	if err != nil {
		if !errors.Is(err, &repo.ErrorUserNotFound{}) {
			slog.Error("Не удалось получить пользователей из кеша", slog.String("error", err.Error()))
		}

		users, err = u.repo.FindAll()
		if err != nil {
			return []entity.User{}, err
		}

		u.cacheRepo.SaveAllUsersToCache(users)
	}

	return users, err
}

func (u UserService) RemoveUser(id string) error {
	u.methodCallerNotifier.NotifyMethodCalled("RemoveUser", map[string]string{
		"id": id,
	})

	err := u.repo.DeleteByID(id)
	if err != nil {
		return err
	}

	u.cacheRepo.InvalidateAllUsersCache()
	return nil
}
