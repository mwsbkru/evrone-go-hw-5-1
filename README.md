# Сервис управления пользователями
Разработать систему управления пользователями с интерфейсами

## Примеры использования
`cd cmd/users-manager && go run .` - запустить программу, которая последовательно вызывает все методы сервиса
с использованием каждого из репозиториев

`docker compose build && docker compose up`



mockgen -source=internal/repo/contracts.go -destination=internal/repo/contracts_mocks.go
mockgen -source=internal/usecase/user-service.go -destination=internal/usecase/user-service-mock.go