package usecase_test

import (
	"context"
	"errors"
	"evrone_go_hw_5_1/internal/entity"
	"evrone_go_hw_5_1/internal/repo"
	"evrone_go_hw_5_1/internal/usecase"
	"github.com/golang/mock/gomock"
	"testing"
)

func TestUserService_CreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repo.NewMockUserRepository(ctrl)
	mockCache := repo.NewMockUserCacheRepository(ctrl)
	mockNotifier := repo.NewMockMethodCalledNotifier(ctrl)

	userService := usecase.NewUserUseCase(mockRepo, mockCache, mockNotifier)

	// Успешный кейс создания пользователя
	mockNotifier.EXPECT().NotifyMethodCalled("CreateUser", map[string]string{
		"name": "Test User", "email": "test@example.com", "role": "admin",
	})
	mockRepo.EXPECT().Save(context.Background(), entity.User{
		Name: "Test User", Email: "test@example.com", Role: entity.UserRoleAdmin,
	}).Return(entity.User{
		ID: "123", Name: "Test User", Email: "test@example.com", Role: entity.UserRoleAdmin,
	}, nil)
	mockCache.EXPECT().InvalidateAllUsersCache(context.Background())

	user, err := userService.CreateUser(context.Background(), "Test User", "test@example.com", entity.UserRoleAdmin)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if user.ID != "123" {
		t.Errorf("Expected user ID '123', got '%s'", user.ID)
	}

	// Ошибка при сохранении пользователя
	mockNotifier.EXPECT().NotifyMethodCalled("CreateUser", map[string]string{
		"name": "Test User 2", "email": "test2@example.com", "role": "user",
	})
	mockRepo.EXPECT().Save(context.Background(), entity.User{
		Name: "Test User 2", Email: "test2@example.com", Role: entity.UserRoleUser,
	}).Return(entity.User{}, errors.New("save error"))

	_, err = userService.CreateUser(context.Background(), "Test User 2", "test2@example.com", entity.UserRoleUser)
	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestUserService_GetUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repo.NewMockUserRepository(ctrl)
	mockCache := repo.NewMockUserCacheRepository(ctrl)
	mockNotifier := repo.NewMockMethodCalledNotifier(ctrl)

	userService := usecase.NewUserUseCase(mockRepo, mockCache, mockNotifier)

	// Случай 1: Пользователь найден в кеше
	mockNotifier.EXPECT().NotifyMethodCalled("GetUser", map[string]string{
		"id": "123",
	})
	mockCache.EXPECT().FetchUserFromCache(context.Background(), "123").Return(entity.User{
		ID: "123", Name: "Cached User", Email: "cached@example.com", Role: entity.UserRoleUser,
	}, nil)

	user, err := userService.GetUser(context.Background(), "123")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if user.Name != "Cached User" {
		t.Errorf("Expected user name 'Cached User', got '%s'", user.Name)
	}

	// Случай 2: Пользователь не найден в кеше, но найден в БД
	mockNotifier.EXPECT().NotifyMethodCalled("GetUser", map[string]string{
		"id": "456",
	})
	mockCache.EXPECT().FetchUserFromCache(context.Background(), "456").Return(entity.User{}, &usecase.ErrUserNotFound{})
	mockRepo.EXPECT().FindByID(context.Background(), "456").Return(entity.User{
		ID: "456", Name: "DB User", Email: "db@example.com", Role: entity.UserRoleAdmin,
	}, nil)
	mockCache.EXPECT().SaveUserToCache(context.Background(), entity.User{
		ID: "456", Name: "DB User", Email: "db@example.com", Role: entity.UserRoleAdmin,
	})

	user, err = userService.GetUser(context.Background(), "456")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if user.Name != "DB User" {
		t.Errorf("Expected user name 'DB User', got '%s'", user.Name)
	}

	// Случай 3: Пользователь не найден ни в кеше, ни в БД
	mockNotifier.EXPECT().NotifyMethodCalled("GetUser", map[string]string{
		"id": "789",
	})
	mockCache.EXPECT().FetchUserFromCache(context.Background(), "789").Return(entity.User{}, &usecase.ErrUserNotFound{})
	mockRepo.EXPECT().FindByID(context.Background(), "789").Return(entity.User{}, &usecase.ErrUserNotFound{})

	_, err = userService.GetUser(context.Background(), "789")
	if err == nil {
		t.Error("Expected ErrUserNotFound, got nil")
	}

	// Случай 4: Ошибка при получении из кеша
	mockNotifier.EXPECT().NotifyMethodCalled("GetUser", map[string]string{
		"id": "000",
	})
	mockCache.EXPECT().FetchUserFromCache(context.Background(), "000").Return(entity.User{}, errors.New("cache error"))
	mockRepo.EXPECT().FindByID(context.Background(), "000").Return(entity.User{
		ID: "000", Name: "DB User 2", Email: "db2@example.com", Role: entity.UserRoleUser,
	}, nil)

	mockCache.EXPECT().SaveUserToCache(context.Background(), entity.User{
		ID: "000", Name: "DB User 2", Email: "db2@example.com", Role: entity.UserRoleUser,
	}).Return(nil)

	user, err = userService.GetUser(context.Background(), "000")
	if err != nil {
		t.Error("Expected User, got error")
	}

	if user.Name != "DB User 2" {
		t.Errorf("Expected user name 'DB User 2', got '%s'", user.Name)
	}

	// Случай 5: Пользователь не найден в кеше, ошибка при получении пользователя из БД
	mockNotifier.EXPECT().NotifyMethodCalled("GetUser", map[string]string{
		"id": "789",
	})
	mockCache.EXPECT().FetchUserFromCache(context.Background(), "789").Return(entity.User{}, &usecase.ErrUserNotFound{})
	mockRepo.EXPECT().FindByID(context.Background(), "789").Return(entity.User{}, errors.New("db error"))

	_, err = userService.GetUser(context.Background(), "789")
	if err == nil {
		t.Error("Expected db error, got nil")
	}
}

func TestUserService_ListUsers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repo.NewMockUserRepository(ctrl)
	mockCache := repo.NewMockUserCacheRepository(ctrl)
	mockNotifier := repo.NewMockMethodCalledNotifier(ctrl)

	userService := usecase.NewUserUseCase(mockRepo, mockCache, mockNotifier)

	// Случай 1: Пользователи найдены в кеше
	mockNotifier.EXPECT().NotifyMethodCalled("ListUsers", map[string]string{})
	mockCache.EXPECT().FetchAllUsersFromCache(context.Background()).Return([]entity.User{
		{ID: "123", Name: "User 1", Email: "user1@example.com", Role: entity.UserRoleUser},
		{ID: "456", Name: "User 2", Email: "user2@example.com", Role: entity.UserRoleAdmin},
	}, nil)

	users, err := userService.ListUsers(context.Background())
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if len(users) != 2 {
		t.Errorf("Expected 2 users, got %d", len(users))
	}

	// Случай 2: Пользователи не найдены в кеше, но найдены в БД
	mockNotifier.EXPECT().NotifyMethodCalled("ListUsers", map[string]string{})
	mockCache.EXPECT().FetchAllUsersFromCache(context.Background()).Return([]entity.User{}, errors.New("cache error"))
	mockRepo.EXPECT().FindAll(context.Background()).Return([]entity.User{
		{ID: "789", Name: "User 3", Email: "user3@example.com", Role: entity.UserRoleUser},
		{ID: "101", Name: "User 4", Email: "user4@example.com", Role: entity.UserRoleAdmin},
	}, nil)
	mockCache.EXPECT().SaveAllUsersToCache(context.Background(), []entity.User{
		{ID: "789", Name: "User 3", Email: "user3@example.com", Role: entity.UserRoleUser},
		{ID: "101", Name: "User 4", Email: "user4@example.com", Role: entity.UserRoleAdmin},
	})

	users, err = userService.ListUsers(context.Background())
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if len(users) != 2 {
		t.Errorf("Expected 2 users, got %d", len(users))
	}

	// Случай 3: Ошибка при получении пользователей из БД
	mockNotifier.EXPECT().NotifyMethodCalled("ListUsers", map[string]string{})
	mockCache.EXPECT().FetchAllUsersFromCache(context.Background()).Return([]entity.User{}, &usecase.ErrUserNotFound{})
	mockRepo.EXPECT().FindAll(context.Background()).Return([]entity.User{}, errors.New("database error"))

	_, err = userService.ListUsers(context.Background())
	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestUserService_RemoveUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repo.NewMockUserRepository(ctrl)
	mockCache := repo.NewMockUserCacheRepository(ctrl)
	mockNotifier := repo.NewMockMethodCalledNotifier(ctrl)

	userService := usecase.NewUserUseCase(mockRepo, mockCache, mockNotifier)

	// Успешное удаление пользователя
	mockNotifier.EXPECT().NotifyMethodCalled("RemoveUser", map[string]string{
		"id": "123",
	})
	mockRepo.EXPECT().DeleteByID(context.Background(), "123").Return(nil)
	mockCache.EXPECT().InvalidateAllUsersCache(context.Background())
	mockCache.EXPECT().InvalidateUserInCache(context.Background(), "123")

	err := userService.RemoveUser(context.Background(), "123")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Ошибка при удалении из БД
	mockNotifier.EXPECT().NotifyMethodCalled("RemoveUser", map[string]string{
		"id": "456",
	})
	mockRepo.EXPECT().DeleteByID(context.Background(), "456").Return(errors.New("delete error"))

	err = userService.RemoveUser(context.Background(), "456")
	if err == nil {
		t.Error("Expected error, got nil")
	}

	// Проверка последовательности вызовов методов кеша
	mockNotifier.EXPECT().NotifyMethodCalled("RemoveUser", map[string]string{
		"id": "789",
	})
	mockRepo.EXPECT().DeleteByID(context.Background(), "789").Return(nil)
	mockCache.EXPECT().InvalidateAllUsersCache(context.Background())
	mockCache.EXPECT().InvalidateUserInCache(context.Background(), "789")

	err = userService.RemoveUser(context.Background(), "789")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}
