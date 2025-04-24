package repo

import "evrone_go_hw_5_1/internal/entity"

type UserRepository interface {
	Save(user entity.User) (entity.User, error)
	FindByID(id string) (entity.User, error)
	FindAll() ([]entity.User, error)
	DeleteByID(id string) error
}
