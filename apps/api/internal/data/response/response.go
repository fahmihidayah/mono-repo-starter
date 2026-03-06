package response

// WebResponse represents the standard API response structure
type WebResponse struct {
	Code    int         `json:"code" example:"200"`
	Message string      `json:"message" example:"Success"`
	Data    interface{} `json:"data,omitempty"`
}
