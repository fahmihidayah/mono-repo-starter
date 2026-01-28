package controller

import (
	"encoding/json"
	"log"
	"mime/multipart"
	"net/http"

	mediaRequest "github.com/fahmihidayah/go-api-orchestrator/internal/data/request/media"
	"github.com/fahmihidayah/go-api-orchestrator/internal/service"
)

// MediaController handles Media endpoints with React Admin compatibility
type MediaController struct {
	BaseController
	mediaService service.IMediaService
}

// MediaControllerProvider creates a new MediaControllerReactAdmin
func MediaControllerProvider(mediaService service.IMediaService) *MediaController {
	return &MediaController{
		mediaService: mediaService,
	}
}

// GetList handles React Admin getList - GET /media?sort=[...]&range=[...]&filter={...}
// Returns array with X-Total-Count header
// @Summary Get list of media files
// @Description Retrieve a paginated list of media files with optional filtering and sorting. Supports React Admin's getList and getMany operations.
// @Tags media
// @Accept json
// @Produce json
// @Param sort query string false "Sort parameters in format [field,order], e.g., [\"file_name\",\"ASC\"]"
// @Param range query string false "Pagination range in format [start,end], e.g., [0,9]"
// @Param filter query string false "Filter parameters as JSON object, e.g., {\"mime_type\":\"image/jpeg\"}. Use {\"ids\":[\"id1\",\"id2\"]} for getMany operation"
// @Success 200 {array} domain.Media "List of media files with X-Total-Count header"
// @Failure 400 {object} response.WebResponse "Invalid query parameters"
// @Failure 500 {object} response.WebResponse "Internal server error"
// @Security BearerAuth
// @Router /api/media [get]
func (c *MediaController) GetList(w http.ResponseWriter, r *http.Request) {
	log.Printf("[MediaController.GetList] Request URL: %s", r.URL.String())
	log.Printf("[MediaController.GetList] Query String: %s", r.URL.RawQuery)

	params, err := c.ParseQueryParams(r)
	if err != nil {
		log.Printf("[MediaController.GetList] Error parsing params: %v", err)
		c.SendBadRequest(w, err.Error())
		return
	}

	// Log parsed parameters
	paramsJSON, _ := json.Marshal(params)
	log.Printf("[MediaController.GetList] Parsed params: %s", string(paramsJSON))

	// Check for getMany case (filter contains "ids")
	if ids, hasIDs := params.GetFilterIDs(); hasIDs {
		log.Printf("[MediaController.GetList] GetMany detected - IDs: %v", ids)
		c.getMany(w, r, ids)
		return
	}

	log.Printf("[MediaController.GetList] Pagination - Limit: %d, Offset: %d", params.Limit, params.Offset)
	log.Printf("[MediaController.GetList] Sort - Field: %s, Order: %s", params.Sort[0], params.Sort[1])

	media, total, err := c.mediaService.GetWithQueryParams(r.Context(), params)
	if err != nil {
		log.Printf("[MediaController.GetList] Service error: %v", err)
		c.SendInternalError(w, err.Error())
		return
	}

	log.Printf("[MediaController.GetList] Results - Total: %d, Media returned: %d", total, len(media))
	c.SendList(w, media, total)
}

// getMany handles React Admin getMany - GET /media?filter={"ids":[123,124,125]}
func (c *MediaController) getMany(w http.ResponseWriter, r *http.Request, ids []string) {
	media, err := c.mediaService.GetByIDs(r.Context(), ids)
	if err != nil {
		c.SendInternalError(w, err.Error())
		return
	}

	// For getMany, don't send total count
	c.SendOne(w, media)
}

// GetOne handles React Admin getOne - GET /media/123
// Returns single object
// @Summary Get media by ID
// @Description Retrieve a single media file by its unique identifier
// @Tags media
// @Accept json
// @Produce json
// @Param id path string true "Media ID"
// @Success 200 {object} domain.Media "Media file details"
// @Failure 400 {object} response.WebResponse "ID is required"
// @Failure 404 {object} response.WebResponse "Media not found"
// @Security BearerAuth
// @Router /api/media/{id} [get]
func (c *MediaController) GetOne(w http.ResponseWriter, r *http.Request) {
	id := c.GetIDFromURL(r)
	if id == "" {
		c.SendBadRequest(w, "ID is required")
		return
	}

	media, err := c.mediaService.GetByID(r.Context(), id)
	if err != nil {
		c.SendNotFound(w, err.Error())
		return
	}

	c.SendOne(w, media)
}

// Create handles React Admin create - POST /media
// Returns created object with ID
// Note: This handles multipart/form-data for file uploads
// @Summary Upload a media file
// @Description Upload a new media file (image, video, etc.) with optional alt text. Supports multipart/form-data with max 32MB file size.
// @Tags media
// @Accept multipart/form-data
// @Produce json
// @Param media formData file true "Media file to upload"
// @Param alt formData string false "Alternative text for the media"
// @Success 200 {object} domain.Media "Uploaded media details including URL and metadata"
// @Failure 400 {object} response.WebResponse "Invalid request or file upload failed"
// @Security BearerAuth
// @Router /api/media [post]
func (c *MediaController) Create(w http.ResponseWriter, r *http.Request) {
	// Parse multipart form data (max 32MB)
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		log.Printf("[MediaController.Create] Failed to parse multipart form: %v", err)
		c.SendBadRequest(w, "Failed to parse multipart form")
		return
	}

	// Get the file from form
	var file multipart.File
	var fileHeader *multipart.FileHeader
	var err error

	file, fileHeader, err = r.FormFile("media")
	if err != nil {
		log.Printf("[MediaController.Create] Failed to get file from form: %v", err)
		c.SendBadRequest(w, "Media file is required")
		return
	}
	defer file.Close()

	// Get alt text from form
	alt := r.FormValue("alt")

	// Create request
	req := &mediaRequest.UploadMediaRequest{
		Media: fileHeader,
		Alt:   alt,
	}

	log.Printf("[MediaController.Create] Uploading media - FileName: %s, Size: %d bytes", fileHeader.Filename, fileHeader.Size)

	media, err := c.mediaService.Upload(r.Context(), req, file)
	if err != nil {
		log.Printf("[MediaController.Create] Failed to upload media: %v", err)
		c.SendBadRequest(w, err.Error())
		return
	}

	log.Printf("[MediaController.Create] Media uploaded successfully - ID: %s", media.ID)

	c.SendOne(w, media)
}

// Update handles React Admin update - PUT /media/123
// Returns updated object
// @Summary Update media file metadata
// @Description Update media file metadata (e.g., alt text) by ID. Note: This does not replace the actual file.
// @Tags media
// @Accept json
// @Produce json
// @Param id path string true "Media ID"
// @Param request body media.UpdateMediaRequest true "Media metadata update details"
// @Success 200 {object} domain.Media "Updated media details"
// @Failure 400 {object} response.WebResponse "Invalid request body, ID required, or validation error"
// @Security BearerAuth
// @Router /api/media/{id} [put]
func (c *MediaController) Update(w http.ResponseWriter, r *http.Request) {
	id := c.GetIDFromURL(r)
	if id == "" {
		c.SendBadRequest(w, "ID is required")
		return
	}

	var req mediaRequest.UpdateMediaRequest
	if err := c.DecodeJSONBody(r, &req); err != nil {
		c.SendBadRequest(w, "Invalid request body")
		return
	}

	media, err := c.mediaService.Update(r.Context(), id, &req)
	if err != nil {
		c.SendBadRequest(w, err.Error())
		return
	}

	c.SendOne(w, media)
}

// UpdateMany handles React Admin updateMany - PUT /media?filter={"id":[123,124,125]}
// Returns array of updated IDs
// @Summary Update multiple media files
// @Description Update multiple media files' metadata with the same data. Supports React Admin's updateMany operation.
// @Tags media
// @Accept json
// @Produce json
// @Param filter query string true "Filter parameters with media IDs, e.g., {\"id\":[\"id1\",\"id2\",\"id3\"]}"
// @Param request body map[string]interface{} true "Update data to apply to all media files"
// @Success 200 {array} string "Array of updated media IDs"
// @Failure 400 {object} response.WebResponse "Invalid query parameters or missing id filter"
// @Security BearerAuth
// @Router /api/media/bulk [put]
func (c *MediaController) UpdateMany(w http.ResponseWriter, r *http.Request) {
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

	updatedIDs, err := c.mediaService.UpdateMany(r.Context(), ids, updates)
	if err != nil {
		c.SendBadRequest(w, err.Error())
		return
	}

	c.SendIDs(w, updatedIDs)
}

// Delete handles React Admin delete - DELETE /media/123
// Returns deleted object
// @Summary Delete media file
// @Description Delete a media file by its ID. This removes both the database record and the physical file from storage.
// @Tags media
// @Accept json
// @Produce json
// @Param id path string true "Media ID"
// @Success 200 {object} domain.Media "Deleted media details"
// @Failure 400 {object} response.WebResponse "ID is required or deletion failed"
// @Failure 404 {object} response.WebResponse "Media not found"
// @Security BearerAuth
// @Router /api/media/{id} [delete]
func (c *MediaController) Delete(w http.ResponseWriter, r *http.Request) {
	id := c.GetIDFromURL(r)
	if id == "" {
		c.SendBadRequest(w, "ID is required")
		return
	}

	// Get media before deletion (to return it)
	media, err := c.mediaService.GetByID(r.Context(), id)
	if err != nil {
		c.SendNotFound(w, err.Error())
		return
	}

	if err := c.mediaService.Delete(r.Context(), id); err != nil {
		c.SendBadRequest(w, err.Error())
		return
	}

	c.SendOne(w, media)
}

// DeleteMany handles React Admin deleteMany - DELETE /media?filter={"id":[123,124,125]}
// Returns array of deleted IDs
// @Summary Delete multiple media files
// @Description Delete multiple media files by their IDs. Removes both database records and physical files from storage. Supports React Admin's deleteMany operation.
// @Tags media
// @Accept json
// @Produce json
// @Param filter query string true "Filter parameters with media IDs, e.g., {\"id\":[\"id1\",\"id2\",\"id3\"]}"
// @Success 200 {array} string "Array of deleted media IDs"
// @Failure 400 {object} response.WebResponse "Invalid query parameters or missing id filter"
// @Security BearerAuth
// @Router /api/media [delete]
func (c *MediaController) DeleteMany(w http.ResponseWriter, r *http.Request) {
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

	if err := c.mediaService.DeleteAll(r.Context(), ids); err != nil {
		c.SendBadRequest(w, err.Error())
		return
	}

	c.SendIDs(w, ids)
}
