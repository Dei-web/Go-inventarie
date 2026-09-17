package main

import (
	"log"

	"github.com/Dei-web/Go-inventarie/cmd/api"
	"github.com/Dei-web/Go-inventarie/internal/config"
	"github.com/Dei-web/Go-inventarie/internal/db"
)

// @title           Go Inventarie API
// @version         1.0
// @description     API CRUD para gestión de usuarios
// @host            localhost:8080
// @BasePath        /
// @schemes         http

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	database, err := db.Connect(cfg)
	if err != nil {
		log.Fatal(err)
	}

	server := api.NewServer(cfg, database)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
