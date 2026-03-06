package request

// LoginUserRequest represents the request body for user login
type LoginUserRequest struct {
	Email    string `json:"email" validate:"required,email" binding:"required,email" example:"user@example.com"`
	Password string `json:"password" validate:"required,min=8" binding:"required,min=8" example:"securepassword123"`
}
