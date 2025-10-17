package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/EduardoMG12/cine/api_v2/internal/domain"
)

type contextKey string

const UserContextKey contextKey = "user"

func AuthMiddleware(authService domain.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, `{"success":false,"error":{"code":"UNAUTHORIZED","message":"Authorization header required"}}`, http.StatusUnauthorized)
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, `{"success":false,"error":{"code":"INVALID_TOKEN_FORMAT","message":"Authorization header must be Bearer token"}}`, http.StatusUnauthorized)
				return
			}

			token := parts[1]
			user, err := authService.ValidateToken(token)
			if err != nil {
				http.Error(w, `{"success":false,"error":{"code":"INVALID_TOKEN","message":"Invalid or expired token"}}`, http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), UserContextKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetUserFromContext(ctx context.Context) (*domain.User, bool) {
	user, ok := ctx.Value(UserContextKey).(*domain.User)
	return user, ok
}
