package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		DatabaseURL: os.Getenv("URLDATABASE"),
	}

	if cfg.DatabaseURL == "" {
		return nil, ErrDatabaseURLMissing
	}

	return cfg, nil
}
