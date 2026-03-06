package service

import (
	"context"
	"errors"

	"github.com/fahmihidayah/go-api-orchestrator/internal/config"
	request "github.com/fahmihidayah/go-api-orchestrator/internal/data/request"
	"github.com/fahmihidayah/go-api-orchestrator/internal/domain"
	"github.com/fahmihidayah/go-api-orchestrator/internal/repository"
	"github.com/fahmihidayah/go-api-orchestrator/internal/utils"
	"github.com/go-playground/validator/v10"
)

type ICategoryService interface {
	Create(ctx context.Context, category *request.CreateCategoryRequest) (*domain.Category, error)
	GetAll(ctx context.Context, page, limit int) ([]domain.Category, int64, error)
	GetByID(ctx context.Context, id string) (*domain.Category, error)
	GetBySlug(ctx context.Context, slug string) (*domain.Category, error)
	Update(ctx context.Context, id string, category *request.CreateCategoryRequest) (*domain.Category, error)
	Delete(ctx context.Context, id string) error
	DeleteAll(ctx context.Context, ids []string) error
	// React Admin specific methods
	GetWithQueryParams(ctx context.Context, queryParams *utils.QueryParams) ([]*domain.Category, *utils.PaginateInfo, error)
	GetByIDs(ctx context.Context, ids []string) ([]domain.Category, error)
	UpdateMany(ctx context.Context, ids []string, updates map[string]interface{}) ([]string, error)
}

type CategoryServiceImpl struct {
	categoryRepository repository.ICategoryRepository
	validator          *validator.Validate
	config             *config.Config
}

func CategoryServiceProvider(
	categoryRepository repository.ICategoryRepository,
	config *config.Config,
) ICategoryService {
	return &CategoryServiceImpl{
		categoryRepository: categoryRepository,
		validator:          validator.New(),
		config:             config,
	}
}

func (s *CategoryServiceImpl) Create(ctx context.Context, category *request.CreateCategoryRequest) (*domain.Category, error) {
	// Validate request
	if err := s.validator.Struct(category); err != nil {
		return nil, err
	}

	// Map request to domain model
	data := &domain.Category{
		ID:    utils.GenerateUUID(),
		Slug:  utils.GenerateSlug(category.Title),
		Title: category.Title,
	}

	// Save to repository
	if err := s.categoryRepository.Create(ctx, data); err != nil {
		return nil, err
	}

	return data, nil
}

func (s *CategoryServiceImpl) GetAll(ctx context.Context, page, limit int) ([]domain.Category, int64, error) {
	// Set defaults
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	// Calculate offset
	offset := (page - 1) * limit

	// Get categories from repository
	categories, err := s.categoryRepository.GetAll(ctx, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	// Get total count
	count, err := s.categoryRepository.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	return categories, count, nil
}

func (s *CategoryServiceImpl) GetByID(ctx context.Context, id string) (*domain.Category, error) {
	return s.categoryRepository.GetByID(ctx, id)
}

func (s *CategoryServiceImpl) GetBySlug(ctx context.Context, slug string) (*domain.Category, error) {
	return s.categoryRepository.FindBySlug(ctx, slug)
}

func (s *CategoryServiceImpl) Update(ctx context.Context, id string, category *request.CreateCategoryRequest) (*domain.Category, error) {
	// Validate request
	if err := s.validator.Struct(category); err != nil {
		return nil, err
	}

	// Check if category exists
	existingCategory, err := s.categoryRepository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update fields
	existingCategory.Title = category.Title
	existingCategory.Slug = utils.GenerateSlug(category.Title)

	// Save to repository
	if err := s.categoryRepository.Update(ctx, existingCategory); err != nil {
		return nil, err
	}

	return existingCategory, nil
}

func (s *CategoryServiceImpl) Delete(ctx context.Context, id string) error {
	// Check if category exists
	_, err := s.categoryRepository.GetByID(ctx, id)
	if err != nil {
		return err
	}

	return s.categoryRepository.Delete(ctx, id)
}

func (s *CategoryServiceImpl) DeleteAll(ctx context.Context, ids []string) error {
	// Validate IDs array
	if len(ids) == 0 {
		return errors.New("IDs array cannot be empty")
	}

	return s.categoryRepository.DeleteAll(ctx, ids)
}

// GetWithQueryParams retrieves categories with React Admin parameters
func (s *CategoryServiceImpl) GetWithQueryParams(ctx context.Context, queryParams *utils.QueryParams) ([]*domain.Category, *utils.PaginateInfo, error) {
	count, err := s.categoryRepository.CountByQuery(ctx, queryParams)
	if err != nil {
		return nil, nil, err
	}

	paginateInfo := queryParams.ToPaginateInfo(count)
	categories, err := s.categoryRepository.GetWithQuery(ctx, queryParams)
	if err != nil {
		return nil, nil, err
	}

	return categories, paginateInfo, nil
}

// GetByIDs retrieves multiple categories by their IDs
func (s *CategoryServiceImpl) GetByIDs(ctx context.Context, ids []string) ([]domain.Category, error) {
	if len(ids) == 0 {
		return []domain.Category{}, nil
	}

	categories := make([]domain.Category, 0, len(ids))
	for _, id := range ids {
		category, err := s.categoryRepository.GetByID(ctx, id)
		if err != nil {
			continue // Skip categories that don't exist
		}
		categories = append(categories, *category)
	}

	return categories, nil
}

// UpdateMany updates multiple categories with the same data
func (s *CategoryServiceImpl) UpdateMany(ctx context.Context, ids []string, updates map[string]interface{}) ([]string, error) {
	if len(ids) == 0 {
		return nil, errors.New("IDs array cannot be empty")
	}

	updatedIDs := make([]string, 0, len(ids))

	for _, id := range ids {
		category, err := s.categoryRepository.GetByID(ctx, id)
		if err != nil {
			continue // Skip categories that don't exist
		}

		// Apply updates
		if name, ok := updates["name"].(string); ok && name != "" {
			category.Title = name
			category.Slug = utils.GenerateSlug(name)
		}

		// Save changes
		if err := s.categoryRepository.Update(ctx, category); err != nil {
			continue // Skip if update fails
		}

		updatedIDs = append(updatedIDs, id)
	}

	if len(updatedIDs) == 0 {
		return nil, errors.New("no categories were updated")
	}

	return updatedIDs, nil
}
