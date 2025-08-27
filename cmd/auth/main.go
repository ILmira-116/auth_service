package main

import (
	"auth-service/config"
	"auth-service/internal/app"
	"auth-service/internal/logger"
	"auth-service/internal/shutdown"
	"fmt"
	"log/slog"
	"time"
)

func main() {
	// 1. Инициализация конфига
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(fmt.Sprintf("failed to load config: %v", err))
	}

	// 2. Инициализация логгера
	log := logger.New(cfg)
	log.Info("Logger is ready")
	log.Debug("Debug message")

	// 3. Приложение
	application, err := app.New(log, cfg.GRPC.ServerPort, &cfg.DB, cfg.TokenTTL, cfg.JWTSecret)
	if err != nil {
		log.Error("failed to create application", slog.String("err", err.Error()))
		return
	}

	// 4. Запуск gRPC сервера в отдельной горутине
	go func() {
		if err := application.GRPCSrv.Run(); err != nil {
			log.Error("grpc server failed", slog.String("err", err.Error()))
		}
	}()

	// 5.Shutdown при сигнале
	shutdown.WaitForSignals(5*time.Second, application.GRPCSrv)

	log.Info("Application stopped")

}
