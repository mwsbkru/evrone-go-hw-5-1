CMD_DIR := cmd/users-manager
GOOSE_MIGRATIONS_DIR := db/migrations

default: run

# запуск приложения
run:
	@echo "Запуск приложения..."
	@docker compose build && docker compose up

# Создание новой миграции
# пример запуска: make create-migration name=create_users
create-migration:
	@echo "Создание миграции..."
	@goose -dir $(GOOSE_MIGRATIONS_DIR) create $(name) sql

# Запуск миграций
# пример запуска: make migration-up
migration-up:
	@echo "Запуск миграций..."
	@goose -dir $(GOOSE_MIGRATIONS_DIR) postgres "postgresql://hw:hw@db.evrone-go-hw-5-1.orb.local:5432/hw?sslmode=disable" up

# Запуск миграций
# пример запуска: make migration-down
migration-down:
	@echo "Откат миграции..."
	@goose -dir $(GOOSE_MIGRATIONS_DIR) postgres "postgresql://hw:hw@db.evrone-go-hw-5-1.orb.local:5432/hw?sslmode=disable" down

# Запуск тестов
# пример запуска: make test
test:
	@echo "Запуск тестов..."
	@go test ./... -coverprofile cover.out && go tool cover -html=cover.out

# Помощь по доступным командам
help:
	@echo "Доступные команды:"
	@echo "  make run       - запуск приложения"
	@echo "  make build    - сборка приложения"
	@echo "  make test     - запуск тестов"
	@echo "  make format   - форматирование кода"
	@echo "  make clean    - очистка бинарей"
	@echo "  make help     - показать эту помощь"