package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/EduardoMG12/cine/api_v2/internal/auth"
	"github.com/EduardoMG12/cine/api_v2/internal/config"
	"github.com/EduardoMG12/cine/api_v2/internal/handler"
	"github.com/EduardoMG12/cine/api_v2/internal/repository"
	"github.com/EduardoMG12/cine/api_v2/internal/service"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

type Server struct {
	cfg        *config.Config
	httpServer *http.Server
	db         *sqlx.DB
	redis      *redis.Client
}

// New creates a new server instance with all dependencies
func New(cfg *config.Config, db *sqlx.DB, redis *redis.Client) *Server {
	return &Server{
		cfg:   cfg,
		db:    db,
		redis: redis,
	}
}

// Start initializes and starts the server
func (s *Server) Start() error {
	// Initialize repositories
	userRepo := repository.NewUserRepository(s.db, s.redis)
	sessionRepo := repository.NewUserSessionRepository(s.db, s.redis)

	// Initialize auth components
	jwtManager := auth.NewJWTManager(s.cfg.JWT.Secret, time.Duration(s.cfg.JWT.Expiration)*time.Hour)
	passwordHasher := auth.NewPasswordHasher()

	// Initialize services
	userService := service.NewUserService(userRepo)
	sessionService := service.NewUserSessionService(sessionRepo, time.Duration(s.cfg.JWT.Expiration)*time.Hour)
	authService := service.NewAuthService(userRepo, sessionRepo, jwtManager, passwordHasher, s.cfg)

	// Initialize handlers
	userHandler := handler.NewUserHandler(userService)
	authHandler := handler.NewAuthHandler(
		userService,
		sessionService,
		jwtManager,
		passwordHasher,
	)

	// Setup router
	router := NewRouter(userHandler, authHandler)

	// Create HTTP server
	s.httpServer = &http.Server{
		Addr:         fmt.Sprintf("%s:%s", s.cfg.Server.Host, s.cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		slog.Info("Server starting", "port", s.cfg.Server.Port, "host", s.cfg.Server.Host)
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Server failed to start", "error", err)
			os.Exit(1)
		}
	}()

	return nil
}

// Stop gracefully shuts down the server
func (s *Server) Stop() error {
	slog.Info("Server shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("server shutdown failed: %w", err)
	}

	slog.Info("Server stopped successfully")
	return nil
}

// WaitForShutdown blocks until a termination signal is received
func (s *Server) WaitForShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}
