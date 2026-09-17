package db

import (
	"github.com/Dei-web/Go-inventarie/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(cfg *config.Config) (*gorm.DB, error) {
	database, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err := Migrate(database); err != nil {
		return nil, err
	}

	return database, nil
}
