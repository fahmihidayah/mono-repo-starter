package controller

import (
	"encoding/json"
	"net/http"

	"github.com/fahmihidayah/go-api-orchestrator/internal/utils"
	"github.com/go-chi/chi/v5"
)

// BaseController provides common controller functionality for React Admin
type BaseController struct{}

// ParseQueryParams parses React Admin query parameters
func (bc *BaseController) ParseQueryParams(r *http.Request) (*utils.QueryParams, error) {
	return utils.ParseQueryListParams(r)
}

// GetIDFromURL extracts ID from URL parameter
func (bc *BaseController) GetIDFromURL(r *http.Request) string {
	return chi.URLParam(r, "id")
}

// DecodeJSONBody decodes JSON request body
func (bc *BaseController) DecodeJSONBody(r *http.Request, v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}

// SendList sends a React Admin list response with pagination
func (bc *BaseController) SendList(w http.ResponseWriter, data interface{}, total int64) {
	utils.SendReactAdminList(w, data, total)
}

func (bc *BaseController) SendListWithPagination(w http.ResponseWriter, data interface{}, params *utils.PaginateInfo) {
	utils.SendReactAdminList(w,
		&map[string]interface{}{
			"message": "success",
			"code":    http.StatusOK,
			"data":    data,
			"pagination": map[string]interface{}{
				"limit":      params.Limit,
				"page":       params.Page,
				"totalPages": params.TotalPages,
				"nextPage":   params.NextPage,
				"prevPage":   params.PrevPage,
				"totalDocs":  params.TotalDocs,
			},
		}, params.TotalDocs)
}

// SendOne sends a single record response
func (bc *BaseController) SendOne(w http.ResponseWriter, data interface{}) {
	utils.SendReactAdminOne(w,
		&map[string]interface{}{
			"message": "success",
			"code":    http.StatusOK,
			"data":    data,
		})
}

// SendIDs sends an array of IDs (for bulk operations)
func (bc *BaseController) SendIDs(w http.ResponseWriter, ids []string) {
	utils.SendReactAdminIDs(w, &map[string]interface{}{
		"message": "success",
		"code":    http.StatusOK,
		"data":    ids,
	})
}

// SendError sends an error response
func (bc *BaseController) SendError(w http.ResponseWriter, message string, status int) {
	utils.SendReactAdminError(w,
		&map[string]interface{}{
			"message": message,
			"code":    status,
		}, status)
}

// SendBadRequest sends a 400 Bad Request error
func (bc *BaseController) SendBadRequest(w http.ResponseWriter, message string) {
	bc.SendError(w, message, http.StatusBadRequest)
}

// SendNotFound sends a 404 Not Found error
func (bc *BaseController) SendNotFound(w http.ResponseWriter, message string) {
	bc.SendError(w, message, http.StatusNotFound)
}

// SendInternalError sends a 500 Internal Server Error
func (bc *BaseController) SendInternalError(w http.ResponseWriter, message string) {
	bc.SendError(w, message, http.StatusInternalServerError)
}

// SendUnauthorized sends a 401 Unauthorized error
func (bc *BaseController) SendUnauthorized(w http.ResponseWriter, message string) {
	bc.SendError(w, message, http.StatusUnauthorized)
}

// HasFilterParam checks if request has filter parameter (for bulk operations)
func (bc *BaseController) HasFilterParam(r *http.Request) bool {
	return r.URL.Query().Get("filter") != ""
}
