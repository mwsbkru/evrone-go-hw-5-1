package app

import (
	"evrone_go_hw_5_1/config"
	"evrone_go_hw_5_1/internal/controller/http"
)

func Run(cfg *config.Config) {
	server := http.NewHttpServer(cfg)
	http.Serve(server, cfg)
}
