package controller

import (
	"encoding/json"
	"log"
	"net/http"

	postRequest "github.com/fahmihidayah/go-api-orchestrator/internal/data/request/posts"
	"github.com/fahmihidayah/go-api-orchestrator/internal/middleware"
	"github.com/fahmihidayah/go-api-orchestrator/internal/service"
)

// PostController handles Post endpoints with React Admin compatibility
type PostController struct {
	BaseController
	postService service.IPostService
}

// PostControllerProvider creates a new PostControllerReactAdmin
func PostControllerProvider(postService service.IPostService) *PostController {
	return &PostController{
		postService: postService,
	}
}

// GetList handles React Admin getList - GET /posts?sort=[...]&range=[...]&filter={...}
// Returns array with X-Total-Count header
// @Summary Get list of posts
// @Description Retrieve a paginated list of posts with optional filtering and sorting. Supports React Admin's getList and getMany operations.
// @Tags posts
// @Accept json
// @Produce json
// @Param sort query string false "Sort parameters in format [field,order], e.g., [\"title\",\"ASC\"]"
// @Param range query string false "Pagination range in format [start,end], e.g., [0,9]"
// @Param filter query string false "Filter parameters as JSON object, e.g., {\"title\":\"My Post\"}. Use {\"ids\":[\"id1\",\"id2\"]} for getMany operation"
// @Success 200 {array} domain.Post "List of posts with X-Total-Count header"
// @Failure 400 {object} response.WebResponse "Invalid query parameters"
// @Failure 500 {object} response.WebResponse "Internal server error"
// @Security BearerAuth
// @Router /api/posts [get]
func (c *PostController) GetList(w http.ResponseWriter, r *http.Request) {
	log.Printf("[PostController.GetList] Request URL: %s", r.URL.String())
	log.Printf("[PostController.GetList] Query String: %s", r.URL.RawQuery)

	params, err := c.ParseQueryParams(r)
	if err != nil {
		log.Printf("[PostController.GetList] Error parsing params: %v", err)
		c.SendBadRequest(w, err.Error())
		return
	}

	// Log parsed parameters
	paramsJSON, _ := json.Marshal(params)
	log.Printf("[PostController.GetList] Parsed params: %s", string(paramsJSON))

	// Check for getMany case (filter contains "ids")
	if ids, hasIDs := params.GetFilterIDs(); hasIDs {
		log.Printf("[PostController.GetList] GetMany detected - IDs: %v", ids)
		c.getMany(w, r, ids)
		return
	}

	log.Printf("[PostController.GetList] Pagination - Limit: %d, Offset: %d", params.Limit, params.Offset)
	log.Printf("[PostController.GetList] Sort - Field: %s, Order: %s", params.Sort[0], params.Sort[1])

	posts, total, err := c.postService.GetWithQueryParams(r.Context(), params)
	if err != nil {
		log.Printf("[PostController.GetList] Service error: %v", err)
		c.SendInternalError(w, err.Error())
		return
	}

	log.Printf("[PostController.GetList] Results - Total: %d, Posts returned: %d", total, len(posts))
	c.SendList(w, posts, total)
}

// getMany handles React Admin getMany - GET /posts?filter={"ids":[123,124,125]}
func (c *PostController) getMany(w http.ResponseWriter, r *http.Request, ids []string) {
	posts, err := c.postService.GetByIDs(r.Context(), ids)
	if err != nil {
		c.SendInternalError(w, err.Error())
		return
	}

	// For getMany, don't send total count
	c.SendOne(w, posts)
}

// GetOne handles React Admin getOne - GET /posts/123
// Returns single object
// @Summary Get post by ID
// @Description Retrieve a single post by its unique identifier
// @Tags posts
// @Accept json
// @Produce json
// @Param id path string true "Post ID"
// @Success 200 {object} domain.Post "Post details"
// @Failure 400 {object} response.WebResponse "ID is required"
// @Failure 404 {object} response.WebResponse "Post not found"
// @Security BearerAuth
// @Router /api/posts/{id} [get]
func (c *PostController) GetOne(w http.ResponseWriter, r *http.Request) {
	id := c.GetIDFromURL(r)
	if id == "" {
		c.SendBadRequest(w, "ID is required")
		return
	}

	post, err := c.postService.GetByID(r.Context(), id)
	if err != nil {
		c.SendNotFound(w, err.Error())
		return
	}

	c.SendOne(w, post)
}

// Create handles React Admin create - POST /posts
// Returns created object with ID
// @Summary Create a new post
// @Description Create a new post with title, content, and categories. Requires authentication.
// @Tags posts
// @Accept json
// @Produce json
// @Param request body request.CreatePostRequest true "Post creation details"
// @Success 200 {object} domain.Post "Created post details"
// @Failure 400 {object} response.WebResponse "Invalid request body or validation error"
// @Failure 401 {object} response.WebResponse "Unauthorized - authentication required"
// @Security BearerAuth
// @Router /api/posts [post]
func (c *PostController) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		c.SendUnauthorized(w, "unauthorized")
		return
	}

	var req postRequest.CreatePostRequest
	if err := c.DecodeJSONBody(r, &req); err != nil {
		c.SendBadRequest(w, "Invalid request body")
		return
	}

	log.Printf("[PostController.Create] Creating post - Title: %s, UserID: %s", req.Title, userID)

	post, err := c.postService.Create(r.Context(), &req, userID)
	if err != nil {
		c.SendBadRequest(w, err.Error())
		return
	}

	c.SendOne(w, post)
}

// Update handles React Admin update - PUT /posts/123
// Returns updated object
// @Summary Update post
// @Description Update an existing post's information by ID. Requires authentication.
// @Tags posts
// @Accept json
// @Produce json
// @Param id path string true "Post ID"
// @Param request body request.CreatePostRequest true "Post update details"
// @Success 200 {object} domain.Post "Updated post details"
// @Failure 400 {object} response.WebResponse "Invalid request body, ID required, or validation error"
// @Failure 401 {object} response.WebResponse "Unauthorized - authentication required"
// @Security BearerAuth
// @Router /api/posts/{id} [put]
func (c *PostController) Update(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		c.SendUnauthorized(w, "unauthorized")
		return
	}

	id := c.GetIDFromURL(r)
	if id == "" {
		c.SendBadRequest(w, "ID is required")
		return
	}

	var req postRequest.CreatePostRequest
	if err := c.DecodeJSONBody(r, &req); err != nil {
		c.SendBadRequest(w, "Invalid request body")
		return
	}

	post, err := c.postService.Update(r.Context(), id, &req, userID)
	if err != nil {
		c.SendBadRequest(w, err.Error())
		return
	}

	c.SendOne(w, post)
}

// UpdateMany handles React Admin updateMany - PUT /posts?filter={"id":[123,124,125]}
// Returns array of updated IDs
// @Summary Update multiple posts
// @Description Update multiple posts with the same data. Supports React Admin's updateMany operation. Requires authentication.
// @Tags posts
// @Accept json
// @Produce json
// @Param filter query string true "Filter parameters with post IDs, e.g., {\"id\":[\"id1\",\"id2\",\"id3\"]}"
// @Param request body map[string]interface{} true "Update data to apply to all posts"
// @Success 200 {array} string "Array of updated post IDs"
// @Failure 400 {object} response.WebResponse "Invalid query parameters or missing id filter"
// @Failure 401 {object} response.WebResponse "Unauthorized - authentication required"
// @Security BearerAuth
// @Router /api/posts/bulk [put]
func (c *PostController) UpdateMany(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		c.SendUnauthorized(w, "unauthorized")
		return
	}

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

	updatedIDs, err := c.postService.UpdateMany(r.Context(), ids, updates, userID)
	if err != nil {
		c.SendBadRequest(w, err.Error())
		return
	}

	c.SendIDs(w, updatedIDs)
}

// Delete handles React Admin delete - DELETE /posts/123
// Returns deleted object
// @Summary Delete post
// @Description Delete a post by its ID. Returns the deleted post object. Requires authentication.
// @Tags posts
// @Accept json
// @Produce json
// @Param id path string true "Post ID"
// @Success 200 {object} domain.Post "Deleted post details"
// @Failure 400 {object} response.WebResponse "ID is required or deletion failed"
// @Failure 401 {object} response.WebResponse "Unauthorized - authentication required"
// @Failure 404 {object} response.WebResponse "Post not found"
// @Security BearerAuth
// @Router /api/posts/{id} [delete]
func (c *PostController) Delete(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		c.SendUnauthorized(w, "unauthorized")
		return
	}

	id := c.GetIDFromURL(r)
	if id == "" {
		c.SendBadRequest(w, "ID is required")
		return
	}

	// Get post before deletion (to return it)
	post, err := c.postService.GetByID(r.Context(), id)
	if err != nil {
		c.SendNotFound(w, err.Error())
		return
	}

	if err := c.postService.Delete(r.Context(), id, userID); err != nil {
		c.SendBadRequest(w, err.Error())
		return
	}

	c.SendOne(w, post)
}

// DeleteMany handles React Admin deleteMany - DELETE /posts?filter={"id":[123,124,125]}
// Returns array of deleted IDs
// @Summary Delete multiple posts
// @Description Delete multiple posts by their IDs. Supports React Admin's deleteMany operation. Requires authentication.
// @Tags posts
// @Accept json
// @Produce json
// @Param filter query string true "Filter parameters with post IDs, e.g., {\"id\":[\"id1\",\"id2\",\"id3\"]}"
// @Success 200 {array} string "Array of deleted post IDs"
// @Failure 400 {object} response.WebResponse "Invalid query parameters or missing id filter"
// @Failure 401 {object} response.WebResponse "Unauthorized - authentication required"
// @Security BearerAuth
// @Router /api/posts [delete]
func (c *PostController) DeleteMany(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		c.SendUnauthorized(w, "unauthorized")
		return
	}

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

	if err := c.postService.DeleteAll(r.Context(), ids, userID); err != nil {
		c.SendBadRequest(w, err.Error())
		return
	}

	c.SendIDs(w, ids)
}
