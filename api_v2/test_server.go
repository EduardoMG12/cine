package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

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
                                                  
 🎬 CineVerse API Test - Health Check Only
`

func main() {
	// Display banner
	fmt.Print(banner)

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		slog.Error("Failed to load configuration", "error", err)
		os.Exit(1)
	}

	// Setup logging
	logger := appMiddleware.SetupLogger(cfg.Server.Environment)
	slog.SetDefault(logger)

	// Initialize i18n
	localizer, err := i18n.NewLocalizer()
	if err != nil {
		slog.Error("Failed to initialize i18n", "error", err)
		os.Exit(1)
	}

	// Setup router
	router := setupRouter(cfg, localizer)

	// Create HTTP server
	server := &http.Server{
		Addr:         cfg.Server.GetServerAddress(),
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  60 * time.Second,
	}

	slog.Info("Test server starting (no database connection)",
		"address", cfg.Server.GetServerAddress(),
		"environment", cfg.Server.Environment,
	)

	if err := server.ListenAndServe(); err != nil {
		slog.Error("Server failed to start", "error", err)
		os.Exit(1)
	}
}

func setupRouter(cfg *config.Config, localizer *i18n.Localizer) *chi.Mux {
	router := chi.NewRouter()

	// Basic middleware
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(appMiddleware.LoggerMiddleware())
	router.Use(middleware.Recoverer)
	router.Use(appMiddleware.Language(localizer))

	// CORS middleware for development
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

	// Initialize handlers
	healthHandler := handler.NewHealthHandler(localizer)

	// Routes
	router.Get("/health", healthHandler.Health)

	// API versioning
	router.Route("/api/v2", func(r chi.Router) {
		// Health check
		r.Get("/health", healthHandler.Health)

		// Welcome endpoint
		r.Get("/", func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"message": "Welcome to CineVerse API v2.0.0", "status": "ok"}`))
		})
	})

	return router
}
