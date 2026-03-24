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

type IRoleService interface {
	Create(ctx context.Context, role *request.CreateRoleRequest) (*domain.Role, error)
	GetAll(ctx context.Context, page, limit int) ([]domain.Role, int64, error)
	GetByID(ctx context.Context, id string) (*domain.Role, error)
	Update(ctx context.Context, id string, role *request.CreateRoleRequest) (*domain.Role, error)
	Delete(ctx context.Context, id string) error
	DeleteAll(ctx context.Context, ids []string) error
	// React Admin specific methods
	GetWithQueryParams(ctx context.Context, queryParams *utils.QueryParams) ([]*domain.Role, *utils.PaginateInfo, error)
	GetByIDs(ctx context.Context, ids []string) ([]domain.Role, error)
	UpdateMany(ctx context.Context, ids []string, updates map[string]interface{}) ([]string, error)
	// Permission management
	AddPermissions(ctx context.Context, roleID string, permissionIDs []string) error
	RemovePermissions(ctx context.Context, roleID string, permissionIDs []string) error
	GetPermissions(ctx context.Context, roleID string) ([]*domain.Permission, error)
}

type RoleServiceImpl struct {
	roleRepository       repository.IRoleRepository
	permissionRepository repository.IPermissionRepository
	validator            *validator.Validate
	config               *config.Config
}

func RoleServiceProvider(
	roleRepository repository.IRoleRepository,
	permissionRepository repository.IPermissionRepository,
	config *config.Config,
) IRoleService {
	return &RoleServiceImpl{
		roleRepository:       roleRepository,
		permissionRepository: permissionRepository,
		validator:            validator.New(),
		config:               config,
	}
}

func (s *RoleServiceImpl) Create(ctx context.Context, role *request.CreateRoleRequest) (*domain.Role, error) {
	// Validate request
	if err := s.validator.Struct(role); err != nil {
		return nil, err
	}

	// Map request to domain model
	data := &domain.Role{
		ID:   utils.GenerateUUID(),
		Name: role.Name,
	}

	// Save to repository
	if err := s.roleRepository.Create(ctx, data); err != nil {
		return nil, err
	}

	return data, nil
}

func (s *RoleServiceImpl) GetAll(ctx context.Context, page, limit int) ([]domain.Role, int64, error) {
	// Set defaults
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	// Calculate offset
	offset := (page - 1) * limit

	// Get roles from repository
	roles, err := s.roleRepository.GetAll(ctx, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	// Get total count
	count, err := s.roleRepository.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	return roles, count, nil
}

func (s *RoleServiceImpl) GetByID(ctx context.Context, id string) (*domain.Role, error) {
	return s.roleRepository.GetByID(ctx, id)
}

func (s *RoleServiceImpl) Update(ctx context.Context, id string, role *request.CreateRoleRequest) (*domain.Role, error) {
	// Validate request
	if err := s.validator.Struct(role); err != nil {
		return nil, err
	}

	// Check if role exists
	existingRole, err := s.roleRepository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update fields
	existingRole.Name = role.Name

	// Save to repository
	if err := s.roleRepository.Update(ctx, existingRole); err != nil {
		return nil, err
	}

	return existingRole, nil
}

func (s *RoleServiceImpl) Delete(ctx context.Context, id string) error {
	// Check if role exists
	_, err := s.roleRepository.GetByID(ctx, id)
	if err != nil {
		return err
	}

	return s.roleRepository.Delete(ctx, id)
}

func (s *RoleServiceImpl) DeleteAll(ctx context.Context, ids []string) error {
	// Validate IDs array
	if len(ids) == 0 {
		return errors.New("IDs array cannot be empty")
	}

	return s.roleRepository.DeleteAll(ctx, ids)
}

// GetWithQueryParams retrieves roles with React Admin parameters
func (s *RoleServiceImpl) GetWithQueryParams(ctx context.Context, queryParams *utils.QueryParams) ([]*domain.Role, *utils.PaginateInfo, error) {
	count, err := s.roleRepository.CountByQuery(ctx, queryParams)
	if err != nil {
		return nil, nil, err
	}

	paginateInfo := queryParams.ToPaginateInfo(count)
	roles, err := s.roleRepository.GetWithQuery(ctx, queryParams)
	if err != nil {
		return nil, nil, err
	}

	return roles, paginateInfo, nil
}

// GetByIDs retrieves multiple roles by their IDs
func (s *RoleServiceImpl) GetByIDs(ctx context.Context, ids []string) ([]domain.Role, error) {
	if len(ids) == 0 {
		return []domain.Role{}, nil
	}

	roles := make([]domain.Role, 0, len(ids))
	for _, id := range ids {
		role, err := s.roleRepository.GetByID(ctx, id)
		if err != nil {
			continue // Skip roles that don't exist
		}
		roles = append(roles, *role)
	}

	return roles, nil
}

// UpdateMany updates multiple roles with the same data
func (s *RoleServiceImpl) UpdateMany(ctx context.Context, ids []string, updates map[string]interface{}) ([]string, error) {
	if len(ids) == 0 {
		return nil, errors.New("IDs array cannot be empty")
	}

	updatedIDs := make([]string, 0, len(ids))

	for _, id := range ids {
		role, err := s.roleRepository.GetByID(ctx, id)
		if err != nil {
			continue // Skip roles that don't exist
		}

		// Apply updates
		if name, ok := updates["name"].(string); ok && name != "" {
			role.Name = name
		}

		// Save changes
		if err := s.roleRepository.Update(ctx, role); err != nil {
			continue // Skip if update fails
		}

		updatedIDs = append(updatedIDs, id)
	}

	if len(updatedIDs) == 0 {
		return nil, errors.New("no roles were updated")
	}

	return updatedIDs, nil
}

// AddPermissions adds permissions to a role
func (s *RoleServiceImpl) AddPermissions(ctx context.Context, roleID string, permissionIDs []string) error {
	if len(permissionIDs) == 0 {
		return errors.New("permission IDs array cannot be empty")
	}

	// Verify role exists
	_, err := s.roleRepository.GetByID(ctx, roleID)
	if err != nil {
		return err
	}

	// Get all permissions by their IDs
	permissions := make([]*domain.Permission, 0, len(permissionIDs))
	for _, permID := range permissionIDs {
		perm, err := s.permissionRepository.GetByID(ctx, permID)
		if err != nil {
			continue // Skip permissions that don't exist
		}
		permissions = append(permissions, perm)
	}

	if len(permissions) == 0 {
		return errors.New("no valid permissions found")
	}

	return s.roleRepository.AddPermissions(ctx, roleID, permissions)
}

// RemovePermissions removes permissions from a role
func (s *RoleServiceImpl) RemovePermissions(ctx context.Context, roleID string, permissionIDs []string) error {
	if len(permissionIDs) == 0 {
		return errors.New("permission IDs array cannot be empty")
	}

	// Verify role exists
	_, err := s.roleRepository.GetByID(ctx, roleID)
	if err != nil {
		return err
	}

	// Get all permissions by their IDs
	permissions := make([]*domain.Permission, 0, len(permissionIDs))
	for _, permID := range permissionIDs {
		perm, err := s.permissionRepository.GetByID(ctx, permID)
		if err != nil {
			continue // Skip permissions that don't exist
		}
		permissions = append(permissions, perm)
	}

	if len(permissions) == 0 {
		return errors.New("no valid permissions found")
	}

	return s.roleRepository.RemovePermissions(ctx, roleID, permissions)
}

// GetPermissions retrieves all permissions for a role
func (s *RoleServiceImpl) GetPermissions(ctx context.Context, roleID string) ([]*domain.Permission, error) {
	return s.roleRepository.GetPermissions(ctx, roleID)
}
