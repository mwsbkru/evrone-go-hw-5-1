package http

import (
	"context"
	"evrone_go_hw_5_1/internal/entity"
)

type UserService interface {
	CreateUser(ctx context.Context, name string, email string, role entity.UserRole) (entity.User, error)
	GetUser(ctx context.Context, id string) (entity.User, error)
	ListUsers(ctx context.Context) ([]entity.User, error)
	RemoveUser(ctx context.Context, id string) error
}
