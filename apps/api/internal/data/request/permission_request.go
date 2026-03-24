package request

// CreatePermissionRequest represents the request body for creating a permission
type CreatePermissionRequest struct {
	BasePermissionRequest
}

// BasePermissionRequest contains the common fields for permission operations
type BasePermissionRequest struct {
	Permission string `json:"permission" validate:"required,min=3,max=255,permission_format" binding:"required,min=3,max=255" example:"users:read"`
}

// UpdatePermissionRequest represents the request body for updating a permission
type UpdatePermissionRequest struct {
	ID string `json:"id" validate:"required,uuid4" binding:"required,uuid4"`
	BasePermissionRequest
}
