package storage

import (
	"fmt"
	"log"

	"github.com/fahmihidayah/go-api-orchestrator/internal/config"
)

// Factory creates storage instances based on configuration
type Factory struct {
	config *config.Config
}

// NewFactory creates a new storage factory
func NewFactory(cfg *config.Config) *Factory {
	return &Factory{
		config: cfg,
	}
}

// CreateStorage creates a storage instance based on the configured storage type
func (f *Factory) CreateStorage() (Storage, error) {
	storageType := f.config.Storage.Type

	switch storageType {
	case "local":
		log.Printf("Initializing local storage: uploadDir=%s", f.config.Storage.LocalUploadDir)
		return NewLocalStorage(f.config.Storage.LocalUploadDir, f.config.BaseURL), nil

	case "s3":
		log.Printf("Initializing S3 storage: bucket=%s, region=%s", f.config.Storage.S3Bucket, f.config.Storage.S3Region)
		s3Config := S3Config{
			Region:          f.config.Storage.S3Region,
			Bucket:          f.config.Storage.S3Bucket,
			AccessKeyID:     f.config.Storage.S3AccessKeyID,
			SecretAccessKey: f.config.Storage.S3SecretAccessKey,
			BaseURL:         f.config.Storage.S3Endpoint, // For custom endpoints (MinIO, etc.)
		}

		storage, err := NewS3Storage(s3Config)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize S3 storage: %v", err)
		}
		return storage, nil

	default:
		return nil, fmt.Errorf("unsupported storage type: %s (must be 'local' or 's3')", storageType)
	}
}

// StorageProvider provides a storage instance for dependency injection
func StorageProvider(cfg *config.Config) (Storage, error) {
	factory := NewFactory(cfg)
	return factory.CreateStorage()
}
