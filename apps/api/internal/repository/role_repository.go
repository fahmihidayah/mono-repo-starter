package repository

import (
	"context"
	"errors"

	"github.com/fahmihidayah/go-api-orchestrator/internal/domain"
	"github.com/fahmihidayah/go-api-orchestrator/internal/utils"
	"gorm.io/gorm"
)

type IRoleRepository interface {
	Create(ctx context.Context, data *domain.Role) error
	GetByID(ctx context.Context, id string) (*domain.Role, error)
	GetAll(ctx context.Context, limit, offset int) ([]domain.Role, error)
	GetAllIds(ctx context.Context, id []string) []*domain.Role
	GetWithQuery(ctx context.Context, queryParams *utils.QueryParams) ([]*domain.Role, error)
	Update(ctx context.Context, data *domain.Role) error
	Delete(ctx context.Context, id string) error
	DeleteAll(ctx context.Context, ids []string) error
	Count(ctx context.Context) (int64, error)
	CountByQuery(ctx context.Context, queryParams *utils.QueryParams) (int64, error)
	AddPermissions(ctx context.Context, roleID string, permissions []*domain.Permission) error
	RemovePermissions(ctx context.Context, roleID string, permissions []*domain.Permission) error
	GetPermissions(ctx context.Context, roleID string) ([]*domain.Permission, error)
}

type RoleRepositoryImpl struct {
	db *gorm.DB
}

func RoleRepositoryProvider(db *gorm.DB) IRoleRepository {
	return &RoleRepositoryImpl{db: db}
}

func (r *RoleRepositoryImpl) Create(ctx context.Context, data *domain.Role) error {
	return r.db.WithContext(ctx).Create(data).Error
}

func (r *RoleRepositoryImpl) GetByID(ctx context.Context, id string) (*domain.Role, error) {
	var data domain.Role
	if err := r.db.WithContext(ctx).First(&data, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("role not found")
		}
		return nil, err
	}
	return &data, nil
}

func (r *RoleRepositoryImpl) GetAll(ctx context.Context, limit, offset int) ([]domain.Role, error) {
	var data []domain.Role
	err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&data).Error
	return data, err
}

func (r *RoleRepositoryImpl) GetAllIds(ctx context.Context, id []string) []*domain.Role {
	var data []*domain.Role
	r.db.WithContext(ctx).Where("id IN ?", id).Find(&data)
	return data
}

func (r *RoleRepositoryImpl) GetWithQuery(ctx context.Context, queryParams *utils.QueryParams) ([]*domain.Role, error) {
	var data []*domain.Role
	tx := r.db.WithContext(ctx).Limit(queryParams.Limit).Offset(queryParams.Offset)

	if len(queryParams.Sort) >= 2 {
		tx = tx.Order(queryParams.Sort[0] + " " + queryParams.Sort[1])
	}

	for key, value := range queryParams.Filter {
		if key == "ids" || key == "id" {
			if ids, ok := queryParams.GetFilterIDs(); ok {
				tx = tx.Where("id IN ?", ids)
			}
		} else {
			if strValue, ok := value.(string); ok {
				tx = tx.Where(key+" LIKE ?", "%"+strValue+"%")
			}
		}
	}

	err := tx.Find(&data).Error
	return data, err
}

func (r *RoleRepositoryImpl) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Delete(&domain.Role{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("role not found")
	}
	return nil
}

func (r *RoleRepositoryImpl) DeleteAll(ctx context.Context, ids []string) error {
	result := r.db.WithContext(ctx).Where("id IN ?", ids).Delete(&domain.Role{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no roles found with the provided IDs")
	}
	return nil
}

func (r *RoleRepositoryImpl) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&domain.Role{}).Count(&count).Error
	return count, err
}

func (r *RoleRepositoryImpl) CountByQuery(ctx context.Context, queryParams *utils.QueryParams) (int64, error) {
	var count int64
	tx := r.db.WithContext(ctx).Model(&domain.Role{})

	for key, value := range queryParams.Filter {
		if key == "ids" || key == "id" {
			if ids, ok := queryParams.GetFilterIDs(); ok {
				tx = tx.Where("id IN ?", ids)
			}
		} else {
			if strValue, ok := value.(string); ok {
				tx = tx.Where(key+" LIKE ?", "%"+strValue+"%")
			}
		}
	}

	err := tx.Count(&count).Error
	return count, err
}

func (r *RoleRepositoryImpl) Update(ctx context.Context, data *domain.Role) error {
	result := r.db.WithContext(ctx).Save(data)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("role not found or no changes made")
	}
	return nil
}

func (r *RoleRepositoryImpl) AddPermissions(ctx context.Context, roleID string, permissions []*domain.Permission) error {
	var role domain.Role
	if err := r.db.WithContext(ctx).First(&role, "id = ?", roleID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("role not found")
		}
		return err
	}

	return r.db.WithContext(ctx).Model(&role).Association("Permissions").Append(permissions)
}

func (r *RoleRepositoryImpl) RemovePermissions(ctx context.Context, roleID string, permissions []*domain.Permission) error {
	var role domain.Role
	if err := r.db.WithContext(ctx).First(&role, "id = ?", roleID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("role not found")
		}
		return err
	}

	return r.db.WithContext(ctx).Model(&role).Association("Permissions").Delete(permissions)
}

func (r *RoleRepositoryImpl) GetPermissions(ctx context.Context, roleID string) ([]*domain.Permission, error) {
	var role domain.Role
	if err := r.db.WithContext(ctx).Preload("Permissions").First(&role, "id = ?", roleID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("role not found")
		}
		return nil, err
	}

	return role.Permissions, nil
} 