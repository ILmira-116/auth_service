# Stage 1: Build
FROM golang:1.24.4-alpine AS builder

# Устанавливаем зависимости для сборки
RUN apk add --no-cache git build-base

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем весь проект
COPY . .

# Собираем бинарник основного сервиса
RUN go build -o auth-service ./cmd/auth/main.go

# Собираем бинарник для миграций
RUN go build -o migrate ./cmd/migrate/main.go

# Stage 2: Final image
FROM alpine:latest

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем бинарники из builder
COPY --from=builder /app/auth-service .
COPY --from=builder /app/migrate .

# Копируем миграции и .env
COPY ./migrations ./migrations
COPY .env .env

# По умолчанию запускаем основной сервис
CMD ["./auth-service"] 