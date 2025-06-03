package main

import (
	"context"
	"evrone_go_hw_5_1/config"
	"evrone_go_hw_5_1/internal/app"
	"log/slog"
	"os/signal"
	"syscall"
)

func main() {
	cfg, err := config.NewConfig()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	if err != nil {
		slog.Error("Не удалось загрузить конфигурацию приложения", slog.String("error", err.Error()))
	}

	app.Run(cfg, ctx)

	// TODO привести в порядок makefile и readme
}
