package app

import (
	"auth-service/config"
	"auth-service/internal/app/grpcapp"
	"auth-service/internal/db"
	"auth-service/internal/repository"
	"auth-service/internal/service"

	"log/slog"
	"time"
)

type App struct {
	log     *slog.Logger
	GRPCSrv *grpcapp.App
}

func New(log *slog.Logger, grpcPort string, dbCfg *config.DBConfig, tokenTTL time.Duration, jwtSecret string) (*App, error) {
	// 1. Инициализация базы данных
	db := db.InitPostgres(dbCfg, log)

	// 2. Создание репозитория пользователей (реализует UserSaver, UserProvider, AppProvider)
	userRepo := repository.NewUserRepository(db)

	// 3. Создание сервиса аутентификации
	authSrv := service.New(
		log,
		userRepo, // UserSaver
		userRepo, // UserProvider
		userRepo, // AppProvider
		tokenTTL,
		jwtSecret,
	)
	// 4. Создание приложения с gRPC сервером
	grpcApp := grpcapp.New(log, grpcPort, authSrv)

	return &App{
		log:     log,
		GRPCSrv: grpcApp,
	}, nil
}
