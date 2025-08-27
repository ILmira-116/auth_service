package db

import (
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	config "auth-service/config"

	_ "github.com/jackc/pgx/v5/stdlib" // драйвер PostgreSQL для database/sql
)

// InitPostgres создает подключение к PostgreSQL и возвращает *sql.DB
func InitPostgres(cfg *config.DBConfig, logger *slog.Logger) *sql.DB {
	const op = "db.Postgres"

	log := logger.With(slog.String("op", op))

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name,
		cfg.SSLMode,
	)

	// Открываем подключение (но пока без проверки)
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Error("Failed to open DB connection", "error", err)
		panic(err)
	}

	// Пинг с повторными попытками
	var pingErr error
	for i := 0; i < 10; i++ {
		pingErr = db.Ping()
		if pingErr == nil {
			log.Info("Successfully connected to PostgreSQL")
			break
		}
		log.Warn("DB not ready yet, retrying", "attempt", i+1, "error", pingErr)
		time.Sleep(2 * time.Second)
	}

	if pingErr != nil {
		log.Error("Failed to connect to PostgreSQL after 10 attempts", "error", pingErr)
		panic(pingErr)
	}

	// Настройка пула соединений
	db.SetMaxOpenConns(25)                 // максимум открытых соединений
	db.SetMaxIdleConns(25)                 // максимум неактивных соединений
	db.SetConnMaxLifetime(5 * time.Minute) // время жизни соединения

	return db
}
