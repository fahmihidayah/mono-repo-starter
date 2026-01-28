package repository

import (
	"context"
	"errors"

	"github.com/fahmihidayah/go-api-orchestrator/internal/domain"
	"github.com/fahmihidayah/go-api-orchestrator/internal/utils"
	"gorm.io/gorm"
)

type ICategoryRepository interface {
	// Define methods for category repository
	FindBySlug(ctx context.Context, slug string) (*domain.Category, error)
	Create(ctx context.Context, data *domain.Category) error
	GetByID(ctx context.Context, id string) (*domain.Category, error)
	GetAll(ctx context.Context, limit, offset int) ([]domain.Category, error)
	GetAllIds(ctx context.Context, id []string) []*domain.Category
	GetWithQuery(ctx context.Context, queryParams *utils.QueryParams) ([]*domain.Category, error)
	Update(ctx context.Context, data *domain.Category) error
	Delete(ctx context.Context, id string) error
	DeleteAll(ctx context.Context, ids []string) error
	Count(ctx context.Context) (int64, error)
	CountByQuery(ctx context.Context, queryParams *utils.QueryParams) (int64, error)
}

type CategoryRepositoryImpl struct {
	db *gorm.DB
}

func CategoryRepositoryProvider(db *gorm.DB) ICategoryRepository {
	return &CategoryRepositoryImpl{db: db}
}

func (r *CategoryRepositoryImpl) Create(ctx context.Context, data *domain.Category) error {
	return r.db.WithContext(ctx).Create(data).Error
}

func (r *CategoryRepositoryImpl) GetByID(ctx context.Context, id string) (*domain.Category, error) {
	var data domain.Category
	if err := r.db.WithContext(ctx).First(&data, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("category not found")
		}
		return nil, err
	}
	return &data, nil
}

func (r *CategoryRepositoryImpl) GetAll(ctx context.Context, limit, offset int) ([]domain.Category, error) {
	var data []domain.Category
	err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&data).Error
	return data, err
}

func (r *CategoryRepositoryImpl) GetAllIds(ctx context.Context, id []string) []*domain.Category {
	var data []*domain.Category
	r.db.WithContext(ctx).Where("id IN ?", id).Find(&data)
	return data
}

func (r *CategoryRepositoryImpl) GetWithQuery(ctx context.Context, queryParams *utils.QueryParams) ([]*domain.Category, error) {
	var data []*domain.Category
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
func (r *CategoryRepositoryImpl) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Delete(&domain.Category{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("category not found")
	}
	return nil
}

func (r *CategoryRepositoryImpl) DeleteAll(ctx context.Context, ids []string) error {
	result := r.db.WithContext(ctx).Where("id IN ?", ids).Delete(&domain.Category{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no categories found with the provided IDs")
	}
	return nil
}

func (r *CategoryRepositoryImpl) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&domain.Category{}).Count(&count).Error
	return count, err
}

func (r *CategoryRepositoryImpl) CountByQuery(ctx context.Context, queryParams *utils.QueryParams) (int64, error) {
	var count int64
	tx := r.db.WithContext(ctx).Model(&domain.Category{})

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

func (r *CategoryRepositoryImpl) Update(ctx context.Context, data *domain.Category) error {
	result := r.db.WithContext(ctx).Save(data)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("category not found or no changes made")
	}
	return nil
}

func (r *CategoryRepositoryImpl) FindBySlug(ctx context.Context, slug string) (*domain.Category, error) {
	var data domain.Category
	if err := r.db.WithContext(ctx).Where("slug = ?", slug).First(&data).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("category not found")
		}
		return nil, err
	}
	return &data, nil
}
