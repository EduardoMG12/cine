package server

import (
	"net/http"
	"time"

	"github.com/EduardoMG12/cine/api_v2/internal/handler"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

// Router creates and configures the main router with all routes
func NewRouter(userHandler *handler.UserHandler, authHandler *handler.AuthHandler) chi.Router {
	r := chi.NewRouter()

	setupMiddleware(r)
	setupRoutes(r, userHandler, authHandler)

	return r
}

func setupMiddleware(r chi.Router) {
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	// CORS configuration
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))
}

func setupRoutes(r chi.Router, userHandler *handler.UserHandler, authHandler *handler.AuthHandler) {
	// Health check
	r.Get("/health", healthCheckHandler)

	// API v1 routes
	r.Route("/api/v1", func(r chi.Router) {
		// Authentication routes
		r.Mount("/auth", authHandler.Routes())

		// User routes
		r.Mount("/users", userHandler.Routes())

		// Movie routes (will be added later)
		r.Route("/movies", func(r chi.Router) {
			// Will be implemented in next step
		})

		// Match session routes (will be added later)
		r.Route("/match-sessions", func(r chi.Router) {
			// Will be implemented in next step
		})
	})
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("CineVerse API v2 - Healthy"))
}
