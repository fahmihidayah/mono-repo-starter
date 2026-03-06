package response

import "github.com/fahmihidayah/go-api-orchestrator/internal/domain"

// UserResponse represents the response containing user data and authentication token
type UserResponse struct {
	Token string      `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	User  domain.User `json:"user"`
	Exp   int64       `json:"exp" example:"1701369600"` // Expiration time as Unix timestamp
}
