// @title CineVerse API
// @version 2.0
// @description A comprehensive social network platform for movie enthusiasts
// @termsOfService http://cineverse.com/terms/

// @contact.name CineVerse API Support
// @contact.url http://www.cineverse.com/support
// @contact.email support@cineverse.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/EduardoMG12/cine/api_v2/internal/config"
	"github.com/EduardoMG12/cine/api_v2/internal/handler"
	customMiddleware "github.com/EduardoMG12/cine/api_v2/internal/middleware"
	"github.com/EduardoMG12/cine/api_v2/internal/repository"
	"github.com/EduardoMG12/cine/api_v2/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	httpSwagger "github.com/swaggo/http-swagger"

	_ "github.com/EduardoMG12/cine/api_v2/docs"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("Failed to load configuration", "error", err)
		os.Exit(1)
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	db, err := sqlx.Connect("postgres", cfg.Database.URL)
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	// Configure database connection pool
	db.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	db.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	db.SetConnMaxLifetime(time.Duration(cfg.Database.ConnMaxLifetime) * time.Second)

	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})
	defer rdb.Close()

	ctx := context.Background()
	if err := rdb.Ping(ctx).Err(); err != nil {
		slog.Error("Failed to connect to Redis", "error", err)
		os.Exit(1)
	}

	userRepo := repository.NewUserRepository(db, rdb)
	userSessionRepo := repository.NewUserSessionRepository(db, rdb)
	friendshipRepo := repository.NewFriendshipRepository(db)
	followRepo := repository.NewFollowRepository(db)

	userService := service.NewUserService(userRepo)
	sessionDuration := time.Duration(cfg.Application.SessionDurationHours) * time.Hour
	userSessionService := service.NewUserSessionService(userSessionRepo, sessionDuration)
	socialService := service.NewSocialService(friendshipRepo, followRepo, userRepo)

	userHandler := handler.NewUserHandler(userService, userSessionService)
	socialHandler := handler.NewSocialHandler(socialService)

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(customMiddleware.LanguageMiddleware)

	// CORS
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("CineVerse API OK - Hot Reload Working!"))
	})

	// Swagger documentation
	if cfg.Application.EnableSwagger {
		r.Get("/swagger/*", httpSwagger.Handler(
			httpSwagger.URL(cfg.Application.SwaggerURL),
		))
	}

	// API routes
	r.Route("/api/v1", func(r chi.Router) {
		r.Mount("/users", userHandler.Routes())

		// Protected social routes
		r.Group(func(r chi.Router) {
			r.Use(customMiddleware.AuthMiddleware)
			r.Mount("/social", socialHandler.Routes())
		})
	})

	// Server
	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port),
		Handler: r,
	}

	// Graceful shutdown
	go func() {
		slog.Info("Server starting", "port", cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Server failed to start", "error", err)
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("Server shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("Server shutdown failed", "error", err)
		os.Exit(1)
	}

	slog.Info("Server stopped")
}
