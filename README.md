 # Auth Service

Сервис аутентификации и авторизации на Go с использованием gRPC, ProtoBuf и JWT.

---

## Основные RPC методы

Сервис `Auth` предоставляет следующие методы:

| RPC         | Request Type       | Response Type      | Описание |
|------------|------------------|-----------------|----------|
| `Login`    | `LoginRequest`    | `LoginResponse`  | Аутентификация пользователя. При успешной аутентификации возвращается JWT токен. Параметры: `email`, `password`, `app_id`. |
| `Register` | `RegisterRequest` | `RegisterResponse` | Регистрация нового пользователя. При успешной регистрации возвращается `user_id`. Параметры: `email`, `password`. |
| `IsAdmin`  | `IsAdminRequest`  | `IsAdminResponse` | Проверка, является ли пользователь администратором. Параметр: `user_id`. |

---

## Технологии и зависимости

- Go 1.24.4  
- gRPC (`google.golang.org/grpc`)  
- Protocol Buffers (`github.com/ILmira-116/protos/gen/auth`)  
- JWT (`github.com/golang-jwt/jwt/v5`)  
- PostgreSQL (`pgx`, `pq`)  
- Миграции базы: `goose`  
- Конфигурации: `cleanenv`
- Логирование: log/slog
- Корректное завершение работы сервиса: graceful shutdown
- Генерация случайных данных для тестов: `gofakeit`  
- Тестирование: `testify`  

---

## Protobuf Contract

- Исходные файлы `.proto` находятся в модуле:  
  `github.com/ILmira-116/protos`

- Сгенерированные Go пакеты доступны по пути:  
  `github.com/ILmira-116/protos/gen/auth`

---

## Запуск сервиса

Сервис можно запустить локально или через Docker Compose.

### Локальный запуск

```bash
git clone git@github.com:ILmira-116/auth_service.git
cd auth-service
go mod tidy
go run cmd/main.go

### Запуск через Docker Compose

```bash
Копировать код
docker-compose up --build

Сервис будет доступен на порту 50051.

### Планы по развитию

В будущем планируется интеграция с сервисом привычек для проверки, является ли пользователь администратором.
На данный момент сервис привычек — это CRUD-сервис для управления привычками пользователей, который находится в разработке. 
