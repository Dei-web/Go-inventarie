package db

import (
	"github.com/Dei-web/Go-inventarie/internal/config"
	"github.com/Dei-web/Go-inventarie/internal/types"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(cfg *config.Config) (*gorm.DB, error) {
	database, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err := AutoMigrate(database); err != nil {
		return nil, err
	}

	return database, nil
}

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(&types.Users{})
}
