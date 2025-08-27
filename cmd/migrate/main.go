package main

import (
	"auth-service/config"
	"auth-service/internal/db"
	"auth-service/internal/logger"

	"log"
)

func main() {
	// 1. Загружаем конфиг
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// 2. Инициализация логгера
	log := logger.New(cfg)
	log.Info("Logger is ready")
	log.Debug("Debug message")

	// 3. Подключаемся к базе
	dbConn := db.InitPostgres(&cfg.DB, log)

	// 4. Применяем миграции
	migrationsDir := "./migrations"
	db.Migrate(dbConn, migrationsDir, log)

	log.Info("All migrations applied successfully")
}
