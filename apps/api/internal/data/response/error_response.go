package response

type ErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Path    string `json:"path,omitempty"`
}
