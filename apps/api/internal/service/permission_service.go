package service

import (
	"context"
	"errors"
	"regexp"

	"github.com/fahmihidayah/go-api-orchestrator/internal/config"
	request "github.com/fahmihidayah/go-api-orchestrator/internal/data/request"
	"github.com/fahmihidayah/go-api-orchestrator/internal/domain"
	"github.com/fahmihidayah/go-api-orchestrator/internal/repository"
	"github.com/fahmihidayah/go-api-orchestrator/internal/utils"
	"github.com/go-playground/validator/v10"
)

type IPermissionService interface {
	Create(ctx context.Context, permission *request.CreatePermissionRequest) (*domain.Permission, error)
	GetAll(ctx context.Context, page, limit int) ([]domain.Permission, int64, error)
	GetByID(ctx context.Context, id string) (*domain.Permission, error)
	GetByPermission(ctx context.Context, permission string) (*domain.Permission, error)
	Update(ctx context.Context, id string, permission *request.CreatePermissionRequest) (*domain.Permission, error)
	Delete(ctx context.Context, id string) error
	DeleteAll(ctx context.Context, ids []string) error
	GetWithQueryParams(ctx context.Context, queryParams *utils.QueryParams) ([]*domain.Permission, *utils.PaginateInfo, error)
	GetByIDs(ctx context.Context, ids []string) ([]domain.Permission, error)
	UpdateMany(ctx context.Context, ids []string, updates map[string]interface{}) ([]string, error)
}

type PermissionServiceImpl struct {
	permissionRepository repository.IPermissionRepository
	validator            *validator.Validate
	config               *config.Config
}

func PermissionServiceProvider(
	permissionRepository repository.IPermissionRepository,
	config *config.Config,
) IPermissionService {
	v := validator.New()

	// Register custom validator for permission format
	v.RegisterValidation("permission_format", validatePermissionFormat)

	return &PermissionServiceImpl{
		permissionRepository: permissionRepository,
		validator:            v,
		config:               config,
	}
}

// validatePermissionFormat validates that permission follows the format [table]:[read/create/update/delete]
func validatePermissionFormat(fl validator.FieldLevel) bool {
	permission := fl.Field().String()

	// Regex pattern: [table_name]:[read|create|update|delete]
	pattern := `^[a-zA-Z_]+:(read|create|update|delete)$`
	matched, err := regexp.MatchString(pattern, permission)
	if err != nil {
		return false
	}

	return matched
}

func (s *PermissionServiceImpl) Create(ctx context.Context, permission *request.CreatePermissionRequest) (*domain.Permission, error) {
	// Validate request
	if err := s.validator.Struct(permission); err != nil {
		return nil, err
	}

	// Check if permission already exists
	existing, _ := s.permissionRepository.GetByPermission(ctx, permission.Permission)
	if existing != nil {
		return nil, errors.New("permission already exists")
	}

	// Map request to domain model
	data := &domain.Permission{
		ID:         utils.GenerateUUID(),
		Permission: permission.Permission,
	}

	// Save to repository
	if err := s.permissionRepository.Create(ctx, data); err != nil {
		return nil, err
	}

	return data, nil
}

func (s *PermissionServiceImpl) GetAll(ctx context.Context, page, limit int) ([]domain.Permission, int64, error) {
	// Set defaults
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	// Calculate offset
	offset := (page - 1) * limit

	// Get permissions from repository
	permissions, err := s.permissionRepository.GetAll(ctx, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	// Get total count
	count, err := s.permissionRepository.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	return permissions, count, nil
}

func (s *PermissionServiceImpl) GetByID(ctx context.Context, id string) (*domain.Permission, error) {
	return s.permissionRepository.GetByID(ctx, id)
}

func (s *PermissionServiceImpl) GetByPermission(ctx context.Context, permission string) (*domain.Permission, error) {
	return s.permissionRepository.GetByPermission(ctx, permission)
}

func (s *PermissionServiceImpl) Update(ctx context.Context, id string, permission *request.CreatePermissionRequest) (*domain.Permission, error) {
	// Validate request
	if err := s.validator.Struct(permission); err != nil {
		return nil, err
	}

	// Check if permission exists
	existingPermission, err := s.permissionRepository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Check if new permission value conflicts with another permission
	if existingPermission.Permission != permission.Permission {
		conflicting, _ := s.permissionRepository.GetByPermission(ctx, permission.Permission)
		if conflicting != nil {
			return nil, errors.New("permission already exists")
		}
	}

	// Update fields
	existingPermission.Permission = permission.Permission

	// Save to repository
	if err := s.permissionRepository.Update(ctx, existingPermission); err != nil {
		return nil, err
	}

	return existingPermission, nil
}

func (s *PermissionServiceImpl) Delete(ctx context.Context, id string) error {
	// Check if permission exists
	_, err := s.permissionRepository.GetByID(ctx, id)
	if err != nil {
		return err
	}

	return s.permissionRepository.Delete(ctx, id)
}

func (s *PermissionServiceImpl) DeleteAll(ctx context.Context, ids []string) error {
	// Validate IDs array
	if len(ids) == 0 {
		return errors.New("IDs array cannot be empty")
	}

	return s.permissionRepository.DeleteAll(ctx, ids)
}

func (s *PermissionServiceImpl) GetWithQueryParams(ctx context.Context, queryParams *utils.QueryParams) ([]*domain.Permission, *utils.PaginateInfo, error) {
	count, err := s.permissionRepository.CountByQuery(ctx, queryParams)
	if err != nil {
		return nil, nil, err
	}

	paginateInfo := queryParams.ToPaginateInfo(count)
	permissions, err := s.permissionRepository.GetWithQuery(ctx, queryParams)
	if err != nil {
		return nil, nil, err
	}

	return permissions, paginateInfo, nil
}

func (s *PermissionServiceImpl) GetByIDs(ctx context.Context, ids []string) ([]domain.Permission, error) {
	if len(ids) == 0 {
		return []domain.Permission{}, nil
	}

	permissions := make([]domain.Permission, 0, len(ids))
	for _, id := range ids {
		permission, err := s.permissionRepository.GetByID(ctx, id)
		if err != nil {
			continue // Skip permissions that don't exist
		}
		permissions = append(permissions, *permission)
	}

	return permissions, nil
}

func (s *PermissionServiceImpl) UpdateMany(ctx context.Context, ids []string, updates map[string]interface{}) ([]string, error) {
	if len(ids) == 0 {
		return nil, errors.New("IDs array cannot be empty")
	}

	updatedIDs := make([]string, 0, len(ids))

	for _, id := range ids {
		permission, err := s.permissionRepository.GetByID(ctx, id)
		if err != nil {
			continue // Skip permissions that don't exist
		}

		// Apply updates
		if permValue, ok := updates["permission"].(string); ok && permValue != "" {
			// Validate the format
			pattern := `^[a-zA-Z_]+:(read|create|update|delete)$`
			matched, _ := regexp.MatchString(pattern, permValue)
			if !matched {
				continue // Skip invalid format
			}
			permission.Permission = permValue
		}

		// Save changes
		if err := s.permissionRepository.Update(ctx, permission); err != nil {
			continue // Skip if update fails
		}

		updatedIDs = append(updatedIDs, id)
	}

	if len(updatedIDs) == 0 {
		return nil, errors.New("no permissions were updated")
	}

	return updatedIDs, nil
}
