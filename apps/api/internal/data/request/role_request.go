package request

// CreateRoleRequest represents the request body for creating a role
type CreateRoleRequest struct {
	BaseRoleRequest
}

// BaseRoleRequest contains the common fields for role operations
type BaseRoleRequest struct {
	Name string `json:"name" validate:"required,min=3,max=255" binding:"required,min=3,max=255" example:"Admin"`
}

// UpdateRoleRequest represents the request body for updating a role
type UpdateRoleRequest struct {
	ID string `json:"id" validate:"required,uuid4" binding:"required,uuid4"`
	BaseRoleRequest
}
