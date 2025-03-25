package repo

import "hw_2_1/internal/entity"

type UserRepository interface {
	Save(user entity.User) error
	FindByID(id string) (entity.User, error)
	FindAll() ([]entity.User, error)
	DeleteByID(id string) error
}
