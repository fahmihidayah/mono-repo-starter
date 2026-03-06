package users

// CreateUserRequest represents the request body for user registration
type CreateUserRequest struct {
	BaseUserRequest
	Password string `json:"password" validate:"required,min=8,max=100" binding:"required,min=8" example:"securepassword123"`
}

type BaseUserRequest struct {
	Email string `json:"email" validate:"required,email" binding:"required,email" example:"user@example.com"`
	Name  string `json:"name" validate:"required,min=2,max=100" binding:"required,min=2,max=100" example:"John Doe"`
}

type UpdateUserRequest struct {
	ID string `json:"id" validate:"required,uuid4" binding:"required,uuid4"`
	BaseUserRequest
}

type ChangePasswordRequest struct {
	ID          string `json:"id"`
	OldPassword string `json:"old_password" validate:"required,min=8,max=100" binding:"required,min=8" example:"oldsecurepassword"`
	NewPassword string `json:"new_password" validate:"required,min=8,max=100" binding:"required,min=8" example:"newsecurepassword"`
}
