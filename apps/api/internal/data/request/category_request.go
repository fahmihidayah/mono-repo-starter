package request

// CreateCategoryRequest represents the request body for creating a category
type CreateCategoryRequest struct {
	BaseCategoryRequest
}

// BaseCategoryRequest contains the common fields for category operations
type BaseCategoryRequest struct {
	Title string `json:"title" validate:"required,min=3,max=255" binding:"required,min=3,max=255" example:"Technology"`
}

// UpdateCategoryRequest represents the request body for updating a category
type UpdateCategoryRequest struct {
	ID string `json:"id" validate:"required,uuid4" binding:"required,uuid4"`
	BaseCategoryRequest
}
