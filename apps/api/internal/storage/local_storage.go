package storage

import (
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fahmihidayah/go-api-orchestrator/internal/utils"
)

// LocalStorage implements Storage interface for local file system storage
type LocalStorage struct {
	uploadDir string
	baseURL   string
}

// NewLocalStorage creates a new local storage instance
func NewLocalStorage(uploadDir, baseURL string) *LocalStorage {
	return &LocalStorage{
		uploadDir: uploadDir,
		baseURL:   baseURL,
	}
}

// Upload uploads a file to local file system
func (s *LocalStorage) Upload(file multipart.File, fileHeader *multipart.FileHeader) (*FileInfo, error) {
	// Get file info
	fileName := fileHeader.Filename
	fileSize := fileHeader.Size
	mimeType := fileHeader.Header.Get("Content-Type")

	// Generate unique filename with timestamp
	ext := filepath.Ext(fileName)
	timestamp := time.Now().Unix()
	uniqueFileName := fmt.Sprintf("%d_%s%s", timestamp, strings.ReplaceAll(utils.GenerateSlug(fileName), ext, ""), ext)

	// Create upload directory if it doesn't exist
	if err := os.MkdirAll(s.uploadDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create upload directory: %v", err)
	}

	// Create destination file
	destPath := filepath.Join(s.uploadDir, uniqueFileName)
	dest, err := os.Create(destPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create destination file: %v", err)
	}
	defer dest.Close()

	// Copy uploaded file to destination
	if _, err := io.Copy(dest, file); err != nil {
		return nil, fmt.Errorf("failed to save file: %v", err)
	}

	// Get image dimensions if it's an image
	width := 0
	height := 0
	if strings.HasPrefix(mimeType, "image/") {
		// Rewind file for reading
		file.Seek(0, 0)
		if img, _, err := image.DecodeConfig(file); err == nil {
			width = img.Width
			height = img.Height
		}
	}

	// Generate URL path (relative path for storage)
	relativePath := fmt.Sprintf("/uploads/%s", uniqueFileName)
	fullURL := fmt.Sprintf("%s%s", s.baseURL, relativePath)

	return &FileInfo{
		FileName:     uniqueFileName,
		OriginalName: fileName,
		MimeType:     mimeType,
		FileSize:     fileSize,
		Width:        width,
		Height:       height,
		URL:          fullURL,
		Path:         relativePath,
	}, nil
}

// Delete deletes a file from local file system
func (s *LocalStorage) Delete(path string) error {
	// Remove leading slash if present
	path = strings.TrimPrefix(path, "/")

	// If path starts with "uploads/", construct full path
	var fullPath string
	if strings.HasPrefix(path, "uploads/") {
		fullPath = filepath.Join(s.uploadDir, strings.TrimPrefix(path, "uploads/"))
	} else {
		fullPath = filepath.Join(s.uploadDir, path)
	}

	// Check if file exists
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return fmt.Errorf("file not found: %s", path)
	}

	// Delete the file
	if err := os.Remove(fullPath); err != nil {
		return fmt.Errorf("failed to delete file: %v", err)
	}

	return nil
}

// GetURL returns the public URL for a file
func (s *LocalStorage) GetURL(path string) string {
	// If path already contains the base URL, return as is
	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		return path
	}

	// Ensure path starts with /
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	return fmt.Sprintf("%s%s", s.baseURL, path)
}

// Exists checks if a file exists in local file system
func (s *LocalStorage) Exists(path string) (bool, error) {
	// Remove leading slash if present
	path = strings.TrimPrefix(path, "/")

	// If path starts with "uploads/", construct full path
	var fullPath string
	if strings.HasPrefix(path, "uploads/") {
		fullPath = filepath.Join(s.uploadDir, strings.TrimPrefix(path, "uploads/"))
	} else {
		fullPath = filepath.Join(s.uploadDir, path)
	}

	_, err := os.Stat(fullPath)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
