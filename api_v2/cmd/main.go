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

	"github.com/EduardoMG12/cine/api_v2/internal/handler"
	"github.com/EduardoMG12/cine/api_v2/internal/repository"
	"github.com/EduardoMG12/cine/api_v2/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

func main() {
	config := loadConfig()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	db, err := sqlx.Connect("postgres", config.DatabaseURL)
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	rdb := redis.NewClient(&redis.Options{
		Addr:     config.RedisAddr,
		Password: config.RedisPassword,
		DB:       0,
	})
	defer rdb.Close()

	ctx := context.Background()
	if err := rdb.Ping(ctx).Err(); err != nil {
		slog.Error("Failed to connect to Redis", "error", err)
		os.Exit(1)
	}

	userRepo := repository.NewUserRepository(db, rdb)

	userService := service.NewUserService(userRepo)

	userHandler := handler.NewUserHandler(userService)

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

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

	// API routes
	r.Route("/api/v1", func(r chi.Router) {
		r.Mount("/users", userHandler.Routes())
	})

	// Server
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", config.Port),
		Handler: r,
	}

	// Graceful shutdown
	go func() {
		slog.Info("Server starting", "port", config.Port)
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

type Config struct {
	Port          string
	DatabaseURL   string
	RedisAddr     string
	RedisPassword string
}

func loadConfig() *Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")

	viper.AutomaticEnv()
	viper.SetEnvPrefix("CINE")

	viper.SetDefault("PORT", "8080")
	viper.SetDefault("DATABASE_URL", "postgres://cineverse:password@postgres:5432/cineverse?sslmode=disable")
	viper.SetDefault("REDIS_ADDR", "redis:6379")
	viper.SetDefault("REDIS_PASSWORD", "")

	if err := viper.ReadInConfig(); err != nil {
		slog.Warn("Config file not found, using environment variables and defaults")
	}

	return &Config{
		Port:          viper.GetString("PORT"),
		DatabaseURL:   viper.GetString("DATABASE_URL"),
		RedisAddr:     viper.GetString("REDIS_ADDR"),
		RedisPassword: viper.GetString("REDIS_PASSWORD"),
	}
}
