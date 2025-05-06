package app

import (
	"context"
	"evrone_go_hw_5_1/config"
	"evrone_go_hw_5_1/internal/controller/http"
	repo2 "evrone_go_hw_5_1/internal/repo"
	"evrone_go_hw_5_1/internal/usecase"
	"github.com/jackc/pgx/v5"
	"log/slog"
	"os"
)

func Run(cfg *config.Config) {
	conn, err := pgx.Connect(context.Background(), cfg.DbConnectionString)
	if err != nil {
		slog.Error("Unable to connect to database", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	repo := repo2.NewPostgreUserRepo(conn)
	userService := usecase.NewUserService(repo)
	server := http.NewHttpServer(cfg, userService)
	http.Serve(server, cfg)
}
