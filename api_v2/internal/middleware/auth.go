package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/EduardoMG12/cine/api_v2/internal/domain"
	"github.com/EduardoMG12/cine/api_v2/internal/dto"
	"github.com/EduardoMG12/cine/api_v2/internal/infrastructure"
	"github.com/google/uuid"
)

type contextKey string

const UserContextKey contextKey = "user"

func JWTAuthMiddleware(jwtService *infrastructure.JWTService, userRepo domain.UserRepository) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				sendErrorResponse(w, http.StatusUnauthorized, "UNAUTHORIZED", "Authorization header required")
				return
			}

			tokenParts := strings.Split(authHeader, " ")
			if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
				sendErrorResponse(w, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid authorization header format")
				return
			}

			token := tokenParts[1]
			userID, err := jwtService.ValidateToken(token)
			if err != nil {
				sendErrorResponse(w, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid or expired token")
				return
			}

			user, err := userRepo.GetUserByID(userID)
			if err != nil {
				sendErrorResponse(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not found")
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

func GetUserIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	user, ok := GetUserFromContext(ctx)
	if !ok {
		return uuid.Nil, false
	}
	return user.ID, true
}

func sendErrorResponse(w http.ResponseWriter, statusCode int, code, message string) {
	response := dto.APIResponse{
		Success: false,
		Error: &dto.APIError{
			Code:    code,
			Message: message,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}
