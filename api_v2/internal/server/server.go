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

	"github.com/EduardoMG12/cine/api_v2/internal/config"
	httpHandler "github.com/EduardoMG12/cine/api_v2/internal/handler/http"
	"github.com/EduardoMG12/cine/api_v2/internal/infrastructure"
	customMiddleware "github.com/EduardoMG12/cine/api_v2/internal/middleware"
	"github.com/EduardoMG12/cine/api_v2/internal/repository"
	"github.com/EduardoMG12/cine/api_v2/internal/usecase/auth"
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
	db, err := sqlx.Connect("postgres", cfg.Database.GetDSN())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	db.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	db.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	db.SetConnMaxLifetime(time.Duration(cfg.Database.ConnMaxLifetime) * time.Minute)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	server := &Server{
		config: cfg,
		db:     db,
		logger: logger,
	}

	server.setupHTTPServer()
	return server, nil
}

func (s *Server) setupHTTPServer() {
	r := chi.NewRouter()
	s.router = r

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

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

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy","service":"cineverse-api"}`))
	})

	passwordService := infrastructure.NewPasswordService()
	jwtService := infrastructure.NewJWTService(s.config.JWT.Secret)

	userRepo := repository.NewUserRepository(s.db)
	sessionRepo := repository.NewSessionRepository(s.db)

	registerUC := auth.NewRegisterUseCase(userRepo, sessionRepo, passwordService, jwtService)
	loginUC := auth.NewLoginUseCase(userRepo, sessionRepo, passwordService, jwtService)
	getMeUC := auth.NewGetMeUseCase(userRepo)
	logoutUC := auth.NewLogoutUseCase(sessionRepo)
	logoutAllUC := auth.NewLogoutAllUseCase(sessionRepo)

	authHandler := httpHandler.NewAuthHandler(registerUC, loginUC, getMeUC, logoutUC, logoutAllUC)

	authMiddleware := customMiddleware.JWTAuthMiddleware(jwtService, userRepo)

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", authHandler.Register)
			r.Post("/login", authHandler.Login)

			r.Group(func(r chi.Router) {
				r.Use(authMiddleware)
				r.Get("/me", authHandler.GetMe)
				r.Post("/logout", authHandler.Logout)
				r.Post("/logout-all", authHandler.LogoutAll)
			})
		})
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

	s.printRoutes()
	return s.httpServer.ListenAndServe()
}

func (s *Server) printRoutes() {
	routes := s.getRoutes()

	fmt.Println("\n┌─────────────────────────────────────────────────────────────────┐")
	fmt.Println("│                     📡 REGISTERED ROUTES                        │")
	fmt.Println("└─────────────────────────────────────────────────────────────────┘")
	fmt.Println()

	publicRoutes := []RouteInfo{}
	protectedRoutes := []RouteInfo{}

	for _, route := range routes {
		if strings.Contains(route.Path, "/auth/me") ||
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

	if s.db != nil {
		if err := s.db.Close(); err != nil {
			s.logger.Error("Failed to close database connection", "error", err)
		}
	}

	return s.httpServer.Shutdown(ctx)
}
