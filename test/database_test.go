package test

import (
	"os"
	"testing"

	"github.com/Dei-web/Go-inventarie/internal/config"
	"github.com/Dei-web/Go-inventarie/internal/db"
	"github.com/joho/godotenv"
)

func TestDatabaseConnection(t *testing.T) {
	if err := godotenv.Load("../.env"); err != nil {
		t.Fatal("No se pudo cargar .env:", err)
	}

	if os.Getenv("JWT_SECRET") == "" {
		os.Setenv("JWT_SECRET", "test-secret")
	}

	cfg, err := config.Load()
	if err != nil {
		t.Fatal("No se pudo cargar la configuración:", err)
	}

	database, err := db.Connect(cfg)
	if err != nil {
		t.Fatal("No se pudo conectar a la DB:", err)
	}

	sqlDB, err := database.DB()
	if err != nil {
		t.Fatal("No se pudo obtener la conexión:", err)
	}

	if err := sqlDB.Ping(); err != nil {
		t.Fatal("La DB no responde:", err)
	}

	t.Log(" Conexión a PostgreSQL exitosa")
}
