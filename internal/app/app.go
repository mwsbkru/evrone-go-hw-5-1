package app

import (
	"context"
	"evrone_go_hw_5_1/config"
	"evrone_go_hw_5_1/internal/controller/http"
	"evrone_go_hw_5_1/internal/notifier"
	repo2 "evrone_go_hw_5_1/internal/repo"
	"evrone_go_hw_5_1/internal/usecase"
	"evrone_go_hw_5_1/internal/user-cache"
	"github.com/jackc/pgx/v5"
	"github.com/nats-io/nats.go"
	"github.com/redis/go-redis/v9"
	"log/slog"
	"os"
	"time"
)

// Run initializes and runs application
func Run(ctx context.Context, cfg *config.Config) {
	redisClient := redis.NewClient(&redis.Options{
		Addr:         cfg.RedisAddr,
		DB:           cfg.RedisDB,
		MaxRetries:   cfg.RedisMaxRetries,
		DialTimeout:  time.Duration(cfg.RedisTimeoutSeconds) * time.Second,
		ReadTimeout:  time.Duration(cfg.RedisTimeoutSeconds) * time.Second,
		WriteTimeout: time.Duration(cfg.RedisTimeoutSeconds) * time.Second,
	})
	if err := redisClient.Ping(ctx).Err(); err != nil {
		slog.Error("Не удалось подключиться к redis", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer redisClient.Close()

	conn, err := pgx.Connect(ctx, cfg.DbConnectionString)
	if err != nil {
		slog.Error("не удалось подключиться к DB", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer conn.Close(ctx)

	natsConn, err := nats.Connect(cfg.NatsURL)
	if err != nil {
		slog.Warn(cfg.NatsURL)
		slog.Error("не удалось подключиться к Nats", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer natsConn.Close()

	repo := repo2.NewPostgreUserRepo(conn)
	cacheRepo := user_cache.NewRedisUserCacheRepo(redisClient, cfg)
	methodCalledNotifier := notifier.NewNatsMethodCalledNotifier(natsConn, cfg)
	userService := usecase.NewUserUseCase(repo, cacheRepo, methodCalledNotifier)
	server := http.NewServer(cfg, userService)

	http.Serve(server, cfg)
}
