package response

// PaginationResponse represents a paginated response structure
type PaginationResponse struct {
	Code       int         `json:"code" example:"200"`
	Message    string      `json:"message" example:"Success"`
	Data       interface{} `json:"data"`
	Pagination Pagination  `json:"pagination"`
}

// Pagination represents pagination metadata
type Pagination struct {
	Page       int   `json:"page" example:"1"`
	Limit      int   `json:"limit" example:"10"`
	TotalRows  int64 `json:"total_rows" example:"100"`
	TotalPages int   `json:"total_pages" example:"10"`
}
