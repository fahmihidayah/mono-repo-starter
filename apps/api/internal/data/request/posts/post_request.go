package posts

// CreatePostRequest represents the request body for creating a post
type CreatePostRequest struct {
	BasePostRequest
}

// BasePostRequest contains the common fields for post operations
type BasePostRequest struct {
	Title       string   `json:"title" validate:"required,min=3,max=255" binding:"required,min=3,max=255" example:"My First Post"`
	Content     string   `json:"content" validate:"required" binding:"required" example:"This is the post content"`
	CategoryIDs []string `json:"category_ids" example:"[\"cat-id-1\",\"cat-id-2\"]"`
}

// UpdatePostRequest represents the request body for updating a post
type UpdatePostRequest struct {
	ID string `json:"id" validate:"required,uuid4" binding:"required,uuid4"`
	BasePostRequest
}
