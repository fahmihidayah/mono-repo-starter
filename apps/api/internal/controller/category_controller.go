package controller

import (
	"encoding/json"
	"log"
	"net/http"

	categoryRequest "github.com/fahmihidayah/go-api-orchestrator/internal/data/request"
	"github.com/fahmihidayah/go-api-orchestrator/internal/service"
)

// CategoryController handles Category endpoints with React Admin compatibility
type CategoryController struct {
	BaseController
	categoryService service.ICategoryService
}

// CategoryControllerProvider creates a new CategoryControllerReactAdmin
func CategoryControllerProvider(categoryService service.ICategoryService) *CategoryController {
	return &CategoryController{
		categoryService: categoryService,
	}
}

// GetList handles React Admin getList - GET /categories?sort=[...]&range=[...]&filter={...}
// Returns array with X-Total-Count header
// @Summary Get list of categories
// @Description Retrieve a paginated list of categories with optional filtering and sorting. Supports React Admin's getList and getMany operations.
// @Tags categories
// @Accept json
// @Produce json
// @Param sort query string false "Sort parameters in format [field,order], e.g., [\"name\",\"ASC\"]"
// @Param range query string false "Pagination range in format [0,9]"
// @Param filter query string false "Filter parameters as JSON object, e.g., {\"name\":\"Tech\"}. Use {\"ids\":[\"1\",\"2\"]} for getMany"
// @Success 200 {array} domain.Category "List of categories with X-Total-Count header"
// @Failure 400 {object} response.WebResponse "Invalid query parameters"
// @Failure 500 {object} response.WebResponse "Internal server error"
// @Router /api/categories [get]
func (c *CategoryController) GetList(w http.ResponseWriter, r *http.Request) {
	log.Printf("[CategoryController.GetList] Request URL: %s", r.URL.String())
	log.Printf("[CategoryController.GetList] Query String: %s", r.URL.RawQuery)

	params, err := c.ParseQueryParams(r)
	if err != nil {
		log.Printf("[CategoryController.GetList] Error parsing params: %v", err)
		c.SendBadRequest(w, err.Error())
		return
	}

	// Log parsed parameters
	paramsJSON, _ := json.Marshal(params)
	log.Printf("[CategoryController.GetList] Parsed params: %s", string(paramsJSON))

	// Check for getMany case (filter contains "ids")
	if ids, hasIDs := params.GetFilterIDs(); hasIDs {
		log.Printf("[CategoryController.GetList] GetMany detected - IDs: %v", ids)
		c.getMany(w, r, ids)
		return
	}

	log.Printf("[CategoryController.GetList] Pagination - Limit: %d, Offset: %d", params.Limit, params.Offset)
	log.Printf("[CategoryController.GetList] Sort - Field: %s, Order: %s", params.Sort[0], params.Sort[1])

	categories, paginateInfo, err := c.categoryService.GetWithQueryParams(r.Context(), params)
	if err != nil {
		log.Printf("[CategoryController.GetList] Service error: %v", err)
		c.SendInternalError(w, err.Error())
		return
	}

	log.Printf("[CategoryController.GetList] Results - Total: %d, Categories returned: %d", paginateInfo.TotalDocs, len(categories))
	c.SendListWithPagination(w, categories, paginateInfo)
}

// getMany handles React Admin getMany - GET /categories?filter={"ids":[123,124,125]}
func (c *CategoryController) getMany(w http.ResponseWriter, r *http.Request, ids []string) {
	categories, err := c.categoryService.GetByIDs(r.Context(), ids)
	if err != nil {
		c.SendInternalError(w, err.Error())
		return
	}

	// For getMany, don't send total count
	c.SendOne(w, categories)
}

// GetOne handles React Admin getOne - GET /categories/123
// Returns single object
// @Summary Get category by ID
// @Description Retrieve a single category by its unique identifier
// @Tags categories
// @Accept json
// @Produce json
// @Param id path string true "Category ID"
// @Success 200 {object} domain.Category "Category details"
// @Failure 400 {object} response.WebResponse "ID is required"
// @Failure 404 {object} response.WebResponse "Category not found"
// @Router /api/categories/{id} [get]
func (c *CategoryController) GetOne(w http.ResponseWriter, r *http.Request) {
	id := c.GetIDFromURL(r)
	if id == "" {
		c.SendBadRequest(w, "ID is required")
		return
	}

	category, err := c.categoryService.GetByID(r.Context(), id)
	if err != nil {
		c.SendNotFound(w, err.Error())
		return
	}

	c.SendOne(w, category)
}

// Create handles React Admin create - POST /categories
// Returns created object with ID
// @Summary Create a new category
// @Description Create a new category.
// @Tags categories
// @Accept json
// @Produce json
// @Param request body categories.CreateCategoryRequest true "Category creation details"
// @Success 200 {object} domain.Category "Created category details"
// @Failure 400 {object} response.WebResponse "Invalid request body or validation error"
// @Security BearerAuth
// @Router /api/categories [post]
func (c *CategoryController) Create(w http.ResponseWriter, r *http.Request) {
	var req categoryRequest.CreateCategoryRequest
	if err := c.DecodeJSONBody(r, &req); err != nil {
		c.SendBadRequest(w, "Invalid request body")
		return
	}

	log.Printf("[CategoryController.Create] Creating category - Name: %s", req.Title)

	category, err := c.categoryService.Create(r.Context(), &req)
	if err != nil {
		c.SendBadRequest(w, err.Error())
		return
	}

	c.SendOne(w, category)
}

// Update handles React Admin update - PUT /categories/123
// Returns updated object
// @Summary Update category
// @Description Update an existing category by ID.
// @Tags categories
// @Accept json
// @Produce json
// @Param id path string true "Category ID"
// @Param request body categories.CreateCategoryRequest true "Category update details"
// @Success 200 {object} domain.Category "Updated category details"
// @Failure 400 {object} response.WebResponse "Invalid request body or ID required"
// @Security BearerAuth
// @Router /api/categories/{id} [put]
func (c *CategoryController) Update(w http.ResponseWriter, r *http.Request) {
	id := c.GetIDFromURL(r)
	if id == "" {
		c.SendBadRequest(w, "ID is required")
		return
	}

	var req categoryRequest.CreateCategoryRequest
	if err := c.DecodeJSONBody(r, &req); err != nil {
		c.SendBadRequest(w, "Invalid request body")
		return
	}

	category, err := c.categoryService.Update(r.Context(), id, &req)
	if err != nil {
		c.SendBadRequest(w, err.Error())
		return
	}

	c.SendOne(w, category)
}

// UpdateMany handles React Admin updateMany - PUT /categories?filter={"id":[123,124,125]}
// Returns array of updated IDs
// @Summary Update multiple categories
// @Description Update multiple categories with the same data. Supports React Admin's updateMany.
// @Tags categories
// @Accept json
// @Produce json
// @Param filter query string true "Filter parameters with category IDs, e.g., {\"id\":[\"1\",\"2\"]}"
// @Param request body map[string]interface{} true "Update data to apply"
// @Success 200 {array} string "Array of updated category IDs"
// @Failure 400 {object} response.WebResponse "Invalid query or missing id filter"
// @Security BearerAuth
// @Router /api/categories/bulk [put]
func (c *CategoryController) UpdateMany(w http.ResponseWriter, r *http.Request) {
	params, err := c.ParseQueryParams(r)
	if err != nil {
		c.SendBadRequest(w, err.Error())
		return
	}

	ids, hasIDs := params.GetFilterIDs()
	if !hasIDs {
		c.SendBadRequest(w, "Missing id filter")
		return
	}

	// Decode the update data
	var updates map[string]interface{}
	if err := c.DecodeJSONBody(r, &updates); err != nil {
		c.SendBadRequest(w, "Invalid request body")
		return
	}

	updatedIDs, err := c.categoryService.UpdateMany(r.Context(), ids, updates)
	if err != nil {
		c.SendBadRequest(w, err.Error())
		return
	}

	c.SendIDs(w, updatedIDs)
}

// Delete handles React Admin delete - DELETE /categories/123
// Returns deleted object
// @Summary Delete category
// @Description Delete a category by its ID. Returns the deleted object.
// @Tags categories
// @Accept json
// @Produce json
// @Param id path string true "Category ID"
// @Success 200 {object} domain.Category "Deleted category details"
// @Failure 400 {object} response.WebResponse "ID is required"
// @Failure 404 {object} response.WebResponse "Category not found"
// @Security BearerAuth
// @Router /api/categories/{id} [delete]
func (c *CategoryController) Delete(w http.ResponseWriter, r *http.Request) {
	id := c.GetIDFromURL(r)
	if id == "" {
		c.SendBadRequest(w, "ID is required")
		return
	}

	// Get category before deletion (to return it)
	category, err := c.categoryService.GetByID(r.Context(), id)
	if err != nil {
		c.SendNotFound(w, err.Error())
		return
	}

	if err := c.categoryService.Delete(r.Context(), id); err != nil {
		c.SendBadRequest(w, err.Error())
		return
	}

	c.SendOne(w, category)
}

// DeleteMany handles React Admin deleteMany - DELETE /categories?filter={"id":[123,124,125]}
// Returns array of deleted IDs
// @Summary Delete multiple categories
// @Description Delete multiple categories by their IDs. Supports React Admin's deleteMany.
// @Tags categories
// @Accept json
// @Produce json
// @Param filter query string true "Filter parameters with category IDs, e.g., {\"id\":[\"1\",\"2\"]}"
// @Success 200 {array} string "Array of deleted category IDs"
// @Failure 400 {object} response.WebResponse "Invalid query or missing id filter"
// @Security BearerAuth
// @Router /api/categories [delete]
func (c *CategoryController) DeleteMany(w http.ResponseWriter, r *http.Request) {
	params, err := c.ParseQueryParams(r)
	if err != nil {
		c.SendBadRequest(w, err.Error())
		return
	}

	ids, hasIDs := params.GetFilterIDs()
	if !hasIDs {
		c.SendBadRequest(w, "Missing id filter")
		return
	}

	if err := c.categoryService.DeleteAll(r.Context(), ids); err != nil {
		c.SendBadRequest(w, err.Error())
		return
	}

	c.SendIDs(w, ids)
}
