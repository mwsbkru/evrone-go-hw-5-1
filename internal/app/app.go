package app

import (
	"evrone_go_hw_5_1/config"
	"evrone_go_hw_5_1/internal/controller/http"
	repo2 "evrone_go_hw_5_1/internal/repo"
	"evrone_go_hw_5_1/internal/usecase"
)

func Run(cfg *config.Config) {
	repo := repo2.NewInMemoryUserRepo()
	userService := usecase.NewUserService(repo)
	server := http.NewHttpServer(cfg, userService)
	http.Serve(server, cfg)
}
