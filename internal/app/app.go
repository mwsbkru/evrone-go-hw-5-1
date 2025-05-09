package app

import (
	"context"
	"evrone_go_hw_5_1/config"
	"evrone_go_hw_5_1/internal/controller/http"
	repo2 "evrone_go_hw_5_1/internal/repo"
	"evrone_go_hw_5_1/internal/usecase"
	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"
	"log/slog"
	"os"
	"time"
)

func Run(cfg *config.Config) {
	ctx := context.Background()
	db := redis.NewClient(&redis.Options{
		Addr:         cfg.RedisAddr,
		Password:     cfg.RedisPassword,
		DB:           cfg.RedisDB,
		Username:     cfg.RedisUserName,
		MaxRetries:   cfg.RedisMaxRetries,
		DialTimeout:  time.Duration(cfg.RedisTimeoutSeconds) * time.Second,
		ReadTimeout:  time.Duration(cfg.RedisTimeoutSeconds) * time.Second,
		WriteTimeout: time.Duration(cfg.RedisTimeoutSeconds) * time.Second,
	})

	if err := db.Ping(ctx).Err(); err != nil {
		slog.Error("Unable to connect to redis", slog.String("error", err.Error()))
		os.Exit(1)
	}

	db.Set(ctx, "testkey", "from_go", 60*time.Second)

	conn, err := pgx.Connect(ctx, cfg.DbConnectionString)
	if err != nil {
		slog.Error("Unable to connect to database", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer conn.Close(ctx)

	repo := repo2.NewPostgreUserRepo(conn)
	userService := usecase.NewUserService(repo)
	server := http.NewHttpServer(cfg, userService)
	http.Serve(server, cfg)
}
