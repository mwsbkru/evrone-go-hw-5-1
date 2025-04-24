package usecase

import (
	"evrone_go_hw_5_1/internal/entity"
	"evrone_go_hw_5_1/internal/repo"
)

type UserService struct {
	repo repo.UserRepository
}

func NewUserService(repo repo.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (u UserService) CreateUser(name string, email string, role entity.UserRole) (entity.User, error) {
	user := entity.User{Name: name, Email: email, Role: role}
	savedUser, err := u.repo.Save(user)
	if err != nil {
		return entity.User{}, err
	}

	return savedUser, nil
}

func (u UserService) GetUser(id string) (entity.User, error) {
	user, err := u.repo.FindByID(id)
	if err != nil {
		return entity.User{}, err
	}

	return user, err
}

func (u UserService) ListUsers() ([]entity.User, error) {
	users, err := u.repo.FindAll()
	if err != nil {
		return []entity.User{}, err
	}

	return users, err
}

func (u UserService) RemoveUser(id string) error {
	err := u.repo.DeleteByID(id)
	if err != nil {
		return err
	}

	return nil
}

func (u UserService) FindByRole(role entity.UserRole) ([]entity.User, error) {
	var result []entity.User

	users, err := u.repo.FindAll()
	if err != nil {
		return result, err
	}

	for _, user := range users {
		if user.Role == role {
			result = append(result, user)
		}
	}

	return result, nil
}
