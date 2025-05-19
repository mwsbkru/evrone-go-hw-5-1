package main

import (
	"evrone_go_hw_5_1/config"
	"evrone_go_hw_5_1/internal/app"
	"log/slog"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		slog.Error("Не удалось загрузить конфигурацию приложения", slog.String("error", err.Error()))
	}

	app.Run(cfg)

	// TODO выставить golang 1.24
	// TODO Осмотреть код перед тем, как пилить тесты
	// TODO пилить тесты
	// TODO убрать у редиса авторизацию или прикрутить номер DB
}
