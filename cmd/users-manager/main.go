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

	// TODO реализовать логику сохранения в БД
	// TODO посмотреть для гошки live-reload
	// TODO прикрутить Redis
	// TODO реализовать кеширование
	// TODO прикрутить Rabbit
	// TODO реализовать отправку сообщения в Rabbit
	// TODO пилить тесты
	// TODO сделать для гуся отдельный контейнер и скрипты запуска миграций
}
