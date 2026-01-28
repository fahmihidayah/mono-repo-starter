package repository

import (
	"context"
	"errors"

	"github.com/fahmihidayah/go-api-orchestrator/internal/domain"
	"github.com/fahmihidayah/go-api-orchestrator/internal/utils"
	"gorm.io/gorm"
)

type IMediaRepository interface {
	Create(ctx context.Context, data *domain.Media) error
	GetByID(ctx context.Context, id string) (*domain.Media, error)
	GetByPath(ctx context.Context, path string) (*domain.Media, error)
	GetAll(ctx context.Context, limit, offset int) ([]domain.Media, error)
	GetAllIds(ctx context.Context, id []string) []*domain.Media
	Update(ctx context.Context, data *domain.Media) error
	Delete(ctx context.Context, id string) error
	DeleteAll(ctx context.Context, ids []string) error
	Count(ctx context.Context) (int64, error)
	GetWithQuery(ctx context.Context, queryParams *utils.QueryParams) ([]domain.Media, error)
	CountByQuery(ctx context.Context, queryParams *utils.QueryParams) (int64, error)
}

type MediaRepositoryImpl struct {
	db *gorm.DB
}

func MediaRepositoryProvider(db *gorm.DB) IMediaRepository {
	return &MediaRepositoryImpl{db: db}
}

func (r *MediaRepositoryImpl) Create(ctx context.Context, data *domain.Media) error {
	return r.db.WithContext(ctx).Create(data).Error
}

func (r *MediaRepositoryImpl) GetByID(ctx context.Context, id string) (*domain.Media, error) {
	var data domain.Media
	if err := r.db.WithContext(ctx).First(&data, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("media not found")
		}
		return nil, err
	}
	return &data, nil
}

func (r *MediaRepositoryImpl) GetAllIds(ctx context.Context, id []string) []*domain.Media {
	var data []*domain.Media
	r.db.WithContext(ctx).Where("id IN ?", id).Find(&data)
	return data
}

func (r *MediaRepositoryImpl) GetByPath(ctx context.Context, path string) (*domain.Media, error) {
	var data domain.Media
	if err := r.db.WithContext(ctx).Where("path = ?", path).First(&data).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("media not found")
		}
		return nil, err
	}
	return &data, nil
}

func (r *MediaRepositoryImpl) GetAll(ctx context.Context, limit, offset int) ([]domain.Media, error) {
	var data []domain.Media
	err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Order("created_at DESC").Find(&data).Error
	return data, err
}

func (r *MediaRepositoryImpl) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Delete(&domain.Media{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("media not found")
	}
	return nil
}

func (r *MediaRepositoryImpl) DeleteAll(ctx context.Context, ids []string) error {
	result := r.db.WithContext(ctx).Where("id IN ?", ids).Delete(&domain.Media{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no media found with the provided IDs")
	}
	return nil
}

func (r *MediaRepositoryImpl) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&domain.Media{}).Count(&count).Error
	return count, err
}

func (r *MediaRepositoryImpl) Update(ctx context.Context, data *domain.Media) error {
	result := r.db.WithContext(ctx).Save(data)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("media not found or no changes made")
	}
	return nil
}

// GetWithQuery retrieves media with query parameters
func (r *MediaRepositoryImpl) GetWithQuery(ctx context.Context, queryParams *utils.QueryParams) ([]domain.Media, error) {
	var data []domain.Media
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

// CountByQuery returns the count of media matching the query parameters
func (r *MediaRepositoryImpl) CountByQuery(ctx context.Context, queryParams *utils.QueryParams) (int64, error) {
	var count int64
	tx := r.db.WithContext(ctx).Model(&domain.Media{})

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
