package handler

import (
	"encoding/json"
	"net/http"

	"github.com/fahmihidayah/go-api-orchestrator/internal/data/response"
)

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)

	json.NewEncoder(w).Encode(response.ErrorResponse{
		Status:  http.StatusNotFound,
		Message: "Route Not Found",
		Path:    r.URL.Path,
	})
}
