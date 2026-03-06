package users

// FilterUserRequest represents the query parameters for filtering users
type FilterUserRequest struct {
	Name  string `json:"name" example:"John"`
	Email string `json:"email" example:"john@example.com"`
	Page  int    `json:"page" example:"1"`
	Limit int    `json:"limit" example:"10"`
}
