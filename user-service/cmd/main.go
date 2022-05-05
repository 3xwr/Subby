package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"

	"user-service/internal/config"
	"user-service/internal/handler"
	"user-service/internal/repository"
	"user-service/internal/service"
)

func main() {
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	cfg, err := config.New()
	if err != nil {
		logger.Fatal().Err(err).Msg("Configuration error")
	}

	r := chi.NewRouter()

	db, err := sql.Open("postgres", cfg.DBConnString)
	if err != nil {
		logger.Fatal().Err(err).Msg("DB initialization error")
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		logger.Fatal().Err(err).Msg("DB pinging error")
	}

	authRepo := repository.NewAuth(db)
	authService := service.NewAuth(&logger, authRepo, []byte(cfg.Secret))
	authHandler := handler.NewAuth(&logger, authService)
	registerHandler := handler.NewRegister(&logger, authService)

	resourceHandler := handler.NewResource(&logger)

	r.Route("/", func(r chi.Router) {
		r.Use(cors.Handler(cors.Options{
			AllowedOrigins: []string{"https://*", "http://*"},
			AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		}))
		r.Use(middleware.RequestLogger(&handler.LogFormatter{Logger: &logger}))
		r.Use(middleware.Recoverer)
		r.Use(handler.JWT([]byte(cfg.Secret)))
		r.Method(http.MethodPost, handler.AuthPath, authHandler)
		r.Method(http.MethodPost, handler.RegisterPath, registerHandler)
		r.Method(http.MethodGet, handler.ResourcePath, resourceHandler)
	})

	srv := http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: r,
	}

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(shutdown)

	go func() {
		logger.Info().Msgf("Server is listening on :%d", cfg.Port)
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			logger.Fatal().Err(err).Msg("Server error")
		}
	}()

	<-shutdown

	logger.Info().Msg("Shutdown signal received")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer func() {
		cancel()
	}()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal().Err(err).Msg("Server shutdown error")
	}

	logger.Info().Msg("Server stopped gracefully")
}
