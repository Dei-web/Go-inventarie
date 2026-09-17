package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string
	Port        string
	JWTSecret   string
	Environment string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		DatabaseURL: os.Getenv("URLDATABASE"),
		Port:        os.Getenv("PORT"),
		JWTSecret:   os.Getenv("JWT_SECRET"),
		Environment: os.Getenv("ENVIRONMENT"),
	}

	if cfg.DatabaseURL == "" {
		return nil, ErrDatabaseURLMissing
	}

	if cfg.Port == "" {
		cfg.Port = "8080"
	}

	if cfg.JWTSecret == "" {
		return nil, ErrJWTSecretMissing
	}

	if cfg.Environment == "" {
		cfg.Environment = "development"
	}

	return cfg, nil
}
