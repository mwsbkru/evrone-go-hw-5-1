package repo

import (
	"evrone_go_hw_5_1/internal/entity"
	"strconv"
)

type InMemoryUserRepo struct {
	memo   map[string]entity.User
	lastId int
}

func NewInMemoryUserRepo() *InMemoryUserRepo {
	return &InMemoryUserRepo{memo: make(map[string]entity.User)}
}

func (i *InMemoryUserRepo) Save(user entity.User) (entity.User, error) {
	i.lastId++
	user.ID = strconv.Itoa(i.lastId)
	i.memo[user.ID] = user
	return user, nil
}

func (i *InMemoryUserRepo) FindByID(id string) (entity.User, error) {
	if user, ok := i.memo[id]; ok {
		return user, nil
	}

	return entity.User{}, &ErrorUserNotFound{}
}

func (i *InMemoryUserRepo) FindAll() ([]entity.User, error) {
	users := make([]entity.User, 0, len(i.memo))

	for _, user := range i.memo {
		users = append(users, user)
	}

	return users, nil
}

func (i *InMemoryUserRepo) DeleteByID(id string) error {
	if _, ok := i.memo[id]; ok {
		delete(i.memo, id)
		return nil
	}

	return &ErrorUserNotFound{}
}
