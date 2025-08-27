package logger

import (
	"auth-service/config"
	"log/slog"
	"os"
)

const (
	EnvLocal = "local"
	EnvDev   = "dev"
	EnvProd  = "prod"
)

// New создаёт slog.Logger на основе конфигурации
func New(cfg *config.Config) *slog.Logger {
	// Определяем уровень логирования из cfg.LogLevel
	level := slog.LevelInfo
	switch cfg.LogLevel {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	}

	// Выбираем формат логов по cfg.Env
	var handler slog.Handler
	switch cfg.Env {
	case EnvLocal:
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: level})
	case EnvDev, EnvProd:
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level})
	default:
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: level})
	}

	return slog.New(handler)
}
