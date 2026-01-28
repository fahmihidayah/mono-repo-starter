package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/fahmihidayah/go-api-orchestrator/internal/config"
	"github.com/fahmihidayah/go-api-orchestrator/internal/data/response"
	"github.com/fahmihidayah/go-api-orchestrator/internal/repository"
	"github.com/fahmihidayah/go-api-orchestrator/internal/security"
)

type contextKey string

const (
	UserContextKey   contextKey = "user"
	UserIDContextKey contextKey = "user_id"
)

// AuthMiddleware validates JWT token from Authorization header
func AuthMiddleware(cfg *config.Config, tokenBlacklistRepo repository.ITokenBlacklistRepository) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				respondWithError(w, http.StatusUnauthorized, "Authorization header required")
				return
			}

			// Check if it's a Bearer token
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				respondWithError(w, http.StatusUnauthorized, "Invalid authorization header format. Expected 'Bearer {token}'")
				return
			}

			token := parts[1]
			if token == "" {
				respondWithError(w, http.StatusUnauthorized, "Token is empty")
				return
			}

			// Check if token is blacklisted
			isBlacklisted, err := tokenBlacklistRepo.IsBlacklisted(token)
			if err != nil {
				respondWithError(w, http.StatusInternalServerError, "Error checking token status")
				return
			}
			if isBlacklisted {
				respondWithError(w, http.StatusUnauthorized, "Token has been revoked")
				return
			}

			// Validate JWT token
			claims, err := security.ValidateJWT(token, cfg.JWTSecret)
			if err != nil {
				if err == security.ErrExpiredToken {
					respondWithError(w, http.StatusUnauthorized, "Token has expired")
					return
				}
				respondWithError(w, http.StatusUnauthorized, "Invalid token")
				return
			}

			// Add user claims to request context
			// Store full claims under "user" key
			ctx := context.WithValue(r.Context(), UserContextKey, claims)
			// Store user_id separately for easy access in controllers
			ctx = context.WithValue(ctx, UserIDContextKey, claims.UserID)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserFromContext retrieves user claims from request context
func GetUserFromContext(ctx context.Context) (*security.JWTClaims, bool) {
	claims, ok := ctx.Value(UserContextKey).(*security.JWTClaims)
	return claims, ok
}

// GetUserIDFromContext retrieves only the user ID from request context
func GetUserIDFromContext(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(UserIDContextKey).(string)
	return userID, ok
}

// respondWithError sends a JSON error response
func respondWithError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(
		response.ErrorResponse{Status: code, Message: message, Path: ""},
	)
}
