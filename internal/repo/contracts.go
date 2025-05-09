package repo

import "evrone_go_hw_5_1/internal/entity"

type UserRepository interface {
	Save(user entity.User) (entity.User, error)
	FindByID(id string) (entity.User, error)
	FindAll() ([]entity.User, error)
	DeleteByID(id string) error
}

type ErrorUserNotFound struct{}

func (e *ErrorUserNotFound) Error() string {
	return "Not Found"
}

type UserCacheRepository interface {
	SaveUserToCache(entity.User) error
	FetchUserFromCache(id string) (entity.User, error)
}
