package middleware

import (
	"net/http"
	"strings"

	"github.com/fahmihidayah/go-api-orchestrator/internal/config"
)

//	r.Use(cors.Handler(cors.Options{
//		AllowedOrigins:   []string{"https://*", "http://*"},
//		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
//		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
//		AllowCredentials: true,
//		MaxAge:           300,
//	}))
func CorsMiddleware(cfg *config.Config) func(http.Handler) http.Handler {
	allowedOrigins := strings.Split(cfg.CORSAllowedOrigins, ",")

	// Validate and sanitize allowed origins - reject wildcards
	validOrigins := make([]string, 0, len(allowedOrigins))
	for _, origin := range allowedOrigins {
		trimmed := strings.TrimSpace(origin)
		// Reject wildcard origins for security
		if trimmed != "" && trimmed != "*" {
			validOrigins = append(validOrigins, trimmed)
		}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			origin := r.Header.Get("Origin")
			isValidOrigin := false

			// Only allow explicitly whitelisted origins
			if origin != "" && origin != "*" {
				for _, allowedOrigin := range validOrigins {
					if allowedOrigin == origin {
						isValidOrigin = true
						w.Header().Set("Access-Control-Allow-Origin", origin)
						break
					}
				}
			}

			// Only set CORS headers if origin is valid
			if isValidOrigin {
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
				w.Header().Set("Access-Control-Allow-Credentials", "true")
				// Expose headers required by React Admin
				w.Header().Set("Access-Control-Expose-Headers", "Content-Range, X-Total-Count")
			}

			// Handle preflight requests
			if r.Method == http.MethodOptions {
				if isValidOrigin {
					w.WriteHeader(http.StatusNoContent)
				} else {
					// Reject preflight from invalid origins
					w.WriteHeader(http.StatusForbidden)
				}
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
