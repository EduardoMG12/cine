package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"
)

// Colors for console output
const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Purple = "\033[35m"
	Cyan   = "\033[36m"
	Gray   = "\033[37m"
	Bold   = "\033[1m"
)

// Logger wraps slog.Logger with additional functionality
type Logger struct {
	*slog.Logger
	isDevelopment bool
}

// NewLogger creates a new logger instance
func NewLogger(isDevelopment bool) *Logger {
	var handler slog.Handler

	if isDevelopment {
		// Development: colorized text handler
		handler = NewColorHandler(os.Stdout, &slog.HandlerOptions{
			Level:     slog.LevelDebug,
			AddSource: true,
		})
	} else {
		// Production: JSON handler
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level:     slog.LevelInfo,
			AddSource: false,
		})
	}

	logger := slog.New(handler)
	return &Logger{
		Logger:        logger,
		isDevelopment: isDevelopment,
	}
}

// LoggingMiddleware returns HTTP logging middleware
func (l *Logger) LoggingMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Wrap response writer to capture status code
			ww := &responseWriterWrapper{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
			}

			// Process request
			next.ServeHTTP(ww, r)

			// Log request
			duration := time.Since(start)
			l.LogRequest(r, ww.statusCode, duration)
		})
	}
}

// LogRequest logs HTTP request details
func (l *Logger) LogRequest(r *http.Request, statusCode int, duration time.Duration) {
	var level slog.Level
	switch {
	case statusCode >= 500:
		level = slog.LevelError
	case statusCode >= 400:
		level = slog.LevelWarn
	default:
		level = slog.LevelInfo
	}

	l.Log(r.Context(), level, "HTTP Request",
		"method", r.Method,
		"path", r.URL.Path,
		"status", statusCode,
		"duration", duration.String(),
		"ip", getClientIP(r),
		"user_agent", r.UserAgent(),
	)
}

// LogStartup logs application startup information
func (l *Logger) LogStartup(config map[string]interface{}) {
	l.printBanner()

	l.Info("🚀 CineVerse API Starting Up",
		"version", "2.0.0",
		"environment", config["environment"],
		"pid", os.Getpid(),
	)

	// Log configuration (without secrets)
	for key, value := range config {
		if !l.isSensitive(key) {
			l.Debug("Configuration loaded",
				"key", key,
				"value", value,
			)
		}
	}
}

// LogDatabaseConnection logs database connection status
func (l *Logger) LogDatabaseConnection(connected bool, dsn string) {
	if connected {
		l.Info("✅ Database connected successfully",
			"driver", "postgresql",
			"host", extractHost(dsn),
		)
	} else {
		l.Error("❌ Database connection failed",
			"driver", "postgresql",
			"dsn", redactDSN(dsn),
		)
	}
}

// LogServerStart logs server startup
func (l *Logger) LogServerStart(address string) {
	l.Info("🌐 HTTP Server started",
		"address", address,
		"timestamp", time.Now().Format(time.RFC3339),
	)

	if l.isDevelopment {
		fmt.Printf("\n%s%s🎬 CineVerse API is running!%s\n", Bold, Green, Reset)
		fmt.Printf("%s📡 Server: %shttp://%s%s\n", Cyan, Bold, address, Reset)
		fmt.Printf("%s📚 Health: %shttp://%s/health%s\n", Cyan, Bold, address, Reset)
		fmt.Printf("%s📖 Docs: %shttp://%s/docs%s\n", Cyan, Bold, address, Reset)
		fmt.Printf("%s⚡ Environment: %s%s%s\n\n", Yellow, Bold, "development", Reset)
	}
}

// LogServerStop logs server shutdown
func (l *Logger) LogServerStop() {
	l.Info("🛑 Server shutdown completed")
}

// LogMigration logs database migration status
func (l *Logger) LogMigration(success bool, count int) {
	if success {
		l.Info("✅ Database migrations completed",
			"migrations_applied", count,
		)
	} else {
		l.Error("❌ Database migrations failed",
			"migrations_attempted", count,
		)
	}
}

// printBanner prints ASCII art banner in development
func (l *Logger) printBanner() {
	if !l.isDevelopment {
		return
	}

	banner := fmt.Sprintf(`
%s%s
 ██████╗██╗███╗   ██╗███████╗██╗   ██╗███████╗██████╗ ███████╗███████╗
██╔════╝██║████╗  ██║██╔════╝██║   ██║██╔════╝██╔══██╗██╔════╝██╔════╝
██║     ██║██╔██╗ ██║█████╗  ██║   ██║█████╗  ██████╔╝███████╗█████╗  
██║     ██║██║╚██╗██║██╔══╝  ╚██╗ ██╔╝██╔══╝  ██╔══██╗╚════██║██╔══╝  
╚██████╗██║██║ ╚████║███████╗ ╚████╔╝ ███████╗██║  ██║███████║███████╗
 ╚═════╝╚═╝╚═╝  ╚═══╝╚══════╝  ╚═══╝  ╚══════╝╚═╝  ╚═╝╚══════╝╚══════╝
%s
%s                    🎬 Movie Social Network API v2.0                     %s
%s                        Clean Architecture | Sprint 0                    %s
`, Bold, Blue, Reset, Cyan, Reset, Gray, Reset)

	fmt.Print(banner)
}

// isSensitive checks if a config key contains sensitive information
func (l *Logger) isSensitive(key string) bool {
	sensitiveKeys := []string{
		"password", "secret", "key", "token", "api_key",
		"jwt_secret", "db_password", "redis_password", "tmdb_api_key",
	}

	keyLower := strings.ToLower(key)
	for _, sensitive := range sensitiveKeys {
		if strings.Contains(keyLower, sensitive) {
			return true
		}
	}
	return false
}

// Helper functions

type responseWriterWrapper struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriterWrapper) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		return strings.Split(xff, ",")[0]
	}

	// Check X-Real-IP header
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}

	// Fall back to remote address
	return strings.Split(r.RemoteAddr, ":")[0]
}

func extractHost(dsn string) string {
	// Simple extraction for logging purposes
	if strings.Contains(dsn, "host=") {
		parts := strings.Split(dsn, "host=")
		if len(parts) > 1 {
			hostPart := strings.Split(parts[1], " ")[0]
			return hostPart
		}
	}
	return "unknown"
}

func redactDSN(dsn string) string {
	// Replace password in DSN for logging
	if strings.Contains(dsn, "password=") {
		parts := strings.Split(dsn, "password=")
		if len(parts) == 2 {
			afterPassword := strings.Split(parts[1], " ")
			afterPassword[0] = "***"
			return parts[0] + "password=" + strings.Join(afterPassword, " ")
		}
	}
	return dsn
}
