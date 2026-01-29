package middleware

import (
	"context"
	"encoding/json"
	"log"
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
			log.Printf("[AuthMiddleware] Processing request: %s %s", r.Method, r.URL.Path)

			// Get Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				log.Printf("[AuthMiddleware] Missing Authorization header for request: %s %s", r.Method, r.URL.Path)
				respondWithError(w, http.StatusUnauthorized, "Authorization header required")
				return
			}

			// Check if it's a Bearer token
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				log.Printf("[AuthMiddleware] Invalid authorization header format for request: %s %s", r.Method, r.URL.Path)
				respondWithError(w, http.StatusUnauthorized, "Invalid authorization header format. Expected 'Bearer {token}'")
				return
			}

			token := parts[1]
			if token == "" {
				log.Printf("[AuthMiddleware] Empty token received for request: %s %s", r.Method, r.URL.Path)
				respondWithError(w, http.StatusUnauthorized, "Token is empty")
				return
			}

			// Mask token for security (show only first 10 and last 4 characters)
			tokenPreview := token
			if len(token) > 14 {
				tokenPreview = token[:10] + "..." + token[len(token)-4:]
			}
			log.Printf("[AuthMiddleware] Validating token: %s for request: %s %s", tokenPreview, r.Method, r.URL.Path)

			// Check if token is blacklisted
			isBlacklisted, err := tokenBlacklistRepo.IsBlacklisted(token)
			if err != nil {
				log.Printf("[AuthMiddleware] Error checking token blacklist status: %v", err)
				respondWithError(w, http.StatusInternalServerError, "Error checking token status")
				return
			}
			if isBlacklisted {
				log.Printf("[AuthMiddleware] Blacklisted token attempted for request: %s %s", r.Method, r.URL.Path)
				respondWithError(w, http.StatusUnauthorized, "Token has been revoked")
				return
			}

			// Validate JWT token
			claims, err := security.ValidateJWT(token, cfg.JWTSecret)
			if err != nil {
				if err == security.ErrExpiredToken {
					log.Printf("[AuthMiddleware] Expired token for request: %s %s", r.Method, r.URL.Path)
					respondWithError(w, http.StatusUnauthorized, "Token has expired")
					return
				}
				log.Printf("[AuthMiddleware] Invalid token for request: %s %s - Error: %v", r.Method, r.URL.Path, err)
				respondWithError(w, http.StatusUnauthorized, "Invalid token")
				return
			}

			log.Printf("[AuthMiddleware] Authentication successful for user ID: %s - Request: %s %s", claims.UserID, r.Method, r.URL.Path)

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
