package utils

import (
	"encoding/json"
	"net/http"

	"github.com/fahmihidayah/go-api-orchestrator/internal/data/response"
)

// SendJSON sends a JSON response with the given status code, message, and data
func SendJSON(w http.ResponseWriter, statusCode int, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response.WebResponse{
		Code:    statusCode,
		Message: message,
		Data:    data,
	})
}

// SendSuccess sends a successful JSON response (200 OK)
func SendSuccess(w http.ResponseWriter, message string, data interface{}) {
	SendJSON(w, http.StatusOK, message, data)
}

// SendCreated sends a created JSON response (201 Created)
func SendCreated(w http.ResponseWriter, message string, data interface{}) {
	SendJSON(w, http.StatusCreated, message, data)
}

// SendBadRequest sends a bad request error response (400 Bad Request)
func SendBadRequest(w http.ResponseWriter, message string) {
	SendJSON(w, http.StatusBadRequest, message, nil)
}

// SendUnauthorized sends an unauthorized error response (401 Unauthorized)
func SendUnauthorized(w http.ResponseWriter, message string) {
	SendJSON(w, http.StatusUnauthorized, message, nil)
}

// SendNotFound sends a not found error response (404 Not Found)
func SendNotFound(w http.ResponseWriter, message string) {
	SendJSON(w, http.StatusNotFound, message, nil)
}

// SendInternalError sends an internal server error response (500 Internal Server Error)
func SendInternalError(w http.ResponseWriter, message string) {
	SendJSON(w, http.StatusInternalServerError, message, nil)
}

// SendPaginated sends a paginated JSON response with pagination metadata
func SendPaginated(w http.ResponseWriter, message string, data interface{}, page, limit int, totalRows int64, totalPages int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response.PaginationResponse{
		Code:    http.StatusOK,
		Message: message,
		Data:    data,
		Pagination: response.Pagination{
			Page:       page,
			Limit:      limit,
			TotalRows:  totalRows,
			TotalPages: totalPages,
		},
	})
}

// Deprecated: Use SendSuccess instead
func HandleSuccess(w http.ResponseWriter, data interface{}, message string, err error) {
	SendSuccess(w, message, data)
}

// Deprecated: Use SendBadRequest instead
func HandleBadRequest(w http.ResponseWriter, message string, err error) {
	SendBadRequest(w, message)
}

// Deprecated: Use SendJSON instead
func HandleError(w http.ResponseWriter, status int, message string, err error) {
	SendJSON(w, status, message, nil)
}
