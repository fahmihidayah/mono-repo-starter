package request

import "mime/multipart"

// UploadMediaRequest represents multipart form data for media upload
type UploadMediaRequest struct {
	Media *multipart.FileHeader `form:"media" validate:"required"`
	Alt   string                `form:"alt" validate:"max=255"`
}

type UpdateMediaRequest struct {
	Alt string `json:"alt" validate:"max=255"`
}
