package controller

import (
	"encoding/json"
	"log"
	"net/http"

	permissionRequest "github.com/fahmihidayah/go-api-orchestrator/internal/data/request"
	"github.com/fahmihidayah/go-api-orchestrator/internal/service"
)

// PermissionController handles Permission endpoints with React Admin compatibility
type PermissionController struct {
	BaseController
	permissionService service.IPermissionService
}

// PermissionControllerProvider creates a new PermissionController
func PermissionControllerProvider(permissionService service.IPermissionService) *PermissionController {
	return &PermissionController{
		permissionService: permissionService,
	}
}

// GetList handles React Admin getList - GET /permissions?sort=[...]&range=[...]&filter={...}
// Returns array with X-Total-Count header
// @Summary Get list of permissions
// @Description Retrieve a paginated list of permissions with optional filtering and sorting. Supports React Admin's getList and getMany operations.
// @Tags permissions
// @Accept json
// @Produce json
// @Param sort query string false "Sort parameters in format [field,order], e.g., [\"permission\",\"ASC\"]"
// @Param range query string false "Pagination range in format [0,9]"
// @Param filter query string false "Filter parameters as JSON object, e.g., {\"permission\":\"users:read\"}. Use {\"ids\":[\"1\",\"2\"]} for getMany"
// @Success 200 {array} domain.Permission "List of permissions with X-Total-Count header"
// @Failure 400 {object} response.WebResponse "Invalid query parameters"
// @Failure 500 {object} response.WebResponse "Internal server error"
// @Router /api/permissions [get]
func (c *PermissionController) GetList(w http.ResponseWriter, r *http.Request) {
	log.Printf("[PermissionController.GetList] Request URL: %s", r.URL.String())
	log.Printf("[PermissionController.GetList] Query String: %s", r.URL.RawQuery)

	params, err := c.ParseQueryParams(r)
	if err != nil {
		log.Printf("[PermissionController.GetList] Error parsing params: %v", err)
		c.SendBadRequest(w, err.Error())
		return
	}

	// Log parsed parameters
	paramsJSON, _ := json.Marshal(params)
	log.Printf("[PermissionController.GetList] Parsed params: %s", string(paramsJSON))

	// Check for getMany case (filter contains "ids")
	if ids, hasIDs := params.GetFilterIDs(); hasIDs {
		log.Printf("[PermissionController.GetList] GetMany detected - IDs: %v", ids)
		c.getMany(w, r, ids)
		return
	}

	log.Printf("[PermissionController.GetList] Pagination - Limit: %d, Offset: %d", params.Limit, params.Offset)
	log.Printf("[PermissionController.GetList] Sort - Field: %s, Order: %s", params.Sort[0], params.Sort[1])

	permissions, paginateInfo, err := c.permissionService.GetWithQueryParams(r.Context(), params)
	if err != nil {
		log.Printf("[PermissionController.GetList] Service error: %v", err)
		c.SendInternalError(w, err.Error())
		return
	}

	log.Printf("[PermissionController.GetList] Results - Total: %d, Permissions returned: %d", paginateInfo.TotalDocs, len(permissions))
	c.SendListWithPagination(w, permissions, paginateInfo)
}

// getMany handles React Admin getMany - GET /permissions?filter={"ids":[123,124,125]}
func (c *PermissionController) getMany(w http.ResponseWriter, r *http.Request, ids []string) {
	permissions, err := c.permissionService.GetByIDs(r.Context(), ids)
	if err != nil {
		c.SendInternalError(w, err.Error())
		return
	}

	// For getMany, don't send total count
	c.SendOne(w, permissions)
}

// GetOne handles React Admin getOne - GET /permissions/123
// Returns single object
// @Summary Get permission by ID
// @Description Retrieve a single permission by its unique identifier
// @Tags permissions
// @Accept json
// @Produce json
// @Param id path string true "Permission ID"
// @Success 200 {object} domain.Permission "Permission details"
// @Failure 400 {object} response.WebResponse "ID is required"
// @Failure 404 {object} response.WebResponse "Permission not found"
// @Router /api/permissions/{id} [get]
func (c *PermissionController) GetOne(w http.ResponseWriter, r *http.Request) {
	id := c.GetIDFromURL(r)
	if id == "" {
		c.SendBadRequest(w, "ID is required")
		return
	}

	permission, err := c.permissionService.GetByID(r.Context(), id)
	if err != nil {
		c.SendNotFound(w, err.Error())
		return
	}

	c.SendOne(w, permission)
}

// Create handles React Admin create - POST /permissions
// Returns created object with ID
// @Summary Create a new permission
// @Description Create a new permission. Format: [table]:[read/create/update/delete]
// @Tags permissions
// @Accept json
// @Produce json
// @Param request body permissions.CreatePermissionRequest true "Permission creation details"
// @Success 200 {object} domain.Permission "Created permission details"
// @Failure 400 {object} response.WebResponse "Invalid request body or validation error"
// @Security BearerAuth
// @Router /api/permissions [post]
func (c *PermissionController) Create(w http.ResponseWriter, r *http.Request) {
	var req permissionRequest.CreatePermissionRequest
	if err := c.DecodeJSONBody(r, &req); err != nil {
		c.SendBadRequest(w, "Invalid request body")
		return
	}

	log.Printf("[PermissionController.Create] Creating permission - Permission: %s", req.Permission)

	permission, err := c.permissionService.Create(r.Context(), &req)
	if err != nil {
		c.SendBadRequest(w, err.Error())
		return
	}

	c.SendOne(w, permission)
}

// Update handles React Admin update - PUT /permissions/123
// Returns updated object
// @Summary Update permission
// @Description Update an existing permission by ID.
// @Tags permissions
// @Accept json
// @Produce json
// @Param id path string true "Permission ID"
// @Param request body permissions.CreatePermissionRequest true "Permission update details"
// @Success 200 {object} domain.Permission "Updated permission details"
// @Failure 400 {object} response.WebResponse "Invalid request body or ID required"
// @Security BearerAuth
// @Router /api/permissions/{id} [put]
func (c *PermissionController) Update(w http.ResponseWriter, r *http.Request) {
	id := c.GetIDFromURL(r)
	if id == "" {
		c.SendBadRequest(w, "ID is required")
		return
	}

	var req permissionRequest.CreatePermissionRequest
	if err := c.DecodeJSONBody(r, &req); err != nil {
		c.SendBadRequest(w, "Invalid request body")
		return
	}

	permission, err := c.permissionService.Update(r.Context(), id, &req)
	if err != nil {
		c.SendBadRequest(w, err.Error())
		return
	}

	c.SendOne(w, permission)
}

// UpdateMany handles React Admin updateMany - PUT /permissions?filter={"id":[123,124,125]}
// Returns array of updated IDs
// @Summary Update multiple permissions
// @Description Update multiple permissions with the same data. Supports React Admin's updateMany.
// @Tags permissions
// @Accept json
// @Produce json
// @Param filter query string true "Filter parameters with permission IDs, e.g., {\"id\":[\"1\",\"2\"]}"
// @Param request body map[string]interface{} true "Update data to apply"
// @Success 200 {array} string "Array of updated permission IDs"
// @Failure 400 {object} response.WebResponse "Invalid query or missing id filter"
// @Security BearerAuth
// @Router /api/permissions/bulk [put]
func (c *PermissionController) UpdateMany(w http.ResponseWriter, r *http.Request) {
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

	updatedIDs, err := c.permissionService.UpdateMany(r.Context(), ids, updates)
	if err != nil {
		c.SendBadRequest(w, err.Error())
		return
	}

	c.SendIDs(w, updatedIDs)
}

// Delete handles React Admin delete - DELETE /permissions/123
// Returns deleted object
// @Summary Delete permission
// @Description Delete a permission by its ID. Returns the deleted object.
// @Tags permissions
// @Accept json
// @Produce json
// @Param id path string true "Permission ID"
// @Success 200 {object} domain.Permission "Deleted permission details"
// @Failure 400 {object} response.WebResponse "ID is required"
// @Failure 404 {object} response.WebResponse "Permission not found"
// @Security BearerAuth
// @Router /api/permissions/{id} [delete]
func (c *PermissionController) Delete(w http.ResponseWriter, r *http.Request) {
	id := c.GetIDFromURL(r)
	if id == "" {
		c.SendBadRequest(w, "ID is required")
		return
	}

	// Get permission before deletion (to return it)
	permission, err := c.permissionService.GetByID(r.Context(), id)
	if err != nil {
		c.SendNotFound(w, err.Error())
		return
	}

	if err := c.permissionService.Delete(r.Context(), id); err != nil {
		c.SendBadRequest(w, err.Error())
		return
	}

	c.SendOne(w, permission)
}

// DeleteMany handles React Admin deleteMany - DELETE /permissions?filter={"id":[123,124,125]}
// Returns array of deleted IDs
// @Summary Delete multiple permissions
// @Description Delete multiple permissions by their IDs. Supports React Admin's deleteMany.
// @Tags permissions
// @Accept json
// @Produce json
// @Param filter query string true "Filter parameters with permission IDs, e.g., {\"id\":[\"1\",\"2\"]}"
// @Success 200 {array} string "Array of deleted permission IDs"
// @Failure 400 {object} response.WebResponse "Invalid query or missing id filter"
// @Security BearerAuth
// @Router /api/permissions [delete]
func (c *PermissionController) DeleteMany(w http.ResponseWriter, r *http.Request) {
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

	if err := c.permissionService.DeleteAll(r.Context(), ids); err != nil {
		c.SendBadRequest(w, err.Error())
		return
	}

	c.SendIDs(w, ids)
}
