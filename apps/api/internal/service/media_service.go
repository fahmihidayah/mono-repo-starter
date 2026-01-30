package service

import (
	"context"
	"errors"
	"mime/multipart"

	"github.com/fahmihidayah/go-api-orchestrator/internal/config"
	request "github.com/fahmihidayah/go-api-orchestrator/internal/data/request/media"
	"github.com/fahmihidayah/go-api-orchestrator/internal/domain"
	"github.com/fahmihidayah/go-api-orchestrator/internal/repository"
	"github.com/fahmihidayah/go-api-orchestrator/internal/storage"
	"github.com/fahmihidayah/go-api-orchestrator/internal/utils"
	"github.com/go-playground/validator/v10"
)

type IMediaService interface {
	Upload(ctx context.Context, media *request.UploadMediaRequest, file multipart.File) (*domain.Media, error)
	GetAll(ctx context.Context, page, limit int) ([]domain.Media, int64, error)
	GetByID(ctx context.Context, id string) (*domain.Media, error)
	GetByPath(ctx context.Context, path string) (*domain.Media, error)
	Update(ctx context.Context, id string, media *request.UpdateMediaRequest) (*domain.Media, error)
	Delete(ctx context.Context, id string) error
	DeleteAll(ctx context.Context, ids []string) error
	// React Admin specific methods
	GetAllReactAdmin(ctx context.Context, limit, offset int, sortField, sortOrder string, filters map[string]interface{}) ([]domain.Media, int64, error)
	GetByIDs(ctx context.Context, ids []string) ([]domain.Media, error)
	UpdateMany(ctx context.Context, ids []string, updates map[string]interface{}) ([]string, error)
	GetWithQueryParams(ctx context.Context, queryParams *utils.QueryParams) ([]domain.Media, int64, error)
}

type MediaServiceImpl struct {
	mediaRepository repository.IMediaRepository
	validator       *validator.Validate
	config          *config.Config
	storage         storage.Storage
}

func MediaServiceProvider(
	mediaRepository repository.IMediaRepository,
	config *config.Config,
	storage storage.Storage,
) IMediaService {
	return &MediaServiceImpl{
		mediaRepository: mediaRepository,
		validator:       validator.New(),
		config:          config,
		storage:         storage,
	}
}

func (s *MediaServiceImpl) Upload(ctx context.Context, media *request.UploadMediaRequest, file multipart.File) (*domain.Media, error) {
	// Validate request
	if err := s.validator.Struct(media); err != nil {
		return nil, err
	}

	// Upload file using storage interface
	fileInfo, err := s.storage.Upload(file, media.Media)
	if err != nil {
		return nil, err
	}

	// Create media record
	data := &domain.Media{
		ID:       utils.GenerateUUID(),
		Alt:      media.Alt,
		Url:      fileInfo.URL,
		Path:     fileInfo.Path,
		FileName: fileInfo.OriginalName,
		MimeType: fileInfo.MimeType,
		FileSize: fileInfo.FileSize,
		Width:    fileInfo.Width,
		Height:   fileInfo.Height,
	}

	// Save to repository
	if err := s.mediaRepository.Create(ctx, data); err != nil {
		// Delete uploaded file if database save fails
		s.storage.Delete(fileInfo.Path)
		return nil, err
	}

	return data, nil
}

func (s *MediaServiceImpl) GetAll(ctx context.Context, page, limit int) ([]domain.Media, int64, error) {
	// Set defaults
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	// Calculate offset
	offset := (page - 1) * limit

	// Get media from repository
	media, err := s.mediaRepository.GetAll(ctx, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	// Get total count
	count, err := s.mediaRepository.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	return media, count, nil
}

func (s *MediaServiceImpl) GetByID(ctx context.Context, id string) (*domain.Media, error) {
	return s.mediaRepository.GetByID(ctx, id)
}

func (s *MediaServiceImpl) GetByPath(ctx context.Context, path string) (*domain.Media, error) {
	return s.mediaRepository.GetByPath(ctx, path)
}

func (s *MediaServiceImpl) Update(ctx context.Context, id string, media *request.UpdateMediaRequest) (*domain.Media, error) {
	// Validate request
	if err := s.validator.Struct(media); err != nil {
		return nil, err
	}

	// Check if media exists
	existingMedia, err := s.mediaRepository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update fields
	existingMedia.Alt = media.Alt

	// Save to repository
	if err := s.mediaRepository.Update(ctx, existingMedia); err != nil {
		return nil, err
	}

	return existingMedia, nil
}

func (s *MediaServiceImpl) Delete(ctx context.Context, id string) error {
	// Check if media exists
	media, err := s.mediaRepository.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Delete file from storage
	if err := s.storage.Delete(media.Path); err != nil {
		// Log error but continue with database deletion
		// In production, you might want to handle this differently
	}

	// Delete from database
	return s.mediaRepository.Delete(ctx, id)
}

func (s *MediaServiceImpl) DeleteAll(ctx context.Context, ids []string) error {
	// Validate IDs array
	if len(ids) == 0 {
		return errors.New("IDs array cannot be empty")
	}

	// Get media records to delete files from storage
	for _, id := range ids {
		media, err := s.mediaRepository.GetByID(ctx, id)
		if err != nil {
			continue // Skip if not found
		}
		// Delete file from storage (ignore errors)
		s.storage.Delete(media.Path)
	}

	// Delete from database
	return s.mediaRepository.DeleteAll(ctx, ids)
}

// GetAllReactAdmin retrieves media with React Admin parameters
func (s *MediaServiceImpl) GetAllReactAdmin(ctx context.Context, limit, offset int, sortField, sortOrder string, filters map[string]interface{}) ([]domain.Media, int64, error) {
	// For now, use the basic GetAll method
	// In a more advanced implementation, you would apply sortField, sortOrder, and filters
	media, err := s.mediaRepository.GetAll(ctx, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	count, err := s.mediaRepository.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	return media, count, nil
}

// GetByIDs retrieves multiple media by their IDs
func (s *MediaServiceImpl) GetByIDs(ctx context.Context, ids []string) ([]domain.Media, error) {
	if len(ids) == 0 {
		return []domain.Media{}, nil
	}

	mediaList := make([]domain.Media, 0, len(ids))
	for _, id := range ids {
		media, err := s.mediaRepository.GetByID(ctx, id)
		if err != nil {
			continue // Skip media that don't exist
		}
		mediaList = append(mediaList, *media)
	}

	return mediaList, nil
}

// UpdateMany updates multiple media with the same data
func (s *MediaServiceImpl) UpdateMany(ctx context.Context, ids []string, updates map[string]interface{}) ([]string, error) {
	if len(ids) == 0 {
		return nil, errors.New("IDs array cannot be empty")
	}

	updatedIDs := make([]string, 0, len(ids))

	for _, id := range ids {
		media, err := s.mediaRepository.GetByID(ctx, id)
		if err != nil {
			continue // Skip media that don't exist
		}

		// Apply updates
		if alt, ok := updates["alt"].(string); ok {
			media.Alt = alt
		}

		// Save changes
		if err := s.mediaRepository.Update(ctx, media); err != nil {
			continue // Skip if update fails
		}

		updatedIDs = append(updatedIDs, id)
	}

	if len(updatedIDs) == 0 {
		return nil, errors.New("no media were updated")
	}

	return updatedIDs, nil
}

// GetWithQueryParams retrieves media with React Admin parameters
func (s *MediaServiceImpl) GetWithQueryParams(ctx context.Context, queryParams *utils.QueryParams) ([]domain.Media, int64, error) {
	count, err := s.mediaRepository.CountByQuery(ctx, queryParams)
	if err != nil {
		return nil, 0, err
	}

	queryParams.FillNextPrevTotal(count)

	media, err := s.mediaRepository.GetWithQuery(ctx, queryParams)
	if err != nil {
		return nil, 0, err
	}

	return media, count, nil
}
