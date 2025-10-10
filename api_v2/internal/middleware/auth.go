package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"strings"

	"github.com/EduardoMG12/cine/api_v2/internal/auth"
	"github.com/EduardoMG12/cine/api_v2/internal/domain"
)

type contextKey string

const (
	UserContextKey    contextKey = "user"
	ClaimsContextKey  contextKey = "claims"
	SessionContextKey contextKey = "session"
)

type AuthMiddleware struct {
	jwtManager         *auth.JWTManager
	userSessionService domain.UserSessionService
}

func NewAuthMiddleware(jwtManager *auth.JWTManager, userSessionService domain.UserSessionService) *AuthMiddleware {
	return &AuthMiddleware{
		jwtManager:         jwtManager,
		userSessionService: userSessionService,
	}
}

func (m *AuthMiddleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := extractToken(r)
		if token == "" {
			http.Error(w, "Authorization token required", http.StatusUnauthorized)
			return
		}

		claims, err := m.jwtManager.ValidateToken(token)
		if err != nil {
			slog.Warn("Invalid JWT token", "error", err)
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		session, err := m.userSessionService.ValidateSession(token)
		if err != nil {
			slog.Warn("Session validation failed", "error", err, "user_id", claims.UserID)
			http.Error(w, "Session expired or invalid", http.StatusUnauthorized)
			return
		}

		// Add user context
		ctx := context.WithValue(r.Context(), ClaimsContextKey, claims)
		ctx = context.WithValue(ctx, SessionContextKey, session)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *AuthMiddleware) OptionalAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := extractToken(r)
		if token == "" {
			next.ServeHTTP(w, r)
			return
		}

		claims, err := m.jwtManager.ValidateToken(token)
		if err != nil {
			// Log but continue without auth
			slog.Debug("Optional auth failed", "error", err)
			next.ServeHTTP(w, r)
			return
		}

		// Validate session
		session, err := m.userSessionService.ValidateSession(token)
		if err != nil {
			slog.Debug("Optional session validation failed", "error", err)
			next.ServeHTTP(w, r)
			return
		}

		// Add user context
		ctx := context.WithValue(r.Context(), ClaimsContextKey, claims)
		ctx = context.WithValue(ctx, SessionContextKey, session)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func extractToken(r *http.Request) string {
	// Check Authorization header
	bearerToken := r.Header.Get("Authorization")
	if bearerToken != "" {
		parts := strings.Split(bearerToken, " ")
		if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
			return parts[1]
		}
	}

	// Check query parameter as fallback
	return r.URL.Query().Get("token")
}

// Helper functions to get user info from context
func GetUserClaims(ctx context.Context) (*auth.Claims, bool) {
	claims, ok := ctx.Value(ClaimsContextKey).(*auth.Claims)
	return claims, ok
}

func GetUserSession(ctx context.Context) (*domain.UserSession, bool) {
	session, ok := ctx.Value(SessionContextKey).(*domain.UserSession)
	return session, ok
}

func GetUserID(ctx context.Context) (int, bool) {
	claims, ok := GetUserClaims(ctx)
	if !ok {
		return 0, false
	}
	return claims.UserID, true
}
