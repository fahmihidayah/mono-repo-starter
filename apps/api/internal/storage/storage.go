package storage

import (
	"io"
	"mime/multipart"
)

// FileInfo contains information about an uploaded file
type FileInfo struct {
	FileName     string
	OriginalName string
	MimeType     string
	FileSize     int64
	Width        int
	Height       int
	URL          string
	Path         string
}

// Storage defines the interface for file storage operations
type Storage interface {
	// Upload uploads a file and returns file information
	Upload(file multipart.File, fileHeader *multipart.FileHeader) (*FileInfo, error)

	// Delete deletes a file by its path
	Delete(path string) error

	// GetURL returns the public URL for a file
	GetURL(path string) string

	// Exists checks if a file exists
	Exists(path string) (bool, error)
}

// UploadResult contains the result of a file upload operation
type UploadResult struct {
	Success bool
	FileInfo *FileInfo
	Error   error
}

// Helper function to copy file content
func copyFile(dst io.Writer, src io.Reader) (int64, error) {
	return io.Copy(dst, src)
}
