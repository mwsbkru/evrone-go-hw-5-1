package repo

import (
	"context"
	"database/sql"
	"errors"
	"evrone_go_hw_5_1/internal/entity"
	"fmt"
	"github.com/jackc/pgx/v5"
	"time"
)

type PostgreUserRepo struct {
	conn *pgx.Conn
}

func NewPostgreUserRepo(conn *pgx.Conn) *PostgreUserRepo {
	return &PostgreUserRepo{conn: conn}
}

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

func (p PostgreUserRepo) FindByID(ctx context.Context, id string) (entity.User, error) {
	// типа мы долго ищем юзера
	time.Sleep(time.Second)

	query := "SELECT * from users WHERE id = $1"
	var user entity.User

	// Выполняем запрос и сканируем результат
	err := p.conn.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Role,
		&user.CreatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return entity.User{}, &ErrorUserNotFound{}
	}
	if err != nil {
		return entity.User{}, err
	}

	return user, nil
}

func (p PostgreUserRepo) FindAll(ctx context.Context) ([]entity.User, error) {
	// типа мы долго ищем юзера
	time.Sleep(time.Second)

	query := "SELECT * FROM users"

	rows, err := p.conn.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении пользователей: %w", err)
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
			return nil, fmt.Errorf("ошибка при получении пользователей, сканированиe строки: %w", err)
		}

		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при получении пользователей, ошибка при обработке строк: %w", err)
	}

	return users, nil
}

func (p PostgreUserRepo) DeleteByID(ctx context.Context, id string) error {
	query := "DELETE FROM users WHERE id = $1"

	_, err := p.conn.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("ошибка при удалении пользователя: %w", err)
	}

	return nil
}
