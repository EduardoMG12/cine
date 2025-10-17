package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/EduardoMG12/cine/api_v2/internal/config"
	"github.com/EduardoMG12/cine/api_v2/internal/handler"
	customMiddleware "github.com/EduardoMG12/cine/api_v2/internal/middleware"
	"github.com/EduardoMG12/cine/api_v2/internal/repository"
	"github.com/EduardoMG12/cine/api_v2/internal/service"
)

type Server struct {
	config     *config.Config
	db         *sqlx.DB
	httpServer *http.Server
	logger     *slog.Logger
}

func NewServer(cfg *config.Config, logger *slog.Logger) (*Server, error) {
	// Connect to database
	db, err := sqlx.Connect("postgres", cfg.Database.GetDSN())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Configure database connection pool
	db.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	db.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	db.SetConnMaxLifetime(time.Duration(cfg.Database.ConnMaxLifetime) * time.Minute)

	// Test database connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	server := &Server{
		config: cfg,
		db:     db,
		logger: logger,
	}

	// Setup HTTP server
	server.setupHTTPServer()

	return server, nil
}

func (s *Server) setupHTTPServer() {
	r := chi.NewRouter()

	// Global middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	// CORS middleware
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, X-CSRF-Token")

			if r.Method == "OPTIONS" {
				return
			}

			next.ServeHTTP(w, r)
		})
	})

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy","service":"cineverse-api"}`))
	})

	// Initialize repositories
	userRepo := repository.NewUserRepository(s.db)
	sessionRepo := repository.NewSessionRepository(s.db)

	// Initialize services
	passwordService := service.NewPasswordService()
	jwtService := service.NewJWTService(s.config.JWT.Secret)
	authService := service.NewAuthService(userRepo, sessionRepo, passwordService, jwtService)

	// Initialize middleware
	authMiddleware := customMiddleware.AuthMiddleware(authService)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authService)

	// Setup API routes
	r.Route("/api/v1", func(r chi.Router) {
		// Authentication routes
		authHandler.RegisterRoutes(r, authMiddleware)
	})

	s.httpServer = &http.Server{
		Addr:         s.config.Server.GetServerAddress(),
		Handler:      r,
		ReadTimeout:  s.config.Server.ReadTimeout,
		WriteTimeout: s.config.Server.WriteTimeout,
	}
}

func (s *Server) Start() error {
	s.logger.Info("Starting CineVerse API server",
		"address", s.config.Server.GetServerAddress(),
		"environment", s.config.Server.Environment)

	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	s.logger.Info("Shutting down server...")

	// Close database connection
	if s.db != nil {
		if err := s.db.Close(); err != nil {
			s.logger.Error("Failed to close database connection", "error", err)
		}
	}

	// Shutdown HTTP server
	return s.httpServer.Shutdown(ctx)
}
