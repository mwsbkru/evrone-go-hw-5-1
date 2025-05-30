package http_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"evrone_go_hw_5_1/config"
	http2 "evrone_go_hw_5_1/internal/controller/http"
	"evrone_go_hw_5_1/internal/entity"
	"evrone_go_hw_5_1/internal/entity/dto"
	"net/http"
	"net/http/httptest"
	"testing"
)

// MockUserService для тестирования
type MockUserService struct {
	CreateUserFunc func(ctx context.Context, name, email string, role entity.UserRole) (entity.User, error)
	GetUserFunc    func(ctx context.Context, id string) (entity.User, error)
	ListUsersFunc  func(ctx context.Context) ([]entity.User, error)
	RemoveUserFunc func(ctx context.Context, id string) error
}

func (m *MockUserService) CreateUser(ctx context.Context, name, email string, role entity.UserRole) (entity.User, error) {
	if m.CreateUserFunc != nil {
		return m.CreateUserFunc(ctx, name, email, role)
	}
	return entity.User{}, nil
}

func (m *MockUserService) GetUser(ctx context.Context, id string) (entity.User, error) {
	if m.GetUserFunc != nil {
		return m.GetUserFunc(ctx, id)
	}
	return entity.User{}, nil
}

func (m *MockUserService) ListUsers(ctx context.Context) ([]entity.User, error) {
	if m.ListUsersFunc != nil {
		return m.ListUsersFunc(ctx)
	}
	return nil, nil
}

func (m *MockUserService) RemoveUser(ctx context.Context, id string) error {
	if m.RemoveUserFunc != nil {
		return m.RemoveUserFunc(ctx, id)
	}
	return nil
}

func TestServer_Save(t *testing.T) {
	cfg := &config.Config{}

	// Успешный кейс
	t.Run("Success", func(t *testing.T) {
		mockService := &MockUserService{
			CreateUserFunc: func(ctx context.Context, name, email string, role entity.UserRole) (entity.User, error) {
				return entity.User{
					ID:    "1",
					Name:  name,
					Email: email,
					Role:  role,
				}, nil
			},
		}

		server := http2.NewServer(cfg, mockService)

		reqBody := dto.SaveUserRequestBody{
			Name:  "Test User",
			Email: "test@example.com",
			Role:  "user",
		}

		body, err := json.Marshal(reqBody)
		if err != nil {
			t.Errorf("test json marshall error: %v", err)
		}

		req := httptest.NewRequest("POST", "/users", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		server.Save(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Ожидаемый статус: %d, полученный: %d", http.StatusOK, w.Code)
		}

		var response dto.UserResponseBody
		json.Unmarshal(w.Body.Bytes(), &response)

		if response.ID != "1" || response.Name != "Test User" || response.Email != "test@example.com" || response.Role != "user" {
			t.Errorf("Некорректный ответ: %v", response)
		}
	})

	// Неверный запрос
	t.Run("InvalidRequest", func(t *testing.T) {
		mockService := &MockUserService{}
		server := http2.NewServer(cfg, mockService)

		req := httptest.NewRequest("POST", "/users", bytes.NewBuffer([]byte("{}")))
		w := httptest.NewRecorder()

		server.Save(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("want status: %d, got: %d", http.StatusBadRequest, w.Code)
		}
	})

	t.Run("BadJson", func(t *testing.T) {
		mockService := &MockUserService{}
		server := http2.NewServer(cfg, mockService)

		req := httptest.NewRequest("POST", "/users", bytes.NewBuffer([]byte("qwe")))
		w := httptest.NewRecorder()

		server.Save(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("want status: %d, got: %d", http.StatusBadRequest, w.Code)
		}
	})

	t.Run("ErrorSaveUser", func(t *testing.T) {
		mockService := &MockUserService{
			CreateUserFunc: func(ctx context.Context, name, email string, role entity.UserRole) (entity.User, error) {
				return entity.User{}, errors.New("Some error")
			},
		}

		server := http2.NewServer(cfg, mockService)

		reqBody := "{\"name\":\"name\",\"email\":\"e@mai.l\",\"role\":\"admin\"}"
		req := httptest.NewRequest("POST", "/users", bytes.NewBuffer([]byte(reqBody)))
		w := httptest.NewRecorder()

		server.Save(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("want status: %d, got: %d", http.StatusInternalServerError, w.Code)
		}
	})
}
