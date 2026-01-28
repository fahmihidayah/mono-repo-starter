package storage

import (
	"bytes"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/fahmihidayah/go-api-orchestrator/internal/utils"
)

// S3Storage implements Storage interface for AWS S3 storage
type S3Storage struct {
	client    *s3.S3
	bucket    string
	region    string
	endpoint  string
	publicURL string
}

// S3Config contains AWS S3 configuration
type S3Config struct {
	Region          string
	Bucket          string
	AccessKeyID     string
	SecretAccessKey string
	BaseURL         string // Optional: for custom S3 endpoints (like MinIO, DigitalOcean Spaces)
}

// NewS3Storage creates a new S3 storage instance
func NewS3Storage(cfg S3Config) (*S3Storage, error) {
	// Configure AWS session
	awsConfig := &aws.Config{
		Region: aws.String(cfg.Region),
		Credentials: credentials.NewStaticCredentials(
			cfg.AccessKeyID,
			cfg.SecretAccessKey,
			"",
		),
	}

	// Add custom endpoint if provided (for MinIO, DigitalOcean Spaces, etc.)
	if cfg.BaseURL != "" {
		awsConfig.Endpoint = aws.String(cfg.BaseURL)
		awsConfig.S3ForcePathStyle = aws.Bool(true) // Required for MinIO and some S3-compatible services
	}

	// Create AWS session
	sess, err := session.NewSession(awsConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create AWS session: %v", err)
	}

	// Create S3 client
	client := s3.New(sess)

	// Determine public URL
	publicURL := cfg.BaseURL
	if publicURL == "" {
		// Standard AWS S3 URL
		if cfg.Region == "us-east-1" {
			publicURL = fmt.Sprintf("https://%s.s3.amazonaws.com", cfg.Bucket)
		} else {
			publicURL = fmt.Sprintf("https://%s.s3.%s.amazonaws.com", cfg.Bucket, cfg.Region)
		}
	} else {
		// Custom endpoint URL (MinIO, DigitalOcean Spaces, etc.)
		publicURL = fmt.Sprintf("%s/%s", strings.TrimSuffix(cfg.BaseURL, "/"), cfg.Bucket)
	}

	return &S3Storage{
		client:    client,
		bucket:    cfg.Bucket,
		region:    cfg.Region,
		endpoint:  cfg.BaseURL,
		publicURL: publicURL,
	}, nil
}

// Upload uploads a file to S3
func (s *S3Storage) Upload(file multipart.File, fileHeader *multipart.FileHeader) (*FileInfo, error) {
	// Get file info
	fileName := fileHeader.Filename
	fileSize := fileHeader.Size
	mimeType := fileHeader.Header.Get("Content-Type")

	// Generate unique filename with timestamp
	ext := filepath.Ext(fileName)
	timestamp := time.Now().Unix()
	uniqueFileName := fmt.Sprintf("%d_%s%s", timestamp, strings.ReplaceAll(utils.GenerateSlug(fileName), ext, ""), ext)

	// S3 key (path in bucket)
	key := fmt.Sprintf("uploads/%s", uniqueFileName)

	// Read file content into buffer
	buf := new(bytes.Buffer)
	if _, err := io.Copy(buf, file); err != nil {
		return nil, fmt.Errorf("failed to read file: %v", err)
	}

	// Get image dimensions if it's an image
	width := 0
	height := 0
	if strings.HasPrefix(mimeType, "image/") {
		// Decode image from buffer
		if img, _, err := image.DecodeConfig(bytes.NewReader(buf.Bytes())); err == nil {
			width = img.Width
			height = img.Height
		}
	}

	// Prepare upload input
	uploadInput := &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(key),
		Body:        bytes.NewReader(buf.Bytes()),
		ContentType: aws.String(mimeType),
		ACL:         aws.String("public-read"), // Make file publicly accessible
	}

	// Upload to S3
	_, err := s.client.PutObject(uploadInput)
	if err != nil {
		return nil, fmt.Errorf("failed to upload to S3: %v", err)
	}

	// Generate public URL
	publicURL := fmt.Sprintf("%s/%s", s.publicURL, key)

	return &FileInfo{
		FileName:     uniqueFileName,
		OriginalName: fileName,
		MimeType:     mimeType,
		FileSize:     fileSize,
		Width:        width,
		Height:       height,
		URL:          publicURL,
		Path:         fmt.Sprintf("/%s", key), // Store relative path
	}, nil
}

// Delete deletes a file from S3
func (s *S3Storage) Delete(path string) error {
	// Remove leading slash if present
	path = strings.TrimPrefix(path, "/")

	// Prepare delete input
	deleteInput := &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(path),
	}

	// Delete from S3
	_, err := s.client.DeleteObject(deleteInput)
	if err != nil {
		return fmt.Errorf("failed to delete from S3: %v", err)
	}

	return nil
}

// GetURL returns the public URL for a file
func (s *S3Storage) GetURL(path string) string {
	// If path already contains the base URL, return as is
	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		return path
	}

	// Remove leading slash if present
	path = strings.TrimPrefix(path, "/")

	return fmt.Sprintf("%s/%s", s.publicURL, path)
}

// Exists checks if a file exists in S3
func (s *S3Storage) Exists(path string) (bool, error) {
	// Remove leading slash if present
	path = strings.TrimPrefix(path, "/")

	// Prepare head object input
	headInput := &s3.HeadObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(path),
	}

	// Head object to check existence
	_, err := s.client.HeadObject(headInput)
	if err != nil {
		// Check if error is "not found"
		if strings.Contains(err.Error(), "NotFound") || strings.Contains(err.Error(), "404") {
			return false, nil
		}
		return false, err
	}

	return true, nil
}
