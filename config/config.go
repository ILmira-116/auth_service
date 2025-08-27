package config

import (
	"log/slog"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	GRPC      GRPCConfig
	TokenTTL  time.Duration
	DB        DBConfig
	LogLevel  string `env:"LOG_LEVEL" env-default:"info"`
	Env       string `env:"ENV" env-default:"local"`
	Logger    *slog.Logger
	JWTSecret string `env:"JWT_SECRET,required"`
}

type DBConfig struct {
	Host     string `env:"DB_HOST" env-default:"auth_db"`
	Port     string `env:"DB_PORT" env-default:"5432"`
	User     string `env:"DB_USER" env-default:"auth"`
	Password string `env:"DB_PASSWORD" env-default:"authpass"`
	Name     string `env:"DB_NAME" env-default:"auth_db"`
	SSLMode  string `env:"DB_SSLMODE" env-default:"disable"`
}

type GRPCConfig struct {
	ServerHost         string        `env:"GRPC_SERVER_HOST" env-default:"0.0.0.0"`
	ServerPort         string        `env:"GRPC_SERVER_PORT" env-default:"50051"`
	ServerReadTimeout  time.Duration `env:"GRPC_SERVER_READ_TIMEOUT" env-default:"5s"`
	ServerWriteTimeout time.Duration `env:"GRPC_SERVER_WRITE_TIMEOUT" env-default:"10s"`
	ServerIdleTimeout  time.Duration `env:"GRPC_SERVER_IDLE_TIMEOUT" env-default:"120s"`
}

func LoadConfig() (*Config, error) {
	cfg := &Config{}
	err := cleanenv.ReadEnv(cfg)

	if err != nil {
		return nil, err
	}

	return cfg, nil
}
