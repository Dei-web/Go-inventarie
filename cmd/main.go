package main

import (
	"log"
	"net/http"

	"github.com/Dei-web/Go-inventarie/internal/config"
	"github.com/Dei-web/Go-inventarie/internal/db"
	"github.com/Dei-web/Go-inventarie/internal/services/users"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	database, err := db.Connect(cfg)
	if err != nil {
		log.Fatal(err)
	}

	userStore := users.NewStore(database)
	userHandler := users.NewHandler(userStore)

	mux := http.NewServeMux()

	users.RegisterRoutes(mux, userHandler)

	log.Println("API running on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
