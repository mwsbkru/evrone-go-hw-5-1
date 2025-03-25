package repo

import (
	"errors"
	"hw_2_1/internal/entity"
)

type InMemoryUserRepo struct {
	memo map[string]entity.User
}

func NewInMemoryUserRepo() UserRepository {
	return InMemoryUserRepo{memo: make(map[string]entity.User)}
}

func (i InMemoryUserRepo) Save(user entity.User) error {
	if user.ID == "" {
		return errors.New("User has empty ID")
	}

	i.memo[user.ID] = user
	return nil
}

func (i InMemoryUserRepo) FindByID(id string) (entity.User, error) {
	if user, ok := i.memo[id]; ok {
		return user, nil
	}

	return entity.User{}, errors.New("User not found")
}

func (i InMemoryUserRepo) FindAll() ([]entity.User, error) {
	users := make([]entity.User, 0, len(i.memo))

	for _, user := range i.memo {
		users = append(users, user)
	}

	return users, nil
}

func (i InMemoryUserRepo) DeleteByID(id string) error {
	delete(i.memo, id)

	return nil
}
