package repo

import (
	"context"
	"database/sql"
	"errors"
	"evrone_go_hw_5_1/internal/entity"
	"evrone_go_hw_5_1/internal/usecase"
	"fmt"
	"github.com/jackc/pgx/v5"
)

// PostgreUserRepo provides functionality for store users in PostgreSQL
type PostgreUserRepo struct {
	conn *pgx.Conn
}

// NewPostgreUserRepo returns new PostgreSQL UserRepo
func NewPostgreUserRepo(conn *pgx.Conn) *PostgreUserRepo {
	return &PostgreUserRepo{conn: conn}
}

// Save saves user in DB
func (p PostgreUserRepo) Save(ctx context.Context, user entity.User) (entity.User, error) {
	var id string
	query := "INSERT INTO users (name, email, role) VALUES ($1, $2, $3) RETURNING id"
	err := p.conn.QueryRow(ctx, query, user.Name, user.Email, user.Role).Scan(&id)

	if err != nil {
		return user, fmt.Errorf("ошибка при вставке пользователя в DB: %w", err)
	}
	user.ID = id

	return user, nil
}

// FindByID finds user in DB by id
func (p PostgreUserRepo) FindByID(ctx context.Context, id string) (entity.User, error) {
	query := "SELECT * from users WHERE id = $1"
	var user entity.User

	err := p.conn.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Role,
		&user.CreatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return entity.User{}, &usecase.ErrUserNotFound{}
	}
	if err != nil {
		return entity.User{}, fmt.Errorf("не удалось получить пользователя из БД: %w", err)
	}

	return user, nil
}

// FindAll fetches all users from db
func (p PostgreUserRepo) FindAll(ctx context.Context) ([]entity.User, error) {
	query := "SELECT * FROM users"

	rows, err := p.conn.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении пользователей из БД: %w", err)
	}
	defer rows.Close()

	var users []entity.User

	for rows.Next() {
		var user entity.User

		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.Role,
			&user.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("ошибка при получении пользователей из БД, сканированиe строки: %w", err)
		}

		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при получении пользователей из БД, ошибка при обработке строк: %w", err)
	}

	return users, nil
}

// DeleteByID removes user with passed id from db
func (p PostgreUserRepo) DeleteByID(ctx context.Context, id string) error {
	query := "DELETE FROM users WHERE id = $1"

	_, err := p.conn.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("ошибка при удалении пользователя из БД: %w", err)
	}

	return nil
}
