package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

// QueryParams represents query parameters for React Admin list operations
type QueryParams struct {
	Sort   []string               // ["field", "order"] where order is "ASC" or "DESC"
	Range  []int                  // [start, end] for pagination (0-indexed, inclusive)
	Filter map[string]interface{} // Filter conditions as map
	Limit  int
	Offset int
	Where  map[string]string
}

// ParseQueryListParams parses React Admin query parameters from the request
// React Admin sends: sort=["title","ASC"]&range=[0,24]&filter={"author_id":12}
func ParseQueryListParams(r *http.Request) (*QueryParams, error) {
	params := &QueryParams{
		Sort:   []string{"id", "ASC"},
		Range:  []int{0, 24},
		Filter: make(map[string]interface{}),
	}

	// Parse sort parameter
	if sortStr := r.URL.Query().Get("sort"); sortStr != "" {
		if err := json.Unmarshal([]byte(sortStr), &params.Sort); err != nil {
			return nil, fmt.Errorf("invalid sort parameter: %w", err)
		}
		if len(params.Sort) != 2 {
			return nil, fmt.Errorf("sort parameter must have exactly 2 elements")
		}
	}

	// Parse range parameter
	if rangeStr := r.URL.Query().Get("range"); rangeStr != "" {
		if err := json.Unmarshal([]byte(rangeStr), &params.Range); err != nil {
			return nil, fmt.Errorf("invalid range parameter: %w", err)
		}
		if len(params.Range) != 2 {
			return nil, fmt.Errorf("range parameter must have exactly 2 elements")
		}
	}
	limit, offset := params.CalculatePagination()
	params.Limit = limit
	params.Offset = offset
	// Parse filter parameter
	if filterStr := r.URL.Query().Get("filter"); filterStr != "" {
		if err := json.Unmarshal([]byte(filterStr), &params.Filter); err != nil {
			return nil, fmt.Errorf("invalid filter parameter: %w", err)
		}
	}

	// log.Printf("[QueryParams] Filters: %s", string(r.URL.Query().Get("filter")))

	// filters := make(map[string]string)
	for key, value := range params.Filter {
		log.Printf("[QueryParams] filter content : %s - %s", key, value)
	}

	return params, nil
}

// CalculatePagination calculates limit and offset from React Admin range
func (p *QueryParams) CalculatePagination() (limit, offset int) {
	offset = p.Range[0]
	limit = p.Range[1] - p.Range[0] + 1
	return
}

// GetSortField returns the field to sort by
func (p *QueryParams) GetSortField() string {
	if len(p.Sort) > 0 {
		return p.Sort[0]
	}
	return "id"
}

// GetSortOrder returns the sort order (ASC or DESC)
func (p *QueryParams) GetSortOrder() string {
	if len(p.Sort) > 1 {
		return p.Sort[1]
	}
	return "ASC"
}

// GetFilterValue returns a filter value by key
func (p *QueryParams) GetFilterValue(key string) (interface{}, bool) {
	value, exists := p.Filter[key]
	return value, exists
}

// GetFilterString returns a filter value as string
func (p *QueryParams) GetFilterString(key string) (string, bool) {
	value, exists := p.Filter[key]
	if !exists {
		return "", false
	}
	if str, ok := value.(string); ok {
		return str, true
	}
	return "", false
}

// GetFilterIDs returns the "ids" or "id" filter as array of strings
// React Admin uses "ids" for getMany and "id" for updateMany/deleteMany
func (p *QueryParams) GetFilterIDs() ([]string, bool) {
	// Try "ids" first (used by getMany)
	if idsValue, exists := p.Filter["ids"]; exists {
		return convertToStringArray(idsValue)
	}

	// Try "id" (used by updateMany/deleteMany)
	if idValue, exists := p.Filter["id"]; exists {
		return convertToStringArray(idValue)
	}

	return nil, false
}

// convertToStringArray converts interface{} to []string
// Handles []interface{}, []string, []int, []float64
func convertToStringArray(value interface{}) ([]string, bool) {
	switch v := value.(type) {
	case []interface{}:
		result := make([]string, len(v))
		for i, item := range v {
			result[i] = fmt.Sprintf("%v", item)
		}
		return result, true
	case []string:
		return v, true
	case []int:
		result := make([]string, len(v))
		for i, item := range v {
			result[i] = strconv.Itoa(item)
		}
		return result, true
	case []float64:
		result := make([]string, len(v))
		for i, item := range v {
			result[i] = fmt.Sprintf("%.0f", item)
		}
		return result, true
	}
	return nil, false
}

// SendReactAdminList sends a React Admin list response with X-Total-Count header
func SendReactAdminList(w http.ResponseWriter, data interface{}, total int64) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Total-Count", strconv.FormatInt(total, 10))
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}

// SendReactAdminOne sends a single record response (plain object)
func SendReactAdminOne(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}

// SendReactAdminIDs sends an array of IDs (for updateMany/deleteMany)
func SendReactAdminIDs(w http.ResponseWriter, ids []string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ids)
}

// SendReactAdminError sends an error response
func SendReactAdminError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": message,
		"status":  status,
	})
}
