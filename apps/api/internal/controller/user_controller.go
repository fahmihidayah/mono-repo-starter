package controller

import (
	"encoding/json"
	"log"
	"net/http"

	userRequest "github.com/fahmihidayah/go-api-orchestrator/internal/data/request/users"
	"github.com/fahmihidayah/go-api-orchestrator/internal/service"
)

// UserController handles User endpoints with React Admin compatibility
type UserController struct {
	BaseController
	userService service.IUserService
}

// UserControllerProvider creates a new UserControllerReactAdmin
func UserControllerProvider(userService service.IUserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

// GetList handles React Admin getList - GET /users?sort=[...]&range=[...]&filter={...}
// Returns array with X-Total-Count header
// @Summary Get list of users
// @Description Retrieve a paginated list of users with optional filtering and sorting. Supports React Admin's getList and getMany operations.
// @Tags users
// @Accept json
// @Produce json
// @Param sort query string false "Sort parameters in format [field,order], e.g., [\"name\",\"ASC\"]"
// @Param range query string false "Pagination range in format [start,end], e.g., [0,9]"
// @Param filter query string false "Filter parameters as JSON object, e.g., {\"name\":\"John\"}. Use {\"ids\":[\"id1\",\"id2\"]} for getMany operation"
// @Success 200 {array} domain.User "List of users with X-Total-Count header"
// @Failure 400 {object} response.WebResponse "Invalid query parameters"
// @Failure 500 {object} response.WebResponse "Internal server error"
// @Security BearerAuth
// @Router /api/users [get]
func (c *UserController) GetList(w http.ResponseWriter, r *http.Request) {
	log.Printf("[UserController.GetList] Request URL: %s", r.URL.String())
	log.Printf("[UserController.GetList] Query String: %s", r.URL.RawQuery)

	params, err := c.ParseQueryParams(r)
	if err != nil {
		log.Printf("[UserController.GetList] Error parsing params: %v", err)
		c.SendBadRequest(w, err.Error())
		return
	}

	// Log parsed parameters
	paramsJSON, _ := json.Marshal(params)
	log.Printf("[UserController.GetList] Parsed params: %s", string(paramsJSON))

	// Check for getMany case (filter contains "ids")
	if ids, hasIDs := params.GetFilterIDs(); hasIDs {
		log.Printf("[UserController.GetList] GetMany detected - IDs: %v", ids)
		c.getMany(w, r, ids)
		return
	}

	log.Printf("[UserController.GetList] Pagination - Limit: %d, Offset: %d", params.Limit, params.Offset)
	log.Printf("[UserController.GetList] Sort - Field: %s, Order: %s", params.Sort[0], params.Sort[1])

	users, total, err := c.userService.GetWithQueryParams(r.Context(), params)
	if err != nil {
		log.Printf("[UserController.GetList] Service error: %v", err)
		c.SendInternalError(w, err.Error())
		return
	}

	log.Printf("[UserController.GetList] Results - Total: %d, Users returned: %d", total, len(users))
	c.SendListWithPagination(w, users, params, total)
}

// getMany handles React Admin getMany - GET /users?filter={"ids":[123,124,125]}
func (c *UserController) getMany(w http.ResponseWriter, r *http.Request, ids []string) {
	users, err := c.userService.GetByIDs(r.Context(), ids)
	if err != nil {
		c.SendInternalError(w, err.Error())
		return
	}

	// For getMany, don't send total count
	c.SendOne(w, users)
}

// GetOne handles React Admin getOne - GET /users/123
// Returns single object
// @Summary Get user by ID
// @Description Retrieve a single user by their unique identifier
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} domain.User "User details"
// @Failure 400 {object} response.WebResponse "ID is required"
// @Failure 404 {object} response.WebResponse "User not found"
// @Security BearerAuth
// @Router /api/users/{id} [get]
func (c *UserController) GetOne(w http.ResponseWriter, r *http.Request) {
	id := c.GetIDFromURL(r)
	if id == "" {
		c.SendBadRequest(w, "ID is required")
		return
	}

	user, err := c.userService.GetByID(r.Context(), id)
	if err != nil {
		c.SendNotFound(w, err.Error())
		return
	}

	c.SendOne(w, user)
}

// Create handles React Admin create - POST /users
// Returns created object with ID
// @Summary Create a new user (Admin)
// @Description Create a new user account with admin privileges. User is automatically verified (no email verification required). For public registration with email verification, use /api/users/auth/register instead.
// @Tags users
// @Accept json
// @Produce json
// @Param request body request.CreateUserRequest true "User creation details"
// @Success 200 {object} domain.User "Created user details"
// @Failure 400 {object} response.WebResponse "Invalid request body or validation error"
// @Failure 500 {object} response.WebResponse "Failed to fetch created user"
// @Security BearerAuth
// @Router /api/users [post]
func (c *UserController) Create(w http.ResponseWriter, r *http.Request) {
	var req userRequest.CreateUserRequest
	if err := c.DecodeJSONBody(r, &req); err != nil {
		c.SendBadRequest(w, "Invalid request body")
		return
	}

	log.Printf("[UserController.Create] Creating user - Name: %s, Email: %s", req.Name, req.Email)

	// Create user using admin creation (auto-verified)
	if err := c.userService.Create(r.Context(), &req); err != nil {
		c.SendBadRequest(w, err.Error())
		return
	}

	// Fetch the created user to return it
	user, err := c.userService.GetByEmail(r.Context(), req.Email)
	if err != nil {
		c.SendInternalError(w, "User created but failed to fetch: "+err.Error())
		return
	}

	c.SendOne(w, user)
}

// Update handles React Admin update - PUT /users/123
// Returns updated object
// @Summary Update user
// @Description Update an existing user's information by ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param request body request.UpdateUserRequest true "User update details"
// @Success 200 {object} domain.User "Updated user details"
// @Failure 400 {object} response.WebResponse "Invalid request body, ID required, or validation error"
// @Failure 500 {object} response.WebResponse "Failed to fetch updated user"
// @Security BearerAuth
// @Router /api/users/{id} [put]
func (c *UserController) Update(w http.ResponseWriter, r *http.Request) {
	id := c.GetIDFromURL(r)
	if id == "" {
		c.SendBadRequest(w, "ID is required")
		return
	}

	var req userRequest.UpdateUserRequest
	if err := c.DecodeJSONBody(r, &req); err != nil {
		c.SendBadRequest(w, "Invalid request body")
		return
	}

	// Set ID from URL
	req.ID = id

	user, err := c.userService.Update(r.Context(), &req)
	if err != nil {
		c.SendBadRequest(w, err.Error())
		return
	}

	c.SendOne(w, user)
}

// Delete handles React Admin delete - DELETE /users/123
// Returns deleted object
// @Summary Delete user
// @Description Delete a user by their ID. Returns the deleted user object.
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} domain.User "Deleted user details"
// @Failure 400 {object} response.WebResponse "ID is required or deletion failed"
// @Failure 404 {object} response.WebResponse "User not found"
// @Security BearerAuth
// @Router /api/users/{id} [delete]
func (c *UserController) Delete(w http.ResponseWriter, r *http.Request) {
	id := c.GetIDFromURL(r)
	if id == "" {
		c.SendBadRequest(w, "ID is required")
		return
	}

	// Get user before deletion (to return it)
	user, err := c.userService.GetByID(r.Context(), id)
	if err != nil {
		c.SendNotFound(w, err.Error())
		return
	}

	if err := c.userService.Delete(r.Context(), id); err != nil {
		c.SendBadRequest(w, err.Error())
		return
	}

	c.SendOne(w, user)
}

// DeleteMany handles React Admin deleteMany - DELETE /users?filter={"id":[123,124,125]}
// Returns array of deleted IDs
// @Summary Delete multiple users
// @Description Delete multiple users by their IDs. Supports React Admin's deleteMany operation.
// @Tags users
// @Accept json
// @Produce json
// @Param filter query string true "Filter parameters with user IDs, e.g., {\"id\":[\"id1\",\"id2\",\"id3\"]}"
// @Success 200 {array} string "Array of deleted user IDs"
// @Failure 400 {object} response.WebResponse "Invalid query parameters or missing id filter"
// @Security BearerAuth
// @Router /api/users [delete]
func (c *UserController) DeleteMany(w http.ResponseWriter, r *http.Request) {
	params, err := c.ParseQueryParams(r)
	if err != nil {
		c.SendBadRequest(w, err.Error())
		return
	}

	ids, hasIDs := params.GetFilterIDs()
	idJsons, _ := json.Marshal(ids)
	log.Printf("ids %s", idJsons)
	if !hasIDs {
		c.SendBadRequest(w, "Missing id filter")
		return
	}

	bulkDeleteReq := &userRequest.BulkDeleteRequest{
		IDs: ids,
	}

	if err := c.userService.DeleteAll(r.Context(), bulkDeleteReq); err != nil {
		c.SendBadRequest(w, err.Error())
		return
	}

	c.SendIDs(w, ids)
}
