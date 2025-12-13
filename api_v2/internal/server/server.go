package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	httpSwagger "github.com/swaggo/http-swagger"

	_ "github.com/EduardoMG12/cine/api_v2/docs"
	"github.com/EduardoMG12/cine/api_v2/internal/config"
	httpHandler "github.com/EduardoMG12/cine/api_v2/internal/handler/http"
	"github.com/EduardoMG12/cine/api_v2/internal/infrastructure"
	customMiddleware "github.com/EduardoMG12/cine/api_v2/internal/middleware"
	"github.com/EduardoMG12/cine/api_v2/internal/repository"
	"github.com/EduardoMG12/cine/api_v2/internal/usecase/auth"
	"github.com/EduardoMG12/cine/api_v2/internal/usecase/movie"
	"github.com/EduardoMG12/cine/api_v2/internal/usecase/user"
	"github.com/EduardoMG12/cine/api_v2/internal/usecase/user_movie"
)

type Server struct {
	config     *config.Config
	db         *sqlx.DB
	httpServer *http.Server
	logger     *slog.Logger
	router     *chi.Mux
}

type RouteInfo struct {
	Method  string
	Path    string
	Handler string
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
	s.router = r

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

	// Initialize infrastructure
	passwordService := infrastructure.NewPasswordService()
	jwtService := infrastructure.NewJWTService(s.config.JWT.Secret)
	_, err := infrastructure.NewRedisService(s.config.Redis.Host, s.config.Redis.Port, s.config.Redis.Password, s.config.Redis.DB)
	if err != nil {
		s.logger.Error("Failed to initialize Redis service", "error", err)
		// Continue without Redis, caching will be disabled
	}
	omdbService := infrastructure.NewOMDbService(s.config.OMDb.APIKey)

	// Initialize repositories
	userRepo := repository.NewUserRepository(s.db)
	sessionRepo := repository.NewSessionRepository(s.db)
	movieRepo := repository.NewMovieRepository(s.db)
	watchedMovieRepo := repository.NewWatchedMovieRepository(s.db)
	favoriteMovieRepo := repository.NewFavoriteMovieRepository(s.db)

	// Initialize movie fetcher chain: OMDb -> Database (with auto-save enabled)
	movieFetcher := infrastructure.NewMovieFetcherChain(omdbService, movieRepo, true)

	// Initialize auth use cases
	registerUC := auth.NewRegisterUseCase(userRepo, sessionRepo, passwordService, jwtService)
	loginUC := auth.NewLoginUseCase(userRepo, sessionRepo, passwordService, jwtService)
	getMeUC := auth.NewGetMeUseCase(userRepo)
	logoutUC := auth.NewLogoutUseCase(sessionRepo)
	logoutAllUC := auth.NewLogoutAllUseCase(sessionRepo)

	// Initialize movie use cases
	getMovieByIDUC := movie.NewGetMovieByIDUseCase(movieFetcher)
	getRandomMovieUC := movie.NewGetRandomMovieUseCase(movieRepo)
	getRandomMovieByGenreUC := movie.NewGetRandomMovieByGenreUseCase(movieRepo)
	searchMoviesUC := movie.NewSearchMoviesUseCase(movieFetcher)
	getTrendingMoviesUC := movie.NewGetTrendingMoviesUseCase(movieRepo, movieFetcher)

	// Initialize user movie use cases
	toggleWatchedMovieUC := user_movie.NewToggleWatchedMovieUseCase(watchedMovieRepo, movieRepo)
	getWatchedMoviesUC := user_movie.NewGetWatchedMoviesUseCase(watchedMovieRepo, movieRepo)
	toggleFavoriteMovieUC := user_movie.NewToggleFavoriteMovieUseCase(favoriteMovieRepo, movieRepo)
	getFavoriteMoviesUC := user_movie.NewGetFavoriteMoviesUseCase(favoriteMovieRepo, movieRepo)

	// Initialize user use cases
	updateUserUC := user.NewUpdateUserUseCase(userRepo)

	// Initialize handlers
	systemHandler := httpHandler.NewSystemHandler()
	authHandler := httpHandler.NewAuthHandler(registerUC, loginUC, getMeUC, logoutUC, logoutAllUC)
	movieHandler := httpHandler.NewMovieHandler(
		getMovieByIDUC,
		getRandomMovieUC,
		getRandomMovieByGenreUC,
		searchMoviesUC,
		getTrendingMoviesUC,
	)
	omdbHandler := httpHandler.NewOMDbHandler(omdbService)
	watchedMovieHandler := httpHandler.NewWatchedMovieHandler(toggleWatchedMovieUC, getWatchedMoviesUC)
	favoriteMovieHandler := httpHandler.NewFavoriteMovieHandler(toggleFavoriteMovieUC, getFavoriteMoviesUC)
	userHandler := httpHandler.NewUserHandler(updateUserUC)

	// System routes
	r.Get("/", systemHandler.Root)
	r.Get("/health", systemHandler.HealthCheck)

	// Initialize middleware
	authMiddleware := customMiddleware.JWTAuthMiddleware(jwtService, userRepo)

	// Setup API routes
	r.Route("/api/v1", func(r chi.Router) {
		// Auth routes
		r.Route("/auth", func(r chi.Router) {
			// Public routes
			r.Post("/register", authHandler.Register)
			r.Post("/login", authHandler.Login)

			// Protected routes
			r.Group(func(r chi.Router) {
				r.Use(authMiddleware)
				r.Get("/me", authHandler.GetMe)
				r.Post("/logout", authHandler.Logout)
				r.Post("/logout-all", authHandler.LogoutAll)
			})
		})

		// Movie routes (all public)
		r.Route("/movies", func(r chi.Router) {
			r.Get("/trending", movieHandler.GetTrendingMovies)
			r.Get("/random", movieHandler.GetRandomMovie)
			r.Get("/random-by-genre", movieHandler.GetRandomMovieByGenre)
			r.Get("/search", movieHandler.SearchMovies)
			r.Get("/{id}", movieHandler.GetMovieByID)
		})

		// Watched movies routes (protected)
		r.Route("/watched", func(r chi.Router) {
			r.Use(authMiddleware)
			r.Get("/", watchedMovieHandler.GetWatchedMovies)
			r.Post("/", watchedMovieHandler.ToggleWatchedMovie)
		})

		// Favorite movies routes (protected)
		r.Route("/favorites", func(r chi.Router) {
			r.Use(authMiddleware)
			r.Get("/", favoriteMovieHandler.GetFavoriteMovies)
			r.Post("/", favoriteMovieHandler.ToggleFavoriteMovie)
		})

		// User routes (protected)
		r.Route("/users", func(r chi.Router) {
			r.Use(authMiddleware)
			r.Patch("/me", userHandler.UpdateUser)
		})

		// OMDb routes (test and search)
		r.Route("/omdb", func(r chi.Router) {
			r.Get("/test", omdbHandler.TestConnection)
			r.Get("/{imdbId}", omdbHandler.GetMovieByIMDbID)
			r.Get("/title", omdbHandler.GetMovieByTitle)
			r.Get("/search", omdbHandler.SearchMovies)
			r.Get("/search-by-type", omdbHandler.SearchMoviesByType)
		})
	})

	// Swagger documentation route
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
	))

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

	// Print routes
	s.printRoutes()

	// Print Swagger documentation URL
	fmt.Println("┌─────────────────────────────────────────────────────────────────┐")
	fmt.Println("│                  📚 API DOCUMENTATION                           │")
	fmt.Println("└─────────────────────────────────────────────────────────────────┘")
	fmt.Printf("\n  For more information, access: http://localhost:%s/swagger/index.html\n\n", s.config.Server.Port)
	fmt.Println("─────────────────────────────────────────────────────────────────")
	fmt.Println()

	return s.httpServer.ListenAndServe()
}

func (s *Server) printRoutes() {
	routes := s.getRoutes()

	fmt.Println("\n┌─────────────────────────────────────────────────────────────────┐")
	fmt.Println("│                     📡 REGISTERED ROUTES                        │")
	fmt.Println("└─────────────────────────────────────────────────────────────────┘")
	fmt.Println()

	// Group routes by prefix
	publicRoutes := []RouteInfo{}
	protectedRoutes := []RouteInfo{}
	movieRoutes := []RouteInfo{}
	userMovieRoutes := []RouteInfo{}
	userRoutes := []RouteInfo{}

	for _, route := range routes {
		if strings.Contains(route.Path, "/movies") {
			movieRoutes = append(movieRoutes, route)
		} else if strings.Contains(route.Path, "/watched") || strings.Contains(route.Path, "/favorites") {
			userMovieRoutes = append(userMovieRoutes, route)
		} else if strings.Contains(route.Path, "/users") {
			userRoutes = append(userRoutes, route)
		} else if strings.Contains(route.Path, "/auth/me") ||
			strings.Contains(route.Path, "/auth/logout") {
			protectedRoutes = append(protectedRoutes, route)
		} else {
			publicRoutes = append(publicRoutes, route)
		}
	}

	if len(publicRoutes) > 0 {
		fmt.Println("  🌍 Public Routes:")
		for _, route := range publicRoutes {
			fmt.Printf("    %-7s %s\n", colorizeMethod(route.Method), route.Path)
		}
		fmt.Println()
	}

	if len(protectedRoutes) > 0 {
		fmt.Println("  🔒 Protected Routes (require JWT):")
		for _, route := range protectedRoutes {
			fmt.Printf("    %-7s %s\n", colorizeMethod(route.Method), route.Path)
		}
		fmt.Println()
	}

	if len(movieRoutes) > 0 {
		fmt.Println("  🎬 Movie Routes:")
		for _, route := range movieRoutes {
			fmt.Printf("    %-7s %s\n", colorizeMethod(route.Method), route.Path)
		}
		fmt.Println()
	}

	if len(userMovieRoutes) > 0 {
		fmt.Println("  ⭐ User Movie Lists (require JWT):")
		for _, route := range userMovieRoutes {
			fmt.Printf("    %-7s %s\n", colorizeMethod(route.Method), route.Path)
		}
		fmt.Println()
	}

	if len(userRoutes) > 0 {
		fmt.Println("  👤 User Profile (require JWT):")
		for _, route := range userRoutes {
			fmt.Printf("    %-7s %s\n", colorizeMethod(route.Method), route.Path)
		}
		fmt.Println()
	}

	fmt.Println("─────────────────────────────────────────────────────────────────")
	fmt.Println()
}

func (s *Server) getRoutes() []RouteInfo {
	var routes []RouteInfo

	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		if route != "" && route != "/*" {
			routes = append(routes, RouteInfo{
				Method:  method,
				Path:    route,
				Handler: "",
			})
		}
		return nil
	}

	chi.Walk(s.router, walkFunc)
	return routes
}

func colorizeMethod(method string) string {
	switch method {
	case "GET":
		return "\033[32mGET\033[0m   "
	case "POST":
		return "\033[33mPOST\033[0m  "
	case "PUT":
		return "\033[34mPUT\033[0m   "
	case "DELETE":
		return "\033[31mDELETE\033[0m"
	case "PATCH":
		return "\033[36mPATCH\033[0m "
	default:
		return method
	}
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
