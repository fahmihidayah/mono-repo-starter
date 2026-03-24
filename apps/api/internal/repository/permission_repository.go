package repository

import (
	"context"
	"errors"

	"github.com/fahmihidayah/go-api-orchestrator/internal/domain"
	"github.com/fahmihidayah/go-api-orchestrator/internal/utils"
	"gorm.io/gorm"
)

type IPermissionRepository interface {
	Create(ctx context.Context, data *domain.Permission) error
	GetByID(ctx context.Context, id string) (*domain.Permission, error)
	GetByPermission(ctx context.Context, permission string) (*domain.Permission, error)
	GetAll(ctx context.Context, limit, offset int) ([]domain.Permission, error)
	GetAllIds(ctx context.Context, id []string) []*domain.Permission
	GetWithQuery(ctx context.Context, queryParams *utils.QueryParams) ([]*domain.Permission, error)
	Update(ctx context.Context, data *domain.Permission) error
	Delete(ctx context.Context, id string) error
	DeleteAll(ctx context.Context, ids []string) error
	Count(ctx context.Context) (int64, error)
	CountByQuery(ctx context.Context, queryParams *utils.QueryParams) (int64, error)
}

type PermissionRepositoryImpl struct {
	db *gorm.DB
}

func PermissionRepositoryProvider(db *gorm.DB) IPermissionRepository {
	return &PermissionRepositoryImpl{db: db}
}

func (r *PermissionRepositoryImpl) Create(ctx context.Context, data *domain.Permission) error {
	return r.db.WithContext(ctx).Create(data).Error
}

func (r *PermissionRepositoryImpl) GetByID(ctx context.Context, id string) (*domain.Permission, error) {
	var data domain.Permission
	if err := r.db.WithContext(ctx).Preload("Roles").First(&data, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("permission not found")
		}
		return nil, err
	}
	return &data, nil
}

func (r *PermissionRepositoryImpl) GetByPermission(ctx context.Context, permission string) (*domain.Permission, error) {
	var data domain.Permission
	if err := r.db.WithContext(ctx).Where("permission = ?", permission).First(&data).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("permission not found")
		}
		return nil, err
	}
	return &data, nil
}

func (r *PermissionRepositoryImpl) GetAll(ctx context.Context, limit, offset int) ([]domain.Permission, error) {
	var data []domain.Permission
	err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&data).Error
	return data, err
}

func (r *PermissionRepositoryImpl) GetAllIds(ctx context.Context, id []string) []*domain.Permission {
	var data []*domain.Permission
	r.db.WithContext(ctx).Where("id IN ?", id).Find(&data)
	return data
}

func (r *PermissionRepositoryImpl) GetWithQuery(ctx context.Context, queryParams *utils.QueryParams) ([]*domain.Permission, error) {
	var data []*domain.Permission
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

func (r *PermissionRepositoryImpl) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Delete(&domain.Permission{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("permission not found")
	}
	return nil
}

func (r *PermissionRepositoryImpl) DeleteAll(ctx context.Context, ids []string) error {
	result := r.db.WithContext(ctx).Where("id IN ?", ids).Delete(&domain.Permission{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no permissions found with the provided IDs")
	}
	return nil
}

func (r *PermissionRepositoryImpl) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&domain.Permission{}).Count(&count).Error
	return count, err
}

func (r *PermissionRepositoryImpl) CountByQuery(ctx context.Context, queryParams *utils.QueryParams) (int64, error) {
	var count int64
	tx := r.db.WithContext(ctx).Model(&domain.Permission{})

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

func (r *PermissionRepositoryImpl) Update(ctx context.Context, data *domain.Permission) error {
	result := r.db.WithContext(ctx).Save(data)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("permission not found or no changes made")
	}
	return nil
}
