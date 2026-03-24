package controller

import (
	"encoding/json"
	"log"
	"net/http"

	roleRequest "github.com/fahmihidayah/go-api-orchestrator/internal/data/request"
	"github.com/fahmihidayah/go-api-orchestrator/internal/service"
)

// RoleController handles Role endpoints with React Admin compatibility
type RoleController struct {
	BaseController
	roleService service.IRoleService
}

// RoleControllerProvider creates a new RoleController
func RoleControllerProvider(roleService service.IRoleService) *RoleController {
	return &RoleController{
		roleService: roleService,
	}
}

// GetList handles React Admin getList - GET /roles?sort=[...]&range=[...]&filter={...}
// Returns array with X-Total-Count header
// @Summary Get list of roles
// @Description Retrieve a paginated list of roles with optional filtering and sorting. Supports React Admin's getList and getMany operations.
// @Tags roles
// @Accept json
// @Produce json
// @Param sort query string false "Sort parameters in format [field,order], e.g., [\"name\",\"ASC\"]"
// @Param range query string false "Pagination range in format [0,9]"
// @Param filter query string false "Filter parameters as JSON object, e.g., {\"name\":\"Admin\"}. Use {\"ids\":[\"1\",\"2\"]} for getMany"
// @Success 200 {array} domain.Role "List of roles with X-Total-Count header"
// @Failure 400 {object} response.WebResponse "Invalid query parameters"
// @Failure 500 {object} response.WebResponse "Internal server error"
// @Router /api/roles [get]
func (c *RoleController) GetList(w http.ResponseWriter, r *http.Request) {
	log.Printf("[RoleController.GetList] Request URL: %s", r.URL.String())
	log.Printf("[RoleController.GetList] Query String: %s", r.URL.RawQuery)

	params, err := c.ParseQueryParams(r)
	if err != nil {
		log.Printf("[RoleController.GetList] Error parsing params: %v", err)
		c.SendBadRequest(w, err.Error())
		return
	}

	// Log parsed parameters
	paramsJSON, _ := json.Marshal(params)
	log.Printf("[RoleController.GetList] Parsed params: %s", string(paramsJSON))

	// Check for getMany case (filter contains "ids")
	if ids, hasIDs := params.GetFilterIDs(); hasIDs {
		log.Printf("[RoleController.GetList] GetMany detected - IDs: %v", ids)
		c.getMany(w, r, ids)
		return
	}

	log.Printf("[RoleController.GetList] Pagination - Limit: %d, Offset: %d", params.Limit, params.Offset)
	log.Printf("[RoleController.GetList] Sort - Field: %s, Order: %s", params.Sort[0], params.Sort[1])

	roles, paginateInfo, err := c.roleService.GetWithQueryParams(r.Context(), params)
	if err != nil {
		log.Printf("[RoleController.GetList] Service error: %v", err)
		c.SendInternalError(w, err.Error())
		return
	}

	log.Printf("[RoleController.GetList] Results - Total: %d, Roles returned: %d", paginateInfo.TotalDocs, len(roles))
	c.SendListWithPagination(w, roles, paginateInfo)
}

// getMany handles React Admin getMany - GET /roles?filter={"ids":[123,124,125]}
func (c *RoleController) getMany(w http.ResponseWriter, r *http.Request, ids []string) {
	roles, err := c.roleService.GetByIDs(r.Context(), ids)
	if err != nil {
		c.SendInternalError(w, err.Error())
		return
	}

	// For getMany, don't send total count
	c.SendOne(w, roles)
}

// GetOne handles React Admin getOne - GET /roles/123
// Returns single object
// @Summary Get role by ID
// @Description Retrieve a single role by its unique identifier
// @Tags roles
// @Accept json
// @Produce json
// @Param id path string true "Role ID"
// @Success 200 {object} domain.Role "Role details"
// @Failure 400 {object} response.WebResponse "ID is required"
// @Failure 404 {object} response.WebResponse "Role not found"
// @Router /api/roles/{id} [get]
func (c *RoleController) GetOne(w http.ResponseWriter, r *http.Request) {
	id := c.GetIDFromURL(r)
	if id == "" {
		c.SendBadRequest(w, "ID is required")
		return
	}

	role, err := c.roleService.GetByID(r.Context(), id)
	if err != nil {
		c.SendNotFound(w, err.Error())
		return
	}

	c.SendOne(w, role)
}

// Create handles React Admin create - POST /roles
// Returns created object with ID
// @Summary Create a new role
// @Description Create a new role.
// @Tags roles
// @Accept json
// @Produce json
// @Param request body request.CreateRoleRequest true "Role creation details"
// @Success 200 {object} domain.Role "Created role details"
// @Failure 400 {object} response.WebResponse "Invalid request body or validation error"
// @Security BearerAuth
// @Router /api/roles [post]
func (c *RoleController) Create(w http.ResponseWriter, r *http.Request) {
	var req roleRequest.CreateRoleRequest
	if err := c.DecodeJSONBody(r, &req); err != nil {
		c.SendBadRequest(w, "Invalid request body")
		return
	}

	log.Printf("[RoleController.Create] Creating role - Name: %s", req.Name)

	role, err := c.roleService.Create(r.Context(), &req)
	if err != nil {
		c.SendBadRequest(w, err.Error())
		return
	}

	c.SendOne(w, role)
}

// Update handles React Admin update - PUT /roles/123
// Returns updated object
// @Summary Update role
// @Description Update an existing role by ID.
// @Tags roles
// @Accept json
// @Produce json
// @Param id path string true "Role ID"
// @Param request body request.CreateRoleRequest true "Role update details"
// @Success 200 {object} domain.Role "Updated role details"
// @Failure 400 {object} response.WebResponse "Invalid request body or ID required"
// @Security BearerAuth
// @Router /api/roles/{id} [put]
func (c *RoleController) Update(w http.ResponseWriter, r *http.Request) {
	id := c.GetIDFromURL(r)
	if id == "" {
		c.SendBadRequest(w, "ID is required")
		return
	}

	var req roleRequest.CreateRoleRequest
	if err := c.DecodeJSONBody(r, &req); err != nil {
		c.SendBadRequest(w, "Invalid request body")
		return
	}

	role, err := c.roleService.Update(r.Context(), id, &req)
	if err != nil {
		c.SendBadRequest(w, err.Error())
		return
	}

	c.SendOne(w, role)
}

// UpdateMany handles React Admin updateMany - PUT /roles?filter={"id":[123,124,125]}
// Returns array of updated IDs
// @Summary Update multiple roles
// @Description Update multiple roles with the same data. Supports React Admin's updateMany.
// @Tags roles
// @Accept json
// @Produce json
// @Param filter query string true "Filter parameters with role IDs, e.g., {\"id\":[\"1\",\"2\"]}"
// @Param request body map[string]interface{} true "Update data to apply"
// @Success 200 {array} string "Array of updated role IDs"
// @Failure 400 {object} response.WebResponse "Invalid query or missing id filter"
// @Security BearerAuth
// @Router /api/roles/bulk [put]
func (c *RoleController) UpdateMany(w http.ResponseWriter, r *http.Request) {
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

	updatedIDs, err := c.roleService.UpdateMany(r.Context(), ids, updates)
	if err != nil {
		c.SendBadRequest(w, err.Error())
		return
	}

	c.SendIDs(w, updatedIDs)
}

// Delete handles React Admin delete - DELETE /roles/123
// Returns deleted object
// @Summary Delete role
// @Description Delete a role by its ID. Returns the deleted object.
// @Tags roles
// @Accept json
// @Produce json
// @Param id path string true "Role ID"
// @Success 200 {object} domain.Role "Deleted role details"
// @Failure 400 {object} response.WebResponse "ID is required"
// @Failure 404 {object} response.WebResponse "Role not found"
// @Security BearerAuth
// @Router /api/roles/{id} [delete]
func (c *RoleController) Delete(w http.ResponseWriter, r *http.Request) {
	id := c.GetIDFromURL(r)
	if id == "" {
		c.SendBadRequest(w, "ID is required")
		return
	}

	// Get role before deletion (to return it)
	role, err := c.roleService.GetByID(r.Context(), id)
	if err != nil {
		c.SendNotFound(w, err.Error())
		return
	}

	if err := c.roleService.Delete(r.Context(), id); err != nil {
		c.SendBadRequest(w, err.Error())
		return
	}

	c.SendOne(w, role)
}

// DeleteMany handles React Admin deleteMany - DELETE /roles?filter={"id":[123,124,125]}
// Returns array of deleted IDs
// @Summary Delete multiple roles
// @Description Delete multiple roles by their IDs. Supports React Admin's deleteMany.
// @Tags roles
// @Accept json
// @Produce json
// @Param filter query string true "Filter parameters with role IDs, e.g., {\"id\":[\"1\",\"2\"]}"
// @Success 200 {array} string "Array of deleted role IDs"
// @Failure 400 {object} response.WebResponse "Invalid query or missing id filter"
// @Security BearerAuth
// @Router /api/roles [delete]
func (c *RoleController) DeleteMany(w http.ResponseWriter, r *http.Request) {
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

	if err := c.roleService.DeleteAll(r.Context(), ids); err != nil {
		c.SendBadRequest(w, err.Error())
		return
	}

	c.SendIDs(w, ids)
}
