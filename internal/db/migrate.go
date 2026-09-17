package db

import (
	"fmt"
	"log"
	"reflect"

	"github.com/Dei-web/Go-inventarie/internal/models"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	log.Println("Starting database migration...")

	for _, model := range models.Models {
		modelName := reflect.TypeOf(model).Elem().Name()
		log.Printf("Migrating model: %s", modelName)

		if err := db.AutoMigrate(model); err != nil {
			return fmt.Errorf("failed to migrate %s: %w", modelName, err)
		}

		log.Printf("✓ %s migrated successfully", modelName)
	}

	log.Println("Database migration completed successfully")
	return nil
}
