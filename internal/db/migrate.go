package db

import (
	"database/sql"
	"log/slog"

	"github.com/pressly/goose/v3"
)

// Migrate применяет миграции из указанной директории к базе
func Migrate(db *sql.DB, dir string, logger *slog.Logger) {
	const op = "db.Migrate"
	log := logger.With(slog.String("op", op))

	// Устанавливаем диалект базы
	if err := goose.SetDialect("postgres"); err != nil {
		log.Error("failed to set dialect", "error", err)
		panic(err)
	}

	// Применяем все миграции из директории
	if err := goose.Up(db, dir); err != nil {
		log.Error("failed to apply migrations", "error", err)
		panic(err)
	}

	log.Info("Migrations applied successfully")
}
