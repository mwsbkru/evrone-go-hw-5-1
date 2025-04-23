package repo

import (
	"evrone_go_hw_5_1/internal/entity"
	"log/slog"
)

type MockUserRepo struct{}

func NewMockUserRepo() UserRepository {
	return MockUserRepo{}
}

func (m MockUserRepo) Save(user entity.User) error {
	slog.Warn("Save user: ", user)
	return nil
}

func (m MockUserRepo) FindByID(id string) (entity.User, error) {
	slog.Warn("FindByID user with id: ", id)
	return entity.User{ID: id}, nil
}

func (m MockUserRepo) FindAll() ([]entity.User, error) {
	slog.Warn("FindAll users")
	return []entity.User{}, nil
}

func (m MockUserRepo) DeleteByID(id string) error {
	slog.Warn("DeleteByID user with id: ", id)
	return nil
}
