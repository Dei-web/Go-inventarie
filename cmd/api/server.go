package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Dei-web/Go-inventarie/internal/config"
	"github.com/Dei-web/Go-inventarie/internal/middleware"
	"github.com/Dei-web/Go-inventarie/internal/middleware/auth"
	"github.com/Dei-web/Go-inventarie/internal/services/user"
	"gorm.io/gorm"
)

type Server struct {
	config *config.Config
	db     *gorm.DB
}

func NewServer(cfg *config.Config, db *gorm.DB) *Server {
	return &Server{config: cfg, db: db}
}

func (s *Server) Run() error {
	userStore := user.NewStore(s.db)
	userService := user.NewService(userStore)
	userHandler := user.NewHandler(userService)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /users", userHandler.GetUsers)
	mux.HandleFunc("GET /users/{id}", userHandler.GetUser)
	mux.HandleFunc("POST /users", userHandler.CreateUser)

	authMw := auth.Middleware(s.config.JWTSecret)
	mux.Handle("PUT /users/{id}", authMw(http.HandlerFunc(userHandler.UpdateUser)))
	mux.Handle("DELETE /users/{id}", authMw(http.HandlerFunc(userHandler.DeleteUser)))

	var handler http.Handler = mux
	handler = middleware.CORS(handler)
	handler = middleware.Logging(handler)
	handler = middleware.Recovery(handler)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", s.config.Port),
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		log.Println("shutting down server...")
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("server forced to shutdown: %v", err)
		}
	}()

	log.Printf("API running on :%s [%s]", s.config.Port, s.config.Environment)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}

	log.Println("server stopped gracefully")
	return nil
}
