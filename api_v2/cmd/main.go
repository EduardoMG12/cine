package main

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/lib/pq"

	"github.com/EduardoMG12/cine/api_v2/internal/config"
	"github.com/EduardoMG12/cine/api_v2/internal/handler"
	"github.com/EduardoMG12/cine/api_v2/internal/i18n"
	appMiddleware "github.com/EduardoMG12/cine/api_v2/internal/middleware"
)

const banner = `
 ██████╗██╗███╗   ██╗███████╗    ██╗   ██╗██████╗ 
██╔════╝██║████╗  ██║██╔════╝    ██║   ██║╚════██╗
██║     ██║██╔██╗ ██║█████╗      ██║   ██║ █████╔╝
██║     ██║██║╚██╗██║██╔══╝      ╚██╗ ██╔╝██╔═══╝ 
╚██████╗██║██║ ╚████║███████╗     ╚████╔╝ ███████╗
 ╚═════╝╚═╝╚═╝  ╚═══╝╚══════╝      ╚═══╝  ╚══════╝
                                                  
 🎬 CineVerse API v2.0.0 - Movie Social Network    
`

func main() {
	fmt.Print(banner)

	cfg, err := config.Load()
	if err != nil {
		slog.Error("Failed to load configuration", "error", err)
		os.Exit(1)
	}

	logger := appMiddleware.SetupLogger(cfg.Server.Environment)
	slog.SetDefault(logger)

	localizer, err := i18n.NewLocalizer()
	if err != nil {
		slog.Error("Failed to initialize i18n", "error", err)
		os.Exit(1)
	}

	db, err := connectDatabase(cfg)
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	router := setupRouter(cfg, localizer)

	server := &http.Server{
		Addr:         cfg.Server.GetServerAddress(),
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		slog.Info("Starting server",
			"address", cfg.Server.GetServerAddress(),
			"environment", cfg.Server.Environment,
		)

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Server failed to start", "error", err)
			os.Exit(1)
		}
	}()

	slog.Info("Server started successfully", "address", cfg.Server.GetServerAddress())

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("Server is shutting down gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("Server forced to shutdown", "error", err)
		os.Exit(1)
	}

	slog.Info("Server stopped successfully")
}

func connectDatabase(cfg *config.Config) (*sql.DB, error) {
	slog.Info("Connecting to database...")

	dsn := cfg.Database.GetDSN()

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	db.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	db.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	db.SetConnMaxLifetime(time.Duration(cfg.Database.ConnMaxLifetime) * time.Minute)

	slog.Info("Database connected successfully")
	return db, nil
}

func setupRouter(cfg *config.Config, localizer *i18n.Localizer) *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(appMiddleware.LoggerMiddleware())
	router.Use(middleware.Recoverer)
	router.Use(appMiddleware.Language(localizer))

	if cfg.Server.Environment == "development" {
		router.Use(middleware.AllowContentType("application/json"))
		router.Use(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Access-Control-Allow-Origin", "*")
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Accept-Language")

				if r.Method == "OPTIONS" {
					w.WriteHeader(http.StatusOK)
					return
				}

				next.ServeHTTP(w, r)
			})
		})
	}

	healthHandler := handler.NewHealthHandler(localizer)

	router.Get("/health", healthHandler.Health)

	router.Route("/api/v2", func(r chi.Router) {
		// Health check (duplicated for convenience)
		r.Get("/health", healthHandler.Health)

		// TODO: Add other API routes here in future sprints
		// r.Route("/auth", authRoutes)
		// r.Route("/users", userRoutes)
		// r.Route("/movies", movieRoutes)
		// etc.
	})

	return router
}
