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
	redisClient := redis.NewClient(&redis.Options{
		Addr:         cfg.RedisAddr,
		Password:     cfg.RedisPassword,
		DB:           cfg.RedisDB,
		Username:     cfg.RedisUserName,
		MaxRetries:   cfg.RedisMaxRetries,
		DialTimeout:  time.Duration(cfg.RedisTimeoutSeconds) * time.Second,
		ReadTimeout:  time.Duration(cfg.RedisTimeoutSeconds) * time.Second,
		WriteTimeout: time.Duration(cfg.RedisTimeoutSeconds) * time.Second,
	})

	if err := redisClient.Ping(ctx).Err(); err != nil {
		slog.Error("Не удалось подключиться к redis", slog.String("error", err.Error()))
		os.Exit(1)
	}

	conn, err := pgx.Connect(ctx, cfg.DbConnectionString)
	if err != nil {
		slog.Error("не удалось подключиться к DB", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer conn.Close(ctx)

	repo := repo2.NewPostgreUserRepo(conn)
	cacheRepo := repo2.NewRedisUserCacheRepo(redisClient, cfg)
	userService := usecase.NewUserService(repo, cacheRepo)
	server := http.NewHttpServer(cfg, userService)
	http.Serve(server, cfg)
}
