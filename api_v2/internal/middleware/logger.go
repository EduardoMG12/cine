package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"
)

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

type Logger struct {
	*slog.Logger
	isDevelopment bool
}

func SetupLogger(environment string) *slog.Logger {
	isDevelopment := environment == "development"
	var handler slog.Handler

	if isDevelopment {
		handler = NewColorHandler(os.Stdout, &slog.HandlerOptions{
			Level:     slog.LevelDebug,
			AddSource: true,
		})
	} else {
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level:     slog.LevelInfo,
			AddSource: false,
		})
	}

	return slog.New(handler)
}

func LoggerMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			ww := &responseWriterWrapper{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
			}

			next.ServeHTTP(ww, r)

			duration := time.Since(start)
			logHTTPRequest(r, ww.statusCode, duration)
		})
	}
}

func logHTTPRequest(r *http.Request, statusCode int, duration time.Duration) {
	var level slog.Level
	switch {
	case statusCode >= 500:
		level = slog.LevelError
	case statusCode >= 400:
		level = slog.LevelWarn
	default:
		level = slog.LevelInfo
	}

	slog.Log(r.Context(), level, "HTTP Request",
		"method", r.Method,
		"path", r.URL.Path,
		"status", statusCode,
		"duration", duration.String(),
		"ip", getClientIP(r),
		"user_agent", r.UserAgent(),
	)
}

func NewLogger(isDevelopment bool) *Logger {
	var handler slog.Handler

	if isDevelopment {

		handler = NewColorHandler(os.Stdout, &slog.HandlerOptions{
			Level:     slog.LevelDebug,
			AddSource: true,
		})
	} else {

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

func (l *Logger) LoggingMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			ww := &responseWriterWrapper{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
			}

			next.ServeHTTP(ww, r)

			duration := time.Since(start)
			l.LogRequest(r, ww.statusCode, duration)
		})
	}
}

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

func (l *Logger) LogStartup(config map[string]interface{}) {
	l.printBanner()

	l.Info("🚀 CineVerse API Starting Up",
		"version", "2.0.0",
		"environment", config["environment"],
		"pid", os.Getpid(),
	)

	for key, value := range config {
		if !l.isSensitive(key) {
			l.Debug("Configuration loaded",
				"key", key,
				"value", value,
			)
		}
	}
}

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

func (l *Logger) LogServerStop() {
	l.Info("🛑 Server shutdown completed")
}

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

type responseWriterWrapper struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriterWrapper) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func getClientIP(r *http.Request) string {

	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		return strings.Split(xff, ",")[0]
	}

	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}

	return strings.Split(r.RemoteAddr, ":")[0]
}

func extractHost(dsn string) string {

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
